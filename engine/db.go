package engine

import (
	"sync"
)

const (
	K = 1024
	M = 1024 * K
	G = 1024 * M
)

const (
	ROWSIZE    = M
	PRIMARYKEY = "pk"
)

type MetaInfo struct {
	Version      uint64
	UpdateStamp  int64
	SavedVersion uint64
	SavedStamp   int64
}

type DB struct {
	sync.RWMutex
	tables    map[string][]interface{}     //one table one store array
	rows      map[string]map[int]int       //one table one map: pk id -> index of store array
	metas     map[string]map[int]*MetaInfo //one table one map: pk id - > version
	indexs    map[string]map[string][]int  //one table one map: index key ->  pk id array， 用于条件查找
	sorting   map[string]map[string]int32  //one table one map: index key -> is sorting
	sortlocks map[string]*sync.Mutex       //one table one map: index key -> sorting tmux
	chans     map[string]chan int          //one table one alloc chan
	//stats   map[string]map[string][][]int //one table one map: key->(v, cnt). todo
	locks map[string]*sync.Mutex //one table one lock
}

var db DB

func init() {
	db.tables = make(map[string][]interface{})
	db.rows = make(map[string]map[int]int)
	db.metas = make(map[string]map[int]*MetaInfo)
	db.chans = make(map[string]chan int)
	db.locks = make(map[string]*sync.Mutex)
	db.indexs = make(map[string]map[string][]int)
	db.sorting = make(map[string]map[string]int32)
	db.sortlocks = make(map[string]*sync.Mutex)
}
