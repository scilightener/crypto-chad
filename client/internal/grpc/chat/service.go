package chat

import (
	"context"
	"crypto-chad-client/internal/domain"
	"crypto-chad-client/internal/grpc/certs"
	pb "crypto-chad-client/internal/grpc/chat/generated"
	"crypto-chad-lib/rsa"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

//go:generate protoc --go_out=./generated --go_opt=paths=source_relative --go-grpc_out=./generated --go-grpc_opt=paths=source_relative chat.proto
type Client struct {
	u           *domain.User
	client      pb.ChatClient
	certsClient *certs.Client
}

func NewClient(u *domain.User, client pb.ChatClient, certsClient *certs.Client) *Client {
	return &Client{
		u:           u,
		client:      client,
		certsClient: certsClient,
	}
}

func (c *Client) SendMessage(receiver string, message string) {
	cert := c.certsClient.RetrieveCert(receiver)
	if cert == nil {
		return
	}
	encrypted := rsa.Encrypt([]byte(message), cert.E, cert.N)
	req := &pb.Message{Sender: c.u.Name, Receiver: receiver, Message: encrypted}
	_, err := c.client.SendMessage(context.Background(), req)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func (c *Client) ReceiveMessages(ctx context.Context) {
	req := &pb.User{Username: c.u.Name}
	stream, err := c.client.ReceiveMessages(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}
	defer stream.CloseSend()
	for {
		select {
		case <-ctx.Done():
			break
		default:
			resp, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v\n", err)
				return
			}
			if resp.GetReceiver() != c.u.Name {
				continue
			}

			message := rsa.Decrypt(resp.GetMessage(), c.u.Keys.PrivateKey.D, c.u.Keys.PrivateKey.N)
			printMessage(resp.GetSender(), string(message))
		}
	}
}

func (c *Client) ShowAllUsers() []string {
	users, err := c.client.ActiveUsers(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return users.GetUsernames()
}

func printMessage(sender, message string) {
	fmt.Printf("[%s]: %s\n", sender, message)
}
