package chat

import (
	"context"
	pb "crypto-chad-server/internal/handlers/grpc/chat/generated"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler interface {
	ForgetUser(name string) bool
}

//go:generate protoc --go_out=./generated --go_opt=paths=source_relative --go-grpc_out=./generated --go-grpc_opt=paths=source_relative chat.proto
type server struct {
	pb.UnimplementedChatServer
	clients map[string]pb.Chat_ReceiveMessagesServer
	handler Handler
}

func (s *server) SendMessage(_ context.Context, req *pb.Message) (*pb.Message, error) {
	if client, ok := s.clients[req.GetReceiver()]; ok {
		err := client.Send(req)
		if err != nil {
			return nil, fmt.Errorf("unable to send message to %s", req.GetReceiver())
		}
		return &pb.Message{}, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (s *server) ReceiveMessages(req *pb.User, messagesServer pb.Chat_ReceiveMessagesServer) error {
	username := req.GetUsername()
	s.clients[username] = messagesServer
	defer func() {
		delete(s.clients, username)
		s.handler.ForgetUser(username)
	}()

	<-messagesServer.Context().Done()
	return nil
}

func (s *server) ActiveUsers(_ context.Context, _ *emptypb.Empty) (*pb.Users, error) {
	users := make([]string, 0, len(s.clients))
	for username := range s.clients {
		users = append(users, username)
	}
	resp := &pb.Users{Usernames: users}
	return resp, nil
}

func Register(s *grpc.Server, chatHandler Handler) {
	pb.RegisterChatServer(s, &server{clients: make(map[string]pb.Chat_ReceiveMessagesServer), handler: chatHandler})
}
