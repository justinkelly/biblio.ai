package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetTextAnalyticsClient() textanalytics.BaseClient {
	var subscriptionKey string = os.Getenv("TEXT_ANALYTICS_KEY")
	var endpoint string = os.Getenv("TEXT_ANALYTICS_ENDPOINT")
	textAnalyticsClient := textanalytics.New(endpoint)
	textAnalyticsClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscriptionKey)

	return textAnalyticsClient
}

func ExtractEntities(item_id int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)
	//azure

	const uriPath = "/text/analytics/v2.1/entities"

	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.MultiLanguageInput{
		{
			Language: to.StringPtr("en"),
			ID:       to.StringPtr("0"),
			Text:     to.StringPtr(item_text),
		},
	}

	batchInput := textanalytics.MultiLanguageBatchInput{Documents: &inputDocuments}
	result, _ := textAnalyticsClient.Entities(ctx, to.BoolPtr(false), &batchInput)

	// Printing extracted entities results
	for _, document := range *result.Documents {
		fmt.Printf("Document ID: %s\n", *document.ID)
		fmt.Printf("\tExtracted Entities:\n")
		for _, entity := range *document.Entities {
			fmt.Printf("\t\tName: %s\tType: %s", *entity.Name, *entity.Type)
			if entity.SubType != nil {
				fmt.Printf("\tSub-Type: %s\n", *entity.SubType)
			}
			fmt.Println()
			for _, match := range *entity.Matches {
				fmt.Printf("\t\t\tOffset: %v\tLength: %v\tScore: %f\n", *match.Offset, *match.Length, *match.EntityTypeScore)

				statement, _ := database.Prepare("INSERT INTO item_text_entity (item_id, value, length, offset, type, sub_type, score ) VALUES (?, ?, ?, ?, ?, ?, ?)")
				result, err := statement.Exec(item_id, *entity.Name, *match.Length, *match.Offset, *entity.Type, entity.SubType, *match.EntityTypeScore)
				fmt.Println("Entity - Last Insert ID")
				iid, err := result.LastInsertId()
				fmt.Println(iid)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		fmt.Println()
	}

	// Printing document errors
	fmt.Println("Document Errors")
	for _, err := range *result.Errors {
		fmt.Printf("Document ID: %s Message : %s\n", *err.ID, *err.Message)
	}
}
func DetectLanguage(item_id int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)
	//azure
	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.LanguageInput{
		{
			ID:   to.StringPtr("0"),
			Text: to.StringPtr(item_text),
		},
	}

	batchInput := textanalytics.LanguageBatchInput{Documents: &inputDocuments}
	result, _ := textAnalyticsClient.DetectLanguage(ctx, to.BoolPtr(false), &batchInput)

	// Printing language detection results
	for _, document := range *result.Documents {
		fmt.Printf("Document ID: %s ", *document.ID)
		fmt.Printf("Detected Languages with Score: ")
		for _, language := range *document.DetectedLanguages {
			fmt.Printf("%s %f,", *language.Name, *language.Score)

			statement, _ := database.Prepare("INSERT INTO item_text_language (item_id, value, score ) VALUES (?, ?, ?)")
			result, err := statement.Exec(item_id, *language.Name, *language.Score)
			fmt.Println("Entity - Last Insert ID")
			iid, err := result.LastInsertId()
			fmt.Println(iid)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println()
	}

	// Printing document errors
	fmt.Println("Document Errors")
	for _, err := range *result.Errors {
		fmt.Printf("Document ID: %s Message : %s\n", *err.ID, *err.Message)
	}
}

func SentimentAnalysis(item_id int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)
	//azure
	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.MultiLanguageInput{
		{
			Language: to.StringPtr("en"),
			ID:       to.StringPtr("0"),
			Text:     to.StringPtr(item_text),
		},
	}

	batchInput := textanalytics.MultiLanguageBatchInput{Documents: &inputDocuments}
	result, _ := textAnalyticsClient.Sentiment(ctx, to.BoolPtr(false), &batchInput)
	var batchResult textanalytics.SentimentBatchResult
	jsonString, _ := json.Marshal(result)
	_ = json.Unmarshal(jsonString, &batchResult)

	// Printing sentiment results
	for _, document := range *batchResult.Documents {
		fmt.Printf("Document ID: %s ", *document.ID)
		fmt.Printf("Sentiment Score: %f\n", *document.Score)

		statement, _ := database.Prepare("INSERT INTO item_text_sentiment (item_id, score ) VALUES (?, ?)")
		result, err := statement.Exec(item_id, *document.Score)
		fmt.Println("Entity - Last Insert ID")
		iid, err := result.LastInsertId()
		fmt.Println(iid)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Printing document errors
	fmt.Println("Document Errors")
	for _, err := range *batchResult.Errors {
		fmt.Printf("Document ID: %s Message : %s\n", *err.ID, *err.Message)
	}
}
func ExtractKeyPhrases(item_id int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)
	//azure
	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.MultiLanguageInput{
		{
			Language: to.StringPtr("en"),
			ID:       to.StringPtr("0"),
			Text:     to.StringPtr(item_text),
		},
	}

	batchInput := textanalytics.MultiLanguageBatchInput{Documents: &inputDocuments}
	result, _ := textAnalyticsClient.KeyPhrases(ctx, to.BoolPtr(false), &batchInput)

	// Printing extracted key phrases results
	for _, document := range *result.Documents {
		fmt.Printf("Document ID: %s\n", *document.ID)
		fmt.Printf("\tExtracted Key Phrases:\n")
		for _, keyPhrase := range *document.KeyPhrases {
			fmt.Printf("\t\t%s\n", keyPhrase)
		}
		fmt.Println()
	}

	// Printing document errors
	fmt.Println("Document Errors")
	for _, err := range *result.Errors {
		fmt.Printf("Document ID: %s Message : %s\n", *err.ID, *err.Message)
	}
}
