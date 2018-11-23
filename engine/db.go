package engine

import (
	"sync"
)

const (
	ROWSIZE = 100000
)

type DB struct {
	sync.RWMutex
	tables map[string][]interface{}    //one table one store array
	rows   map[string]map[int]int      //one table one map: pk id -> index of store array
	ids    map[string][]int            //one table one map: pk ids. 用于全查找
	indexs map[string]map[string][]int //one talbe one map: index key ->  pk id array， 用于条件查找
	chans  map[string]chan int         //one table one alloc chan
	locks  map[string]sync.RWMutex     //one table one lock
}

var db DB

func init() {
	db.tables = make(map[string][]interface{})
	db.rows = make(map[string]map[int]int)
	db.ids = make(map[string][]int)
	db.chans = make(map[string]chan int)
	db.locks = make(map[string]sync.RWMutex)
	db.indexs = make(map[string]map[string][]int)
}
