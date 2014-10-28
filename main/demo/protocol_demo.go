package main

import (
	. "GoRedis/protocol"
	"fmt"
)

type RedisServer struct {
	RedisHandler
}

func (r *RedisServer) On(session Session, cmd Command) (Reply, error) {
	return StatusReply("OK"), nil
}

func main() {
	startServer()
}

func startServer() {
	fmt.Println("start ...")
	s := &RedisServer{}
	err := ListenAndServe(":3001", s)
	if err != nil {
		panic(err)
	}
}

func cmd() {
	cmd := Command{[]byte("SET"), []byte("name"), []byte("latermoon")}
	fmt.Println(cmd)

	ok := StatusReply("OK")
	fmt.Println(ok, ok.Bytes())

	er := ErrorReply("404 NotFound")
	fmt.Println(er, er.Bytes())

	size := IntegerReply(100)
	fmt.Println(size, size.Bytes())

	bulk := BulkReply([]byte("latermoon"))
	fmt.Println(string(bulk), bulk.Bytes())

	bulks := MultiBulkReply{10, "name", []byte("latermoon")}
	fmt.Println(bulks, bulks.Bytes())

}
