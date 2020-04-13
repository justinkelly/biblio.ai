package main

import (
  "context"
    "encoding/json"
    "fmt"
    "os"
    "log"
    /*
    "io/ioutil"
    "net/http"
    "strings"
    "time"
    "github.com/tidwall/gjson"
    "github.com/savaki/jq"
    */
     "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"
        "github.com/Azure/go-autorest/autorest"
        "github.com/Azure/go-autorest/autorest/to"


    _ "github.com/mattn/go-sqlite3"
)

/*
func text_analytic(item_id int64) {
  stmt, err := database.Prepare("select value from item_text where item_id = ? limit 1")
  if err != nil {
    log.Fatal(err)
  }
  defer stmt.Close()
  var item_text string
  //var last_insert_id int
  err = stmt.QueryRow(item_id).Scan(&item_text)
  //azure
  var subscriptionKey string = os.Getenv("TEXT_ANALYTICS_KEY")
  var endpoint string = os.Getenv("TEXT_ANALYTICS_ENDPOINT")

  const uriPath = "/text/analytics/v2.1/entities"

  var uri = endpoint + uriPath
  fmt.Println(item_text)
  s := item_text
  data := []map[string]string{
    {"id": "1", "language": "en", "text": s},
  }

  documents, err := json.Marshal(&data)
  fmt.Println("\nText Analytics: POST docu")
  fmt.Println(string(documents))
  if err != nil {
    fmt.Printf("Error marshaling data: %v\n", err)
    return
  }

  r := strings.NewReader("{\"documents\": " + string(documents) + "}")

  client := &http.Client{
    Timeout: time.Second * 20,
  }

  req, err := http.NewRequest("POST", uri, r)
  if err != nil {
    fmt.Printf("Error creating request: %v\n", err)
    return
  }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Ocp-Apim-Subscription-Key", subscriptionKey)

  resp, err := client.Do(req)
  if err != nil {
    fmt.Printf("Error on request: %v\n", err)
    return
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  fmt.Println("\nText Analytics: body")
  fmt.Println(body)
  fmt.Println("\nText Analytics: resp body")
  fmt.Println(resp.Body)

  if err != nil {
    fmt.Printf("Error reading response body: %v\n", err)
    return
  }
  fmt.Println("\nText Analytics: ")


  // Define structs for which to unmarshal the JSON.
  type Celebrities struct {
    Name       string  `json:"name"`
    Type float64 `json:"type"`
  }

  type CelebrityResult struct {
    Celebrities []Celebrities `json:"entities"`
  }

  var celebrityResult CelebrityResult

  // Unmarshal the data.
  err = json.Unmarshal(body, &celebrityResult)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("No text detected.")

  //	Check if any celebrities detected.
  if len(celebrityResult.Celebrities) == 0 {
    fmt.Println("No text detected.")
  } else {
    for _, celebrity := range celebrityResult.Celebrities {
      fmt.Printf("name: %v\n", celebrity.Name)
      fmt.Printf("type: %v\n", celebrity.Type)

      statement, _ := database.Prepare("INSERT INTO item_text_analytic (item_id, value, score ) VALUES (?, ?, ?)")
      result, err := statement.Exec(item_id, celebrity.Name, celebrity.Type)
      fmt.Println("Entity - Last Insert ID")
      iid, err := result.LastInsertId()
      fmt.Println(iid)
      if err != nil {
        fmt.Println(err)
        return
      }
    }
  }
  var f interface{}
  json.Unmarshal(body, &f)

  jsonFormatted, err := json.MarshalIndent(f, "", "  ")
  if err != nil {
    fmt.Printf("Error producing JSON: %v\n", err)
    return
  }
  fmt.Println()
  fmt.Println("JQ")
  fmt.Println(jsonFormatted)
  fmt.Println("JQ")
  fmt.Println(string(jsonFormatted))
  //op, _ := jq.Parse(".documents[0].entities")           // create an Op
  op, _ := jq.Parse(".documents[0]")           // create an Op
  //data2 := []byte(`{"hello":"world"}`)   // sample input
  data2 := jsonFormatted   // sample input
  value2, _ := op.Apply(data2)            // value == '"world"'
  fmt.Println("JQ")
  fmt.Println(string(value2))
  fmt.Println("JQ - end")


  println("gjson")
  result := gjson.Get(string(body), "documents")
  for _, name := range result.Array() {
    println("\n")
    println(name.String())
    println("\n")
  }
}
*/

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
            result, err := statement.Exec(item_id, *entity.Name , *match.Length, *match.Offset, *entity.Type, entity.SubType, *match.EntityTypeScore)
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
            result, err := statement.Exec(item_id, *language.Name , *language.Score)
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

