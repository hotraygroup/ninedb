package engine

import (
	"log"
)

type Transaction struct {
	Cmd       string
	TableName string
	ID        int
	Version   uint64
}

type Response struct {
	Code         string
	TableName    string
	ID           int
	SavedVersion uint64
	SavedStamp   int64
}

func Prepare() {

}

func putTrx(trx *Transaction) {
	select {
	case ReqChan <- trx:
	default:
		log.Printf("reqChain is full")
	}
}

func GetTrx() *Transaction {
	return <-ReqChan
}

func PutResp(res *Response) {
	select {
	case RespChan <- res:
	default:
		log.Printf("respChain is full")
	}
}

func getResp() *Response {
	return <-RespChan
}

func work() {
	go func() {
		for {
			res := getResp()
			switch res.Code {
			case "OK":
				log.Printf("table %s's record %d, version %d done at %d", res.TableName, res.ID, res.SavedVersion, res.SavedStamp)
			case "SKIP":
				log.Printf("table %s's record %d, version %d skipped at %d", res.TableName, res.ID, res.SavedVersion, res.SavedStamp)
			}
		}
	}()
}

const (
	WORKCHANSIZE = 100000
)

var ReqChan chan *Transaction
var RespChan chan *Response

func init() {
	ReqChan = make(chan *Transaction, WORKCHANSIZE)
	RespChan = make(chan *Response, WORKCHANSIZE)
	work()
}
