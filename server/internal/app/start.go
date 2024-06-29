package app

import (
	"crypto-chad-server/internal/handlers/grpc/certs"
	"crypto-chad-server/internal/handlers/grpc/chat"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := NewApp(logger)
	startApp(app)
}

func startApp(app *App) {
	tcp, err := net.Listen("tcp", "localhost:9876")
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
