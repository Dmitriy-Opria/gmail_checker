package main

import (
	"fmt"
	"gmail_check/quickstart"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	pathToKeyFile = "client_secret.json"
	userID        = "demoni421@gmail.com"
)

func main() {
	for {
		oauthHttpClient, err := NewClient()
		if err != nil {
			fmt.Printf("Can`t make client error: %v\n", err)
			os.Exit(1)
		}
		gmailService, err := gmail.New(oauthHttpClient)
		if err != nil {
			fmt.Printf("Can`t get new gmail service")
			os.Exit(1)
		}

		fmt.Printf("Unreded messages(1): %v\n", getUnreadCount(gmailService))

		listUnreadMessages, err := gmailService.Users.Messages.List(userID).LabelIds("UNREAD").Do()
		if err != nil {
			fmt.Printf("Can`t get unreaded message list, error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Unreaded messages(2): %d\n", len(listUnreadMessages.Messages))

		for _, msg := range listUnreadMessages.Messages {
			makeRead := gmail.ModifyMessageRequest{
				RemoveLabelIds: []string{"UNREAD"},
			}
			_, err := gmailService.Users.Messages.Modify(userID, msg.Id, &makeRead).Do()
			if err != nil {
				fmt.Printf("Can`t modify email, error: %v\n", err)
			} else {
				fmt.Printf("Read msg by id: %s\n", msg.Id)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func getUnreadCount(gmailService *gmail.Service) (count int64) {
	label, err := gmailService.Users.Labels.Get(userID, "INBOX").Do()
	if err != nil {
		fmt.Printf("Can`t get labels, error: %v\n", err)
		return
	}
	return label.MessagesUnread
}

func NewClient() (client *http.Client, err error) {

	jsonKey, err := ioutil.ReadFile(pathToKeyFile)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(jsonKey, gmail.MailGoogleComScope) // gmail.MailGoogleComScope
	if err != nil {
		return nil, err
	}
	client = quickstart.GetClient(config)
	return
}
