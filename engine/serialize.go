package engine

import (
	"log"
)

type Transaction struct {
	Cmd  string
	Data interface{}
}

type Response struct {
	Code   string
	Result string
	Data   interface{}
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
				id := res.Data.(Row).GetID()
				log.Printf("id %d done", id)

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
