package store

import (
	"fmt"
	"github.com/erwaen/pneumo/minivalkey"
)

type StorageService struct {
	valkeyClient *minivalkey.Client
}

func CreateStore() *StorageService {
	valkeyClient, err := minivalkey.CreateClient("valkey:6379")
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


func (ss *StorageService) StorePneumo(fullUrl, pneumo string) error {
	_, err := ss.valkeyClient.Set(pneumo, fullUrl)
	if err != nil {
		return fmt.Errorf("error storing pneumo: %w", err)
	}

	return nil
}

func (ss *StorageService) RetrieveFromPneumo(pneumo string) (string, error) {
	fullUrl, err := ss.valkeyClient.Get(pneumo)
	if err != nil {
		return "", fmt.Errorf("error retrieving full pneumo: %w", err)
	}

	if fullUrl == "" {
		return "", fmt.Errorf("pneumo not found: %s", pneumo)
	}

	return fullUrl, nil
}


