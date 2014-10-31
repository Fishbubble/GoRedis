package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Redis Command
type Command [][]byte

func NewCommand(args ...[]byte) Command {
	return Command(args)
}

func (c Command) Len() int {
	return len(c)
}

func (c Command) ArgAtIndex(i int) (arg []byte, err error) {
	if i >= c.Len() {
		return nil, errors.New(fmt.Sprintf("index out of range %d/%d", i, c.Len()))
	} else {
		return c[i], nil
	}
}

func (c Command) StringAtIndex(i int) (string, error) {
	if arg, err := c.ArgAtIndex(i); err == nil {
		return string(arg), nil
	} else {
		return "", err
	}
}

func (c Command) IntAtIndex(i int) (n int, err error) {
	if f, err := c.FloatAtIndex(i); err == nil {
		return int(f), nil
	} else {
		return 0, nil
	}
}

func (c Command) Int64AtIndex(i int) (int64, error) {
	if f, err := c.FloatAtIndex(i); err == nil {
		return int64(f), nil
	} else {
		return 0, err
	}
}

func (c Command) FloatAtIndex(i int) (float64, error) {
	if arg, err := c.StringAtIndex(i); err != nil {
		return 0, err
	} else if n, err := strconv.ParseFloat(arg, 64); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func (c Command) Bytes() []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte('*')
	argCount := c.Len()
	buf.WriteString(itoa(argCount)) //<number of arguments>
	buf.WriteString(CRLF)
	for i := 0; i < argCount; i++ {
		buf.WriteByte('$')
		buf.WriteString(itoa(len(c[i]))) //<number of bytes of argument i>
		buf.WriteString(CRLF)
		buf.Write(c[i]) //<argument data>
		buf.WriteString(CRLF)
	}
	return buf.Bytes()
}

func (c Command) String() string {
	arr := make([]string, len(c))
	for i := range c {
		arr[i] = string(c[i])
	}
	b, _ := json.Marshal(arr)
	return string(b)
}
