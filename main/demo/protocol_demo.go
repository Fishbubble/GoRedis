package main

import (
	. "GoRedis/protocol"
	"fmt"
)

func main() {
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
