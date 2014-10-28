package protocol

import (
	"bufio"
	"net"
)

type RedisHandler interface {
	// SessionOpened(session Session)
	// SessionClosed(session Session, err error)
	On(session Session, cmd Command) (Reply, error)
	ExceptionCaught(err error)
}

func ListenAndServe(addr string, handler RedisHandler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			handler.ExceptionCaught(err)
			continue
		}
		go handleConnection(handler, Session(conn))
	}
	return nil
}

func handleConnection(handler RedisHandler, session Session) {
	reader := &SessionReader{bufio.NewReader(session)}
	for {
		cmd, err := reader.ReadCommand()
		if err != nil {
			session.Close()
			break
		}

		reply, err := handler.On(session, cmd)
		if err != nil {
			session.Close()
			break
		}

		// if reply == nil {

		// }

		if _, err := session.Write(reply.Bytes()); err != nil {
			session.Close()
			break
		}
	}
}
