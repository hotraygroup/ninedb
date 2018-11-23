package engine

import (
	"sync"
)

const (
	ROWSIZE    = 100000
	PRIMARYKEY = "pk"
)

type Version struct {
	Version      int64
	SavedVersion uint64
}

type DB struct {
	sync.RWMutex
	tables    map[string][]interface{}    //one table one store array
	rows      map[string]map[int]int      //one table one map: pk id -> index of store array
	versions  map[string]map[int]*Version //one table one map: pk id - > version
	indexs    map[string]map[string][]int //one table one map: index key ->  pk id array， 用于条件查找
	sorting   map[string]map[string]int32 //one table one map: index key -> is sorting
	sortlocks map[string]*sync.Mutex      //one table one map: index key -> sorting tmux
	chans     map[string]chan int         //one table one alloc chan
	//stats   map[string]map[string][][]int //one table one map: key->(v, cnt). todo
	locks map[string]*sync.Mutex //one table one lock
}

var db DB

func init() {
	db.tables = make(map[string][]interface{})
	db.rows = make(map[string]map[int]int)
	db.versions = make(map[string]map[int]*Version)
	db.chans = make(map[string]chan int)
	db.locks = make(map[string]*sync.Mutex)
	db.indexs = make(map[string]map[string][]int)
	db.sorting = make(map[string]map[string]int32)
	db.sortlocks = make(map[string]*sync.Mutex)
}
