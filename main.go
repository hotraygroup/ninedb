package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
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

	u1 := models.User{Base: models.Base{ID: 1}, GID: 0, TCC: 10000}
	engine.CreateTable(&u1)
	engine.Insert(&u1)
	engine.Insert(&u1)

	u2 := models.User{Base: models.Base{ID: 2}, GID: 0, TCC: 10000}
	engine.Insert(&u2)

	m1 := models.TchMachine{}
	engine.CreateTable(&m1)

	//peformance
	cnt := 1000000
	start := time.Now().Unix()

	//插入10台矿机
	for i := 0; i < 10; i++ {
		m := models.TchMachine{Base: models.Base{ID: i}, GID: 0, UID: i}
		//log.Printf("m:+%v", m)
		engine.Insert(&m)
	}

	/*
		for i := 0; i < cnt; i++ {
			m := models.TchMachine{Base: models.Base{ID: i}, GID: 0, UID: i % 10}
			//log.Printf("m:+%v", m)
			engine.Insert(&m)
		}
	*/
	end := time.Now().Unix()
	log.Printf("insert %d records in %d second", 10, end-start)

	//更新cnt次
	start = time.Now().Unix()
	for i := 0; i < cnt; i++ {
		m := models.TchMachine{Base: models.Base{ID: i % 10}, GID: 0, UID: i % 10}
		engine.Update(&m)
	}
	end = time.Now().Unix()
	log.Printf("update %d records in %d second", cnt, end-start)

}
