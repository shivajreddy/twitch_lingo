package main

import (
	// "bytes"
	// "encoding/json"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// TODO 1: connect to twitch stream

// TODO 2: Print the stream

// TODO 3: Initiate translator

// TODO 4.1: each message into separate thread
// TODO 4.2: multithread convert messages
// TODO 4.3: collect messages, and print (order doesnt matter for now)

// GLOBAL VARIABLES
var (
	CLIENT_ID           string
	CLIENT_SECRET       string
	TWITCH_ACCESS_TOKEN string
)

func logFatalErr(err any) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnvFile() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	CLIENT_ID = os.Getenv("CLIENT_ID")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
}

func getTwitchAccessToken() {
	// Create form data
	FORM_DATA := url.Values{}
	FORM_DATA.Set("client_id", CLIENT_ID)
	FORM_DATA.Set("client_secret", CLIENT_SECRET)
	FORM_DATA.Set("grant_type", "client_credentials")

	// Make the request
	resp, err := http.Post(
		"https://id.twitch.tv/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(FORM_DATA.Encode()),
	)
	logFatalErr(err)
	defer resp.Body.Close()

	// Parse JSON response
	type AuthTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	var tokenResp AuthTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	logFatalErr(err)

	fmt.Println("access_token:", tokenResp.AccessToken)
	fmt.Println("expires_in:", tokenResp.ExpiresIn)
	fmt.Println("token_type:", tokenResp.TokenType)
	TWITCH_ACCESS_TOKEN = tokenResp.AccessToken
}

func getToken2() {
	URL := "https://id.twitch.tv/oauth2/authorize"

	req, err := http.NewRequest("GET", URL, nil)
	logFatalErr(err)

	// Add query params
	q := req.URL.Query()

	q.Add("client_id", CLIENT_ID)
	q.Add("client_secret", CLIENT_SECRET)
	q.Add("response_type", "token")
	q.Add("redirect_uri", "http://localhost:3000")
	// q.Add("scope", "")
	// q.Add("state", "")

	// final query looks like this
	req.URL.RawQuery = q.Encode()
	fmt.Println(q)

	// Send Request
	client := &http.Client{}
	resp, err := client.Do(req)
	logFatalErr(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func makeTwitchCall() {
	URL := "https://api.twitch.tv/helix/users"

	req, err := http.NewRequest("GET", URL, nil)
	logFatalErr(err)

	// Add query params
	q := req.URL.Query()
	q.Add("login", "twitchdev")
	req.URL.RawQuery = q.Encode()

	// Add Headers
	req.Header.Set("Authorization", "Bearer "+TWITCH_ACCESS_TOKEN)
	req.Header.Set("Client-Id", CLIENT_ID)

	// Send Request
	client := &http.Client{}
	resp, err := client.Do(req)
	logFatalErr(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	/*
		{
			"data":[
			{"id":"141981764",
				"login":"twitchdev",
				"display_name":"TwitchDev",
				"type":"",
				"broadcaster_type":"partner",
				"description":"Supporting third-party developers building Twitch integrations from chatbots to game integrations.",
				"profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png",
				"offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png",
				"view_count":0,
				"created_at":"2016-12-14T20:32:28Z"}
			]
		}
	*/

	// curl -X POST 'https://id.twitch.tv/oauth2/token' \
	// -H 'Content-Type: application/x-www-form-urlencoded' \
	// -d 'client_id=<your client id goes here>
	// &client_secret=<your client secret goes here>
	// &code=17038swieks1jh1hwcdr36hekyui
	// &grant_type=authorization_code
	// &redirect_uri=http://localhost:3000'

	type RespShape struct {
		AccessToken  string   `json:"access_token"`
		ExpiresIn    int      `json:"expires_in"`
		RefreshToken string   `json:"refresh_token"`
		Scope        []string `json:"scope"`
		TokenType    string   `json:"token_type"`
	}

}

func readTwitchChat() {
	channelId := "qsnake"

	server := "irc.chat.twitch.tv:6667"
	conn, err := net.Dial("tcp", server)
	logFatalErr(err)
	defer conn.Close()

	fmt.Println("Connected to Twitch IRC")

	// Authenticate
	fmt.Fprintf(conn, "PASS oauth:%s\r\n", TWITCH_ACCESS_TOKEN)
	fmt.Fprintf(conn, "NICK %s\r\n", "simplecolon")

	// Join Channel
	fmt.Fprintf(conn, "JOIN #%s\r\n", channelId)

	fmt.Printf("Joined channel: #%s\n", channelId)
	fmt.Println("Listening for messages...")

	// Read Messages
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			break
		}
		fmt.Println("message: ", line)
	}

}

func main() {
	fmt.Println("TWITCH-LINGO")

	loadEnvFile()

	// getToken2()

	// getTwitchAccessToken()
	// makeTwitchCall()

	// readTwitchChat()
}
