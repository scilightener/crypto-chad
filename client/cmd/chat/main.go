package main

import (
	"bufio"
	"context"
	"crypto-chad-client/internal/common"
	"crypto-chad-client/internal/config"
	"crypto-chad-client/internal/domain"
	"crypto-chad-client/internal/grpc/certs"
	pbcerts "crypto-chad-client/internal/grpc/certs/generated"
	"crypto-chad-client/internal/grpc/chat"
	pbchat "crypto-chad-client/internal/grpc/chat/generated"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
)

const defaultConfigPath = "./config.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <username> <config_path>\n"+
			"<username> is your username under which you will communicate with others in this chat\n"+
			"<config_path> is the path to the config with which the app should run [if it's not provided, it's assumed to be %s]",
			os.Args[0], defaultConfigPath)
		return
	}

	configPath := defaultConfigPath
	if len(os.Args) >= 3 {
		configPath = os.Args[2]
	}
	cfg := config.MustLoad(configPath)
	common.SetServerPubKey(cfg.Server.RSAKey.E, cfg.Server.RSAKey.N)

	// don't forget to set Server PubKey
	if common.ServerPubKey == nil {
		panic("server pubkey is not set")
	}

	conn, err := grpc.NewClient(cfg.Server.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	user := &domain.User{
		Name: os.Args[1],
		Keys: nil,
	}

	certClient := certs.NewClient(user, pbcerts.NewCertsClient(conn))
	certClient.IssueCert()
	chatClient := chat.NewClient(user, pbchat.NewChatClient(conn), certClient)
	interactWithUser(chatClient)
}
func interactWithUser(chatClient *chat.Client) {
	go chatClient.ReceiveMessages(context.Background())

	fmt.Println("Введите !quit чтобы выйти")
	fmt.Println("Введите !show чтобы показать всех активных пользователей")
	fmt.Println("Вводите сообщения в виде 'alice: hello', т.е. сначала к кому вы обращаетесь, " +
		"затем двоеточие, и само сообщение")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()

		if input == "!quit" {
			fmt.Println("Exiting program...")
			break
		}

		if input == "!show" {
			fmt.Println(chatClient.ShowAllUsers())
			continue
		}
		parts := splitMessage(input)

		if len(parts) != 2 {
			log.Println("Некорректный формат ввода")
			continue
		}

		chatClient.SendMessage(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
}

func splitMessage(input string) []string {
	found := false
	return strings.FieldsFunc(input, func(r rune) bool {
		if !found && r == ':' {
			found = true
			return true
		}

		return false
	})
}
