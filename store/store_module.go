package store

import (
	"fmt"
	"github.com/erwaen/pneumo/minivalkey"
)

type StorageService struct {
	valkeyClient *minivalkey.Client
}

func CreateStore() *StorageService {
	valkeyClient, err := minivalkey.CreateClient("localhost:6379")
	if err != nil {
		panic(fmt.Sprintf("Could not connect to valkey: %v", err))
	}

	// check health
	pong, err := valkeyClient.SendRespCommand("PING")
	if err != nil {
		panic(fmt.Sprintf("Error sending PING command: %v\n", err))
	}

	fmt.Printf("Connected to valkey storage successfully: pong = %s\n", pong)
	return &StorageService{valkeyClient: valkeyClient}
}

func (ss *StorageService) Close() error {
	return ss.valkeyClient.Close()
}

