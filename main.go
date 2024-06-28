package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"strings"
)

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
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
func tokenFromFile(file string) (*oauth2.Token, error) {
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
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func readParagraphElement(element docs.ParagraphElement) string {
	text := element.TextRun
	if text == nil {
		return ""
	}
	return text.Content
}

func readStructuralElements(elements []*docs.StructuralElement) string {
	text := ""
	for _, element := range elements {
		if element.Paragraph == nil {
			text += "\n"
			continue
		}
		if element.Paragraph.ParagraphStyle.NamedStyleType == "HEADING_1" {
			text += "# "
		}
		if element.Paragraph.ParagraphStyle.NamedStyleType == "HEADING_2" {
			text += "## "
		}
		if element.Paragraph.ParagraphStyle.NamedStyleType == "HEADING_3" {
			text += "### "
		}
		if element.Paragraph.Bullet != nil {
			text += "- "
		}
		for _, paragraphElement := range element.Paragraph.Elements {
			fragment := paragraphElement.TextRun.Content
			if paragraphElement.TextRun.Content != "" && paragraphElement.TextRun.Content != "\n" {
				if paragraphElement.TextRun.TextStyle.Bold {
					fragment = "**" + paragraphElement.TextRun.Content + "**"
				}
				if paragraphElement.TextRun.TextStyle.Italic {
					fragment = "_" + paragraphElement.TextRun.Content + "_"
				}
			}
			text += fragment

		}
	}
	return strings.ReplaceAll(text, string(rune(11)), "")
}

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	docId := os.Getenv("DOCUMENT_ID")
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	body := readStructuralElements(doc.Body.Content)
	fmt.Printf(body)
}
