package main

import (
    "context"
   // "encoding/json"
    "fmt"
    //"log"
    "os"
    "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/to"
)

func GetTextAnalyticsClient() textanalytics.BaseClient {
    var key string = os.Getenv("TEXT_ANALYTICS_KEY")
    var endpoint string = os.Getenv("TEXT_ANALYTICS_ENDPOINT")

    textAnalyticsClient := textanalytics.New(endpoint)
    textAnalyticsClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(key)

    return textAnalyticsClient
}

func ExtractEntities() {
    textAnalyticsClient := GetTextAnalyticsClient()
    ctx := context.Background()
    inputDocuments := []textanalytics.MultiLanguageInput{
        {
            Language: to.StringPtr("en"),
            ID:       to.StringPtr("0"),
            Text:     to.StringPtr("Microsoft was founded by Bill Gates and Paul Allen on April 4, 1975, to develop and sell BASIC interpreters for the Altair 8800."),
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

