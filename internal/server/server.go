package server

import (
	"context"
	"net"
)

func Start(ctx context.Context) (int, <-chan net.Conn, <-chan error, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, nil, nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port

	connectionErrors := make(chan error)
	connections := make(chan net.Conn)

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	go func() {
		defer close(connectionErrors)
		defer close(connections)

		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					connectionErrors <- err
					continue
				}
			}
			connections <- conn
		}
	}()

	return port, connections, connectionErrors, nil
}
