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

func ExtractEntities(item_id int64, timestamp int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? order by timestamp limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)
	//azure

	stmt2, err := database.Prepare("select code from item_text_language where item_id = ? order by timestamp DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt2.Close()
	var lang_code string
	//var last_insert_id int
	err = stmt2.QueryRow(item_id).Scan(&lang_code)
	const uriPath = "/text/analytics/v2.1/entities"

	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.MultiLanguageInput{
		{
			Language: to.StringPtr(lang_code),
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

				statement, _ := database.Prepare("INSERT INTO item_text_entity (item_id, timestamp,value, length, offset, type, sub_type, score ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
				result, err := statement.Exec(item_id,timestamp, *entity.Name, *match.Length, *match.Offset, *entity.Type, entity.SubType, *match.EntityTypeScore)
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
func DetectLanguage(item_id int64, timestamp int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? order by timestamp DESC limit 1")
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
			fmt.Printf("%s %s %f,", *language.Name, *language.Iso6391Name,*language.Score)

			statement, _ := database.Prepare("INSERT INTO item_text_language (item_id, timestamp, value, code, score ) VALUES (?, ?, ?, ?, ?)")
			result, err := statement.Exec(item_id, timestamp, *language.Name, *language.Iso6391Name, *language.Score)
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

func SentimentAnalysis(item_id int64, timestamp int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? order by timestamp DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)
	//azure
	stmt2, err := database.Prepare("select code from item_text_language where item_id = ? order by timestamp DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt2.Close()
	var lang_code string
	//var last_insert_id int
	err = stmt2.QueryRow(item_id).Scan(&lang_code)

	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.MultiLanguageInput{
		{
			Language: to.StringPtr(lang_code),
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

		statement, _ := database.Prepare("INSERT INTO item_text_sentiment (item_id,timestamp, score ) VALUES (?, ?, ?)")
		result, err := statement.Exec(item_id, timestamp, *document.Score)
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
func ExtractKeyPhrases(item_id int64, timestamp int64) {

	stmt, err := database.Prepare("select value from item_text where item_id = ? order by timestamp DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_text string
	//var last_insert_id int
	err = stmt.QueryRow(item_id).Scan(&item_text)

	stmt2, err := database.Prepare("select code from item_text_language where item_id = ? order by timestamp DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt2.Close()
	var lang_code string
	//var last_insert_id int
	err = stmt2.QueryRow(item_id).Scan(&lang_code)
	//azure
	textAnalyticsClient := GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := []textanalytics.MultiLanguageInput{
		{
			Language: to.StringPtr(lang_code),
			ID:       to.StringPtr("0"),
			Text:     to.StringPtr(item_text),
		},
	}

	batchInput := textanalytics.MultiLanguageBatchInput{Documents: &inputDocuments}
	result, _ := textAnalyticsClient.KeyPhrases(ctx, to.BoolPtr(false), &batchInput)

	// Printing extracted key phrases results
	for _, document := range *result.Documents {
		fmt.Printf("Document ID: %s\n", *document.ID)
		fmt.Printf("Document Language: %s\n", lang_code)
		fmt.Printf("\tExtracted Key Phrases:\n")
		for _, keyPhrase := range *document.KeyPhrases {
			fmt.Printf("\t\t%s\n", keyPhrase)

                        statement, _ := database.Prepare("INSERT INTO item_text_key_phrase (item_id,timestamp, value ) VALUES (?, ?, ?)")
                        result, err := statement.Exec(item_id, timestamp, keyPhrase)
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
func Get(item_id int64) {

  fmt.Println("Print Items: Key Phrases")
	stmt, err := database.Prepare("WITH RECURSIVE  x as (select value, timestamp, max(timestamp) as max from item_text_key_phrase where item_id = ?) select item_id, value from item_text_key_phrase WHERE timestamp = (select max from x);")
        result, err := stmt.Query(item_id)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var value string

        for result.Next() {
          err := result.Scan(&item_id, &value)
          if err != nil {
            log.Fatal(err)
          }
          fmt.Println(item_id, value)
        }
      }
