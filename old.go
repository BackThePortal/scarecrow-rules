package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Hello World")

	ctx := context.Background()
	svc, svcErr := docs.NewService(ctx, option.WithScopes(docs.DocumentsReadonlyScope))
	if svcErr != nil {
		log.Fatal(svcErr)
		return
	}

	do, err := svc.Documents.Get(os.Getenv("DOCUMENT_ID")).Do()
	if err != nil {
		return
	}

	fmt.Println(do.Body.Content)
}
