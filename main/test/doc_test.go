package test

import (
	"testing"
)

/*
test case:
docset mydoc '{"$incr":{"name":1}}'
docset mydoc '{"$incr":{"version":1}}'
docset mydoc '{"$incr":{"version":"ok"}}'

docset mydoc '{"$rpush":"ok"}'
docset mydoc '{"$rpush":[1,2]}'
docset mydoc '{"$rpush":{"ok":1}}'
docset mydoc '{"$rpush":{"ok":[1, 3]}}'
docset mydoc '{"$rpush":{"name":["a.jpg"]}}'

docset user:100422:profile '{"sex.a":3}'
*/

func TestDoc(t *testing.T) {
	conn, err := NewRedisConn(host)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// clean
	if _, err := conn.Do("DEL", "mydoc"); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		`{"name":"latermoon"}`,
		`{"$incr":{"version":1}}`,
		`{"$incr":{"version":2}}`,
		`{"$rpush":{"photos":["a.jpg", "b.jpg"]}}`,
		`{"$rpush":{"photos":["c.jpg"]}}`,
		`{"setting.mute":{"start":22, "end":7, "status":true}}`,
		`{"setting.mute.start":23}`,
		`{"$del":["setting.mute.start"]}`}

	for _, line := range lines {
		if reply, err := conn.Do("DOCSET", "mydoc", line); err != nil {
			t.Fatal(err)
		} else if reply.(string) != "OK" {
			t.Error("bad reply")
		}
	}

	badcommands := []string{
		`{"name":"latermoon"`,  // bad json format
		`{"$incr":{"name":1}}`, // `name` is not int
	}

	for _, line := range badcommands {
		if _, err := conn.Do("DOCSET", "mydoc", line); err == nil { // must be error
			t.Error("no error is error")
		}
	}

	if reply, err := conn.Do("DOCGET", "mydoc", `version`); err != nil {
		t.Fatal(err)
	} else if string(reply.([]byte)) != `{"version":3}` {
		t.Error("bad reply")
	}
	if reply, err := conn.Do("DOCGET", "mydoc"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(reply.([]byte)))
	}
}
