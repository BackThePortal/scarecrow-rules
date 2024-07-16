package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kardianos/osext"
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
func getClient(config *oauth2.Config, interact bool, folderPath string) *http.Client {
	tokFile := folderPath + "/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		if interact {
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
		} else {
			log.Println("No token found. Manual action required.")
			os.Exit(3)
		}

	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
	fmt.Printf("M\n%v\n", authURL)

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
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	docId := flag.String("doc", "", "id of the document")
	interact := flag.Bool("interact", true, "enable user interaction")
	testOnly := flag.Bool("test-only", false, "don't retrieve document content, just check that it's readable")
	flag.Parse()

	ctx := context.Background()
	b, err := os.ReadFile(folderPath + "/credentials.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		os.Exit(2)

	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		os.Exit(1)

	}
	client := getClient(config, *interact, folderPath)

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Docs client: %v", err)
		os.Exit(4)

	}

	if *docId == "" {
		os.Exit(0)
	}

	//docId := os.Getenv("DOCUMENT_ID")
	doc, err := srv.Documents.Get(*docId).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from document: %v", err)
		os.Exit(5)
		return
	}
	if !*testOnly {
		body := readStructuralElements(doc.Body.Content)
		fmt.Printf(body)
	}
	os.Exit(0)
}
