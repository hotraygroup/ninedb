package main

import (
	"flag"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	_ "net/http/pprof"
	"ninedb/controller"
	"ninedb/engine"
	"ninedb/models"
	_ "ninedb/store"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	CpuProfile  = flag.String("cpu-profile", "", "write cpu profile to file")
	HeapProfile = flag.String("heap-profile", "", "write heap profile to file")
)

func main() {
	log.Printf("main")
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if *CpuProfile != "" {
		file, err := os.Create(*CpuProfile)
		if err != nil {
			log.Panicln(err)
		}
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	if *HeapProfile != "" {
		file, err := os.Create(*HeapProfile)
		if err != nil {
			log.Panicln(err)
		}
		defer pprof.WriteHeapProfile(file)
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//for test
	sample()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	log.Println(<-sigChan)

}

func sample() {

	//////////////////建表/////////////////////////////////////////
	u1 := models.User{UID: 1, GID: 0, TCC: decimal.New(99, 2), ETH: decimal.New(199, 2), NASH: decimal.New(299, 2), Worker: map[int]bool{2: true}}
	engine.CreateTable(&u1)
	engine.Insert(&u1)
	engine.Insert(&u1)

	u2 := models.User{UID: 2, GID: 0, TCC: decimal.New(99, 2), ETH: decimal.New(199, 2), NASH: decimal.New(299, 2), Worker: map[int]bool{1: true}}
	engine.Insert(&u2)

	m1 := models.TchMachine{}
	engine.CreateTable(&m1)
	//////////////////建表/////////////////////////////////////////

	///////////////////插入/////////////////////////////////////////
	ucnt := 10
	//插入ucnt个用户
	for i := 0; i <= ucnt; i++ {
		u := models.User{UID: i, GID: 0, TCC: decimal.New(99, 2), ETH: decimal.New(199, 2), NASH: decimal.New(299, 2), Worker: map[int]bool{1: true}}
		engine.Insert(&u)
	}
	mcnt := 10
	start := time.Now().Unix()

	//插入mcnt台矿机
	for i := 0; i < mcnt; i++ {
		m := models.TchMachine{ID: i, GID: 0, UID: i % ucnt}
		//log.Printf("m:+%v", m)
		engine.Insert(&m, "load")
	}
	end := time.Now().Unix()
	log.Printf("insert %d records in %d second", mcnt, end-start)
	///////////////////插入/////////////////////////////////////////

	///////////////////删除/////////////////////////////////////////
	u10 := models.User{UID: 10}
	engine.Delete(&u10)
	///////////////////删除/////////////////////////////////////////

	///////////////////更新/////////////////////////////////////////
	start = time.Now().Unix()
	for i := 0; i < mcnt; i++ {
		m := models.TchMachine{ID: i % 10, GID: 0, UID: i % ucnt}
		engine.Update(&m)
	}
	end = time.Now().Unix()
	log.Printf("update %d records in %d second", mcnt, end-start)
	///////////////////更新/////////////////////////////////////////

	///////////////////转账/////////////////////////////////////////////
	//engine.UpdateFunc((controller.Transfer(nil, nil, nil)).(engine.CallBack))
	log.Printf("before transfer: user1: %+v, user2: %+v", engine.Get(&u1), engine.Get(&u2))

	controller.Transfer(1, 2, "TCC", decimal.New(1, 1))
	controller.Transfer(1, 2, "TCC", decimal.New(1, 100))

	controller.Transfer(1, 2, "ETH", decimal.New(1, -1))

	controller.Transfer(1, 2, "NASH", decimal.New(1, 2))

	log.Printf("after transfer: user1: %+v, user2: %+v", engine.Get(&u1), engine.Get(&u2))

	log.Printf("before u1 is %+v", u1)
	engine.UpdateField(&u1, "Desc", "REPLACE", "ssss")
	engine.UpdateField(&u1, "Worker", "REPLACE", map[int]bool{3: true})

	engine.UpdateField(&u1, "I1", "ZERO", 0)
	engine.UpdateField(&u1, "I1", "REPLACE", 1000)
	engine.UpdateField(&u1, "I1", "INC", 100)
	engine.UpdateField(&u1, "I1", "DESC", 10)

	log.Printf("after u1 is %+v", u1)

	///////////////////转账/////////////////////////////////////////////

}
