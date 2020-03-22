package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    //"log"
    "net/http"
    "os"
    "strings"
    "time"
    "database/sql"
    "strconv"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
  //database
	database, _ := 
		sql.Open("sqlite3", "./bogo.db")
	statement, _ := 
		database.Prepare("CREATE TABLE IF NOT EXISTS item_raw  (id INTEGER PRIMARY KEY, ocr_text TEXT, text_analytics TEXT)")
    statement.Exec()
	statement, _ = 
		database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
    statement.Exec("Rob", "Gronkowski")
	rows, _ := 
		database.Query("SELECT id, firstname, lastname FROM people")
    var id int
    var firstname string
    var lastname string
    for rows.Next() {
        rows.Scan(&id, &firstname, &lastname)
        fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
    }
  //azure
    var subscriptionKey string = os.Getenv("TEXT_ANALYTICS_KEY")
    var endpoint string = os.Getenv("TEXT_ANALYTICS_ENDPOINT")
    
    //const uriPath = "/text/analytics/v2.1/entities"
    const uriPath = "/text/analytics/v3.0-preview.1/entities/recognition/general"

    var uri = endpoint + uriPath
    var text =`S.S. Shropshire
At Sea
4th March 19120

Mr F.W. Green.
Dear Sir,

We are off to Freemantle for another 40,000 Boxes of apples after loading 110,00 at Hobart. We did not call at Melbourne again. We will land 6,000 Boxes of Butter at Freemantle from Sydney.

Will you please remember me to Mr Waters & to my fiends at My College and with best wishes to yourself I remain

Mr F.W. Green.
Yours faithfully
John Duncan`

input := string(text)
s := strings.Replace(input, "\n"," ",-1)

    data := []map[string]string{
        {"id": "1", "language": "en", "text": s},
    }

    documents, err := json.Marshal(&data)
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
    if err != nil {
        fmt.Printf("Error reading response body: %v\n", err)
        return
    }

    var f interface{}
    json.Unmarshal(body, &f)

    jsonFormatted, err := json.MarshalIndent(f, "", "  ")
    if err != nil {
        fmt.Printf("Error producing JSON: %v\n", err)
        return
    }
    fmt.Println(string(jsonFormatted))
}
