package goredis_server

/*
实现类似于mongo的document base指令，
对一个key提供docuemnt存储，以及原子操作

user:100422:profile = {
	name: "latermoon", // string
	sex: 1 // int
	photos: ["a.jpg", "b.jpg", "c.jpg"], // array<string>
	setting: { // hash
		mute: {
			start: 23,
			end: 8
		}
	},
	is_vip: true, // bool
	version: 172 // int
}

// Update/Insert
docset(key, {"name":"latermoon"})
docset(key, {"$incr":{"version":1}})
docset(key, {"$rpush":{"photos":["a.jpg","b.jpg"]}})
docset(key, {"setting.mute":{"start":23, "end":8}})
docset(key, {"setting.mute.start":23, "setting.mute.end":8})
docset(key, {"$del":["name", "setting.mute.start"])
docset(key, {"$set":{"name":"latermoon", "sex":"M"}, "$inc":{"profile.version":1}})

// Get All
docget(key)
docget(key, "name,sex,photos,setting.mute,version")

*/

import (
	. "GoRedis/goredis"
	"encoding/json"
	"fmt"
	"strings"
)

/*
docset hi '{"name":"latermoon", "sex":"M", "version":10, "setting":{"start":23, "end":8}}'
docset hi '{"$inc":{"version":1}}'
docset hi '{"$del":["version", "setting.start"]}'
*/
func (server *GoRedisServer) OnDOCSET(cmd *Command) (reply *Reply) {
	key := cmd.StringAtIndex(1)
	// 传入的json字节
	jsonbytes, err := cmd.ArgAtIndex(2)
	if err != nil {
		return ErrorReply(err)
	}
	// 反序列化为map
	jsonObj := make(map[string]interface{})
	err = json.Unmarshal(jsonbytes, &jsonObj)
	if err != nil {
		return ErrorReply(fmt.Sprintf("bad json format: %s", err))
	}
	// 调用LevelDocument更新数据
	doc := server.levelRedis.GetDoc(key)
	err = doc.Set(jsonObj)
	if err != nil {
		return ErrorReply(err)
	}
	reply = StatusReply("OK")
	return
}

func (server *GoRedisServer) OnDOCGET(cmd *Command) (reply *Reply) {
	key := cmd.StringAtIndex(1)
	fields := strings.Split(cmd.StringAtIndex(2), ",")
	doc := server.levelRedis.GetDoc(key)
	result := doc.Get(fields...)
	if result == nil {
		return BulkReply(nil)
	}
	data, err := json.Marshal(result)
	if err != nil {
		return ErrorReply(err)
	}
	reply = BulkReply(data)
	return
}
