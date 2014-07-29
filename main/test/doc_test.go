package test

import (
	"testing"
)

/*
Bad Command:
docset mydoc '{"$incr":{"name":1}}'
docset mydoc '{"$incr":{"version":1}}'
docset mydoc '{"$incr":{"version":"ok"}}'

docset mydoc '{"$rpush":"ok"}'
docset mydoc '{"$rpush":[1,2]}'
docset mydoc '{"$rpush":{"ok":1}}'
docset mydoc '{"$rpush":{"ok":[1, 3]}}'
docset mydoc '{"$rpush":{"name":["a.jpg"]}}'
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

	if reply, err := conn.Do("DOCSET", "mydoc", `{"name":"latermoon"}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"$incr":{"name":1}}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"$incr":{"version":1}}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"$incr":{"version":2}}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCGET", "mydoc", `version`); err != nil {
		t.Fatal(err)
	} else if string(reply.([]byte)) != `{"version":3}` {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"$rpush":{"photos":["a.jpg", "b.jpg"]}}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"$rpush":{"photos":["c.jpg"]}}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"setting.mute":{"start":22, "end":7, "status":true}}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"setting.mute.start":23}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCSET", "mydoc", `{"$del":["setting.mute.start"]}`); err != nil {
		t.Fatal(err)
	} else if reply.(string) != "OK" {
		t.Error("bad reply")
	}

	if reply, err := conn.Do("DOCGET", "mydoc"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(reply.([]byte)))
	}
}
