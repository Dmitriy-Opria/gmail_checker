package quickstart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("../client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	srv, err := gmail.New(GetClient(config))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}
