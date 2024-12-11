package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hiennq12/my-money/caculator_data"
	"github.com/hiennq12/my-money/noti"
	"github.com/hiennq12/my-money/struct_modal"
	"github.com/robfig/cron/v3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"time"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	//allProcess()

	c := cron.New(
	//cron.WithSeconds(), // Cho phép lập lịch theo giây thi bieu thuc cron co 6 dau *
	)

	// // */15 * * * * 15p 1 laanf,  @hourly
	jobID, err := c.AddFunc("*/15 * * * *", allProcess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Job ID: %v\n", jobID)

	c.Start()

	// Đợi một thời gian để xem kết quả
	time.Sleep(time.Hour)

	// Dừng cron an toàn
	ctx := c.Stop()
	<-ctx.Done()
}

func allProcess() {
	ctx := context.Background()
	b, err := os.ReadFile("/Users/hiennguyen/Documents/learn/my-money/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1rVQtA77ILhvj03bCANNu5mRP4vgJXKRyQUpuScUuppI/edit?gid=0#gid=0
	spreadsheetId := "1rVQtA77ILhvj03bCANNu5mRP4vgJXKRyQUpuScUuppI"
	readRange := "T12/2024!A1:Z40"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	moneySpending, err := caculator_data.MoneySpending(&struct_modal.DataRows{
		ValueRange: resp,
	})

	if err != nil {
		log.Fatalf("Error when calculator money spend in day: %v", err.Error())
	}

	configTele, message := noti.PrepareData(moneySpending)
	if configTele != nil {
		err = noti.SendTelegramMessage(configTele, message)
	}
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	fmt.Println("Message sent successfully!", message)
}
