package store

import (
	//"encoding/json"
	"log"
	"ninedb/engine"
	"time"
)

const (
	WORKCNT = 1
)

func work() {
	for i := 0; i < WORKCNT; i++ {
		go func() {
			for {
				trx := engine.GetTrx()
				switch trx.Cmd {
				case "INSERT", "UPDATE":
					log.Printf("insert/update record %d", trx.ID)
					//todo save to db
					resp := &engine.Response{Code: "OK", TableName: trx.TableName, ID: trx.ID, SavedVersion: trx.Version, SavedStamp: time.Now().Unix()}
					engine.PutResp(resp)
				case "DELETE":
					log.Printf("delete record %d", trx.ID)
					resp := &engine.Response{Code: "OK", TableName: trx.TableName, ID: trx.ID, SavedVersion: trx.Version, SavedStamp: time.Now().Unix()}
					engine.PutResp(resp)
				}
			}
		}()

	}
}

func init() {
	work()
}
