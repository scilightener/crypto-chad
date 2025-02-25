package app

import (
	"crypto-chad-server/internal/handlers/grpc/certs"
	"crypto-chad-server/internal/handlers/grpc/chat"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

const defaultAddr = "0.0.0.0:9876"

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := NewApp(logger)
	startApp(app)
}

func startApp(app *App) {
	host, ok := os.LookupEnv("HOST_ADDRESS")
	if !ok {
		host = defaultAddr
	}
	tcp, err := net.Listen("tcp", host)
	if err != nil {
		app.logger.Error("Failed to listen: %v", slog.String("err", err.Error()))
		panic(err)
	}

	server := grpc.NewServer()
	chat.Register(server, app)
	certs.Register(server, app)

	app.logger.Info("app started", slog.String("addr", tcp.Addr().String()))
	if err = server.Serve(tcp); err != nil {
		app.logger.Error("failed to serve", slog.String("err", err.Error()))
		panic(err)
	}
}
