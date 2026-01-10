package main

import (
	"fmt"
	"os"
)

type TwitchClient struct {
	ClientId     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	UserID       string
}

// Load configuration from environment
func NewTwitchClient() (*TwitchClient, error) {
	client := &TwitchClient{
		ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITHC_CLIENT_SECRET"),
		AccessToken:  os.Getenv("TWITCH_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("TWITCH_REFRESH_TOKEN"),
		UserID:       os.Getenv("TWITCH_USER_ID"),
	}

	if client.ClientId == "" {
		return nil, fmt.Errorf("TWITCH_CLIENT_ID env variable not found")
	}
	if client.AccessToken == "" {
		return nil, fmt.Errorf("TWITCH_ACCESS_TOKEN env variable not found")
	}
	return client, nil
}

type SubscriptionRequest struct {
	Type      string `json:"type"`
	Version   string `json:"version"`
	Condition string `json:"condition"`
	Transport string `json:"transport"`
}

func eventSubSubscriptions() {
}
