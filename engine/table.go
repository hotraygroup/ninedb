package engine

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"sync"
	"time"
)

type Table interface {
	TableName() string
}

func CreateTable(obj interface{}) {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	tableName := typ.Name()

	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("cannot use non-ptr struct %s", tableName))
	}

	log.Printf("tableName is %s", tableName)

	db.Lock()
	defer db.Unlock()

	if _, ok := db.tables[tableName]; ok {
		panic(fmt.Errorf("%s has been created", tableName))
	}

	db.tables[tableName] = make([]interface{}, 0) //store array
	db.chans[tableName] = make(chan int, ROWSIZE)
	db.rows[tableName] = make(map[int]int) //rid  -> array index
	db.locks[tableName] = sync.RWMutex{}
	db.sorting[tableName] = make(map[string]bool)

	//存在索引
	if indexs := obj.(Row).Index(); indexs != nil {
		db.indexs[tableName] = make(map[string][]int)
	}

}

////////////////////////////内部调用////////////////////////////////////////////////////////
func get(tableName string) int {
	//check table exist and lock outer
	select {
	case index := <-db.chans[tableName]:
		return index
	default:
		allocSize := len(db.tables[tableName])
		toAppend := make([]interface{}, ROWSIZE/2)
		db.tables[tableName] = append(db.tables[tableName], toAppend...)
		for i := 0; i < ROWSIZE/2; i++ {
			db.chans[tableName] <- allocSize + i
		}
		return <-db.chans[tableName]
	}
}

func put(tableName string, rid int) {
	//check table exist and lock outer
	select {
	case db.chans[tableName] <- rid:
		return
	default:
		log.Printf("table %s's chan is full", tableName)
	}

}

func sortIndex(tableName string, index string) {
	if db.sorting[tableName][index] == true {
		return
	}
	db.sorting[tableName][index] = true
	time.AfterFunc(5*time.Second, func() {
		start := time.Now().Unix()
		lock := db.locks[tableName]
		lock.Lock()
		defer lock.Unlock()
		var arr []int
		if index == "global" {
			arr = db.ids[tableName]
		} else {
			arr = db.indexs[tableName][index]
		}
		sort.IntSlice(arr).Sort()
		end := time.Now().Unix()
		log.Printf("sort index %s:%s %d records finished in %d second", tableName, index, len(arr), end-start)
		db.sorting[tableName][index] = false
	})

}

////////////////////////////////////////////////////////////////////////////////////
func Insert(obj interface{}) error {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	tableName := typ.Name()

	db.RLock()
	if _, ok := db.tables[tableName]; !ok {
		panic(fmt.Errorf("table %s is not exsit", tableName))
	}
	db.RUnlock()

	id := obj.(Row).GetID()

	lock := db.locks[tableName]
	lock.Lock()
	defer lock.Unlock()

	if _, ok := db.rows[tableName][id]; ok { //exist
		rid := db.rows[tableName][id]
		log.Printf("record id[%d] is exist in table %s %d row", id, tableName, rid)
		return fmt.Errorf("record %d is exist in %s", id, tableName)
	}

	rid := get(tableName)

	db.tables[tableName][rid] = obj
	db.rows[tableName][id] = rid
	//添加到pk ids
	if len(db.ids[tableName]) == 0 {
		db.sorting[tableName]["global"] = false
	}
	db.ids[tableName] = append(db.ids[tableName], id)
	sortIndex(tableName, "global")

	//log.Printf("insert record id[%d] in table %s's %d row", id, tableName, rid)

	indexs := obj.(Row).Index()
	if indexs == nil {
		return nil
	}

	//存在索引，创建索引
	for i := 0; i < len(indexs); i++ {
		if len(indexs[i]) == 0 {
			continue
		}
		pk := fmt.Sprintf("[%s][%d][%d]", tableName, len(indexs), len(indexs[i]))
		sort.StringSlice(indexs[i]).Sort()
		for j := 0; j < len(indexs[i]); j++ {
			pk += fmt.Sprintf(":%s:%v", indexs[i][j], reflect.Indirect(val).FieldByName(indexs[i][j]))
		}
		if len(db.indexs[tableName][pk]) == 0 {
			db.sorting[tableName][pk] = false
		}
		db.indexs[tableName][pk] = append(db.indexs[tableName][pk], id)

		sortIndex(tableName, pk)
	}
	return nil
}

