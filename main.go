package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"ninedb/engine"
	"ninedb/models"
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

	u1 := models.User{GID: 0, UID: 1, Version: 0, TCC: 10000}
	engine.CreateTable(&u1)
	engine.Insert(&u1)
	engine.Insert(&u1)

	u2 := models.User{GID: 0, UID: 2, Version: 0, TCC: 10000}
	engine.Insert(&u2)

	m1 := models.TchMachine{ID: 1, GID: 0, UID: 1, Version: 0}
	engine.CreateTable(&m1)
	engine.Insert(&m1)

	for i := 3; i > 0; i-- {
		m := models.TchMachine{ID: i, GID: 0, UID: 1, Version: 0}
		engine.Insert(&m)
	}

	for i := 6; i > 3; i-- {
		m := models.TchMachine{ID: i, GID: 0, UID: 1, Version: 0}
		engine.Insert(&m)
	}

	m1.ID = 5
	engine.Delete(&m1)

	//peformance
	cnt := 10000
	start := time.Now().Unix()
	for i := 0; i < cnt; i++ {
		m := models.TchMachine{ID: i, GID: 0, UID: i % 1000, Version: 0}
		engine.Insert(&m)
	}
	end := time.Now().Unix()
	log.Printf("insert %d records in %d second", cnt, end-start)

	start = time.Now().Unix()
	for i := 0; i < cnt; i++ {
		m := models.TchMachine{ID: i, GID: 0, UID: i % 1000, Version: 1}
		engine.Update(&m)
	}
	end = time.Now().Unix()
	log.Printf("update %d records in %d second", cnt, end-start)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	log.Println(<-sigChan)

}
