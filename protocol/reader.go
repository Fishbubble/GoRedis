package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type SessionReader struct {
	rw *bufio.Reader
}

func (s *SessionReader) ReadCommand() (Command, error) {
	// Read ( *<number of arguments> CR LF )
	var err error
	if err = s.skipByte('*'); err != nil {
		return nil, err // io.EOF
	}
	// number of arguments
	var argCount int
	if argCount, err = s.readInt(); err != nil {
		return nil, err
	}
	args := make([][]byte, argCount)
	for i := 0; i < argCount; i++ {
		// Read ( $<number of bytes of argument 1> CR LF )
		if err = s.skipByte('$'); err != nil {
			return nil, err
		}

		var argSize int
		if argSize, err = s.readInt(); err != nil {
			return nil, err
		}

		// Read ( <argument data> CR LF )
		args[i] = make([]byte, argSize)
		_, err = io.ReadFull(s.rw, args[i])
		if err != nil {
			return nil, err
		}

		if err = s.skipBytes([]byte{CR, LF}); err != nil {
			return nil, err
		}
	}
	return Command(args), nil
}

// 读取一行
func (s *SessionReader) readLine() (line []byte, err error) {
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
func (s *SessionReader) readString() (string, error) {
	if line, err := s.readLine(); err != nil {
		return "", err
	} else {
		return string(line), nil
	}
}

func (s *SessionReader) readInt() (i int, err error) {
	var line string
	if line, err = s.readString(); err != nil {
		return
	}
	i, err = strconv.Atoi(line)
	return
}

func (s *SessionReader) readInt64() (i int64, err error) {
	var line string
	if line, err = s.readString(); err != nil {
		return
	}
	i, err = strconv.ParseInt(line, 10, 64)
	return
}

// 验证并跳过指定的字节，用于开始符和结束符的判断
func (s *SessionReader) skipByte(c byte) (err error) {
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

func (s *SessionReader) skipBytes(bs []byte) (err error) {
	for _, c := range bs {
		err = s.skipByte(c)
		if err != nil {
			break
		}
	}
	return
}
