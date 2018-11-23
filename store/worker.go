package store

import (
	//"encoding/json"
	"log"
	"ninedb/engine"
)

func work() {
	go func() {
		for {
			trx := engine.GetTrx()
			switch trx.Cmd {
			case "INSERT":
				id := trx.Data.(engine.Row).GetID()
				log.Printf("insert record %d", id)
				resp := &engine.Response{Code: "OK", Data: trx.Data}
				engine.PutResp(resp)
			case "UPDATE":
				id := trx.Data.(engine.Row).GetID()
				log.Printf("update record %d", id)
				resp := &engine.Response{Code: "OK", Data: trx.Data}
				engine.PutResp(resp)
			case "DELETE":
				id := trx.Data.(engine.Row).GetID()
				log.Printf("delete record %d", id)
				resp := &engine.Response{Code: "OK", Data: trx.Data}
				engine.PutResp(resp)
			}
		}
	}()

}

func init() {
	work()
}
