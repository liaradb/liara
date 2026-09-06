package listener

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

func Listen(port int, server *grpc.Server) error {
	listener, err := getListener(port)
	if err != nil {
		return err
	}

	go listen(listener, server)

	return nil
}

func getListener(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func listen(listener net.Listener, server *grpc.Server) {
	slog.Info(fmt.Sprintf("listening at %v", listener.Addr()))
	if err := server.Serve(listener); err != nil {
		slog.Error("failed to serve",
			"error", err)
		panic(err)
	}
}