//全覆盖更新
func Update(obj interface{}) error {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	tableName := typ.Name()

	db.RLock()
	if _, ok := db.tables[tableName]; !ok {
		panic(fmt.Errorf("table %s is not exsit", tableName))
	}
	db.RUnlock()

	id := obj.(Row).GetID()

	lock := db.locks[tableName]
	lock.Lock()
	defer lock.Unlock()

	if rid, ok := db.rows[tableName][id]; ok {
		db.tables[tableName][rid] = obj
		//log.Printf("update record id[%d] in table %s's %d row", id, tableName, rid)

	} else {
		log.Printf("record %d is not exist in table %s", id, tableName)
		return fmt.Errorf("record %d is not exist in table %s", id, tableName)
	}
	return nil
}

//回调更新，用户转账、打工等场景  todo
func UpdateFunc() error {
	return nil
}

func Get(obj interface{}) interface{} {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	tableName := typ.Name()

	db.RLock()
	if _, ok := db.tables[tableName]; !ok {
		panic(fmt.Errorf("table %s is not exsit", tableName))
	}
	db.RUnlock()

	id := obj.(Row).GetID()
	lock := db.locks[tableName]
	lock.Lock()
	defer lock.Unlock()

	if rid, ok := db.rows[tableName][id]; ok {
		return db.tables[tableName][rid]
	}

	log.Printf("record %d is not exist in table %s", id, tableName)
	return nil
}

func Delete(obj interface{}) {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	tableName := typ.Name()

	db.RLock()
	if _, ok := db.tables[tableName]; !ok {
		panic(fmt.Errorf("table %s is not exsit", tableName))
	}
	db.RUnlock()

	id := obj.(Row).GetID()
	lock := db.locks[tableName]
	lock.Lock()
	defer lock.Unlock()

	rid, ok := db.rows[tableName][id]
	if !ok {
		return
	}

	put(tableName, rid)
	delete(db.rows[tableName], id)
	//删除pk id
	for i := 0; i < len(db.ids[tableName]); i++ {
		if db.ids[tableName][i] == id {
			db.ids[tableName][i] = db.ids[tableName][len(db.ids[tableName])-1]
			db.ids[tableName] = db.ids[tableName][:len(db.ids[tableName])-1]
		}
	}
	sortIndex(tableName, "global")

	log.Printf("delete recoed %d from %s", id, tableName)

	indexs := obj.(Row).Index()
	if indexs == nil {
		return
	}

	//存在索引，删除索引
	for i := 0; i < len(indexs); i++ {
		if len(indexs[i]) == 0 {
			continue
		}
		pk := fmt.Sprintf("[%s][%d][%d]", tableName, len(indexs), len(indexs[i]))
		sort.StringSlice(indexs[i]).Sort()
		for j := 0; j < len(indexs[i]); j++ {
			pk += fmt.Sprintf(":%s:%v", indexs[i][j], reflect.Indirect(val).FieldByName(indexs[i][j]))
		}
		for k := 0; k < len(db.indexs[tableName][pk]); k++ {
			if db.indexs[tableName][pk][k] == id {
				db.indexs[tableName][pk][k] = db.indexs[tableName][pk][len(db.indexs[tableName][pk])-1]
				db.indexs[tableName][pk] = db.indexs[tableName][pk][:len(db.indexs[tableName][pk])-1]
			}
		}
		sortIndex(tableName, pk)
	}
	//log.Printf("index is %+v", db.indexs[tableName])
}
