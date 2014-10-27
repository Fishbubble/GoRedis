package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Session struct {
	conn net.Conn
	rw   *bufio.Reader
}

func NewSession(conn net.Conn) *Session {
	s := &Session{
		conn: conn,
		rw:   bufio.NewReader(conn),
	}
	return s
}

func (s *Session) Read(b []byte) (n int, err error) {
	return s.rw.Read(b)
}

// 读取一行
func (s *Session) readLine() (line []byte, err error) {
	line, err = s.rw.ReadSlice(LF)
	if err == bufio.ErrBufferFull {
		return nil, errors.New("line too long")
	}
	if err != nil {
		return
	}
	i := len(line) - 2
	if i < 0 || line[i] != CR {
		err = errors.New("bad line terminator:" + string(line))
	}
	return line[:i], nil
}

// 读取字符串，遇到CRLF换行为止
func (s *Session) readString() (str string, err error) {
	var line []byte
	if line, err = s.readLine(); err != nil {
		return
	}
	str = string(line)
	return
}

func (s *Session) readInt() (i int, err error) {
	var line string
	if line, err = s.readString(); err != nil {
		return
	}
	i, err = strconv.Atoi(line)
	return
}

func (s *Session) readInt64() (i int64, err error) {
	var line string
	if line, err = s.readString(); err != nil {
		return
	}
	i, err = strconv.ParseInt(line, 10, 64)
	return
}

// 验证并跳过指定的字节，用于开始符和结束符的判断
func (s *Session) skipByte(c byte) (err error) {
	var tmp byte
	tmp, err = s.rw.ReadByte()
	if err != nil {
		return
	}
	if tmp != c {
		err = errors.New(fmt.Sprintf("Illegal Byte [%d] != [%d]", tmp, c))
	}
	return
}

func (s *Session) skipBytes(bs []byte) (err error) {
	for _, c := range bs {
		err = s.skipByte(c)
		if err != nil {
			break
		}
	}
	return
}

func (s *Session) ReadWith() interface{} {
	return nil
}

func ReadCommand(s *Session) (Command, error) {
	return nil, nil
}

func ReadReply(r io.Reader) (Reply, error) {
	return nil, nil
}
