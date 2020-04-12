package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	//"github.com/Azure/azure-sdk-for-go/services/preview/cognitiveservices/v3.0-preview/computervision"
        "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Declare global so don't have to pass it to all of the tasks.
var computerVisionContext context.Context
var database, _ = sql.Open("sqlite3", "./azure.db")

func main() {
//	imageURL := "https://commons.swinburne.edu.au/file/cd53e247-3e39-458e-8582-9fa0a2a2e120/1/cor-duncan_to_green_1920.jpg"
	//        imageURL := "https://rosetta.slv.vic.gov.au/delivery/DeliveryManagerServlet?dps_func=stream&dps_pid=FL16406745"
	imageURL := "https://a57.foxnews.com/static.foxnews.com/foxnews.com/content/uploads/2020/03/931/524/Ellen-DeGeneres-Jennifer-Aniston-Getty.jpg"
	//        imageURL :=  "https://rosetta.slv.vic.gov.au/delivery/DeliveryManagerServlet?dps_func=stream&dps_pid=FL18983698"
	//        imageURL := "https://rosetta.slv.vic.gov.au/delivery/DeliveryManagerServlet?dps_func=stream&dps_pid=FL18980978"
	//        imageURL := "https://rosetta.slv.vic.gov.au/delivery/DeliveryManagerServlet?dps_func=stream&dps_pid=FL16464085"
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS item (id INTEGER PRIMARY KEY, url TEXT)")
	statement.Exec()
	statement_entity, _ := database.Prepare("CREATE TABLE IF NOT EXISTS item_text (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_description (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_category (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_tag (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_object (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT, x TEXT, y TEXT, width TEXT, height TEXT, score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_brand (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT, x TEXT, y TEXT, width TEXT, height TEXT, score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_face (id INTEGER PRIMARY KEY, item_id INTEGER, gender TEXT, age TEXT, left TEXT, top TEXT, width TEXT, height TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_color (id INTEGER PRIMARY KEY, item_id INTEGER, black_and_white TEXT, accent_color TEXT, dominant_color_background TEXT, dominant_color_foreground TEXT, dominant_colors TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_adult (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_racy (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_celebrity (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_landmark (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT,score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_text_entity (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT, length INTERER, offset INTEGER, type TEXT, sub_type TEXT, score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_text_language (id INTEGER PRIMARY KEY, item_id INTEGER, value TEXT, score TEXT)")
	statement_entity.Exec()
	statement_entity, _ = database.Prepare("CREATE TABLE IF NOT EXISTS item_text_sentiment (id INTEGER PRIMARY KEY, item_id INTEGER, score TEXT)")
	statement_entity.Exec()

	stmt, err := database.Prepare("select id, url from item where url = ? limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var item_id int64
	//var last_insert_id int
	var url string
	err = stmt.QueryRow(imageURL).Scan(&item_id, &url)
	if item_id < 1 {

		statement, _ = database.Prepare("INSERT INTO item (url) VALUES (?)")
		result, err := statement.Exec(imageURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Last Insert ID")
		item_id, err = result.LastInsertId()
		fmt.Println(item_id)
	}
	fmt.Println("Item ID")
	fmt.Println(item_id)
	fmt.Println("Item URL")
	fmt.Println(imageURL)

	/*
	 * Configure the Computer Vision client
	 * Set environment variables for COMPUTER_VISION_SUBSCRIPTION_KEY and COMPUTER_VISION_ENDPOINT,
	 * then restart your command shell or your IDE for changes to take effect.
	 */
	computerVisionKey := os.Getenv("COMPUTER_VISION_SUBSCRIPTION_KEY")

	if computerVisionKey == "" {
		log.Fatal("\n\nPlease set a COMPUTER_VISION_SUBSCRIPTION_KEY environment variable.\n" +
			"**You may need to restart your shell or IDE after it's set.**\n")
	}

	endpointURL := os.Getenv("COMPUTER_VISION_ENDPOINT")
	if endpointURL == "" {
		log.Fatal("\n\nPlease set a COMPUTER_VISION_ENDPOINT environment variable.\n" +
			"**You may need to restart your shell or IDE after it's set.**")
	}

	computerVisionClient := computervision.New(endpointURL)
	computerVisionClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(computerVisionKey)

	computerVisionContext = context.Background()
	/*
	 * END - Configure the Computer Vision client
	 */
	/*printedImageURL := "https://i.imgur.com/I9r02n7.png"
	 */
	//	printedImageURL := "https://s3-ap-southeast-2.amazonaws.com/awm-media/collection/PR82/193.023/large/4164690.JPG"
	//        printedImageURL := "https://i.imgur.com/6n0uxk9.png" /*SLV avoca*/
	//        printedImageURL := "https://i.imgur.com/i41tezf.jpg" /*SLV eureka*/
	//printedImageURL := "https://i.imgur.com/YkqQZfB.png" /*George Swinburne*/
	//        printedImageURL := "https://i.imgur.com/XkJUPRL.png" /*SLV eureka*/
	//printedImageURL := "https://commons.swinburne.edu.au/file/cd53e247-3e39-458e-8582-9fa0a2a2e120/1/cor-duncan_to_green_1920.jpg"
	///* SWIn letteer
	// Analyze text in an image, remote
	BatchReadFileRemoteImage(computerVisionClient, imageURL, item_id)

	// Analyze features of an image, remote
	DescribeRemoteImage(computerVisionClient, imageURL, item_id)
	CategorizeRemoteImage(computerVisionClient, imageURL, item_id)
	TagRemoteImage(computerVisionClient, imageURL, item_id)
	DetectFacesRemoteImage(computerVisionClient, imageURL, item_id)
	DetectObjectsRemoteImage(computerVisionClient, imageURL, item_id)
	DetectBrandsRemoteImage(computerVisionClient, imageURL, item_id)
	DetectAdultOrRacyContentRemoteImage(computerVisionClient, imageURL, item_id)
	DetectColorSchemeRemoteImage(computerVisionClient, imageURL, item_id)
	DetectDomainSpecificContentRemoteImage(computerVisionClient, imageURL, item_id)
        
        //text_analytic(item_id) 
        ExtractEntities(item_id)
        DetectLanguage(item_id)
        SentimentAnalysis(item_id)


}

func DescribeRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DESCRIBE IMAGE - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	var caption_value string
	var caption_confidence float64
	remoteImage.URL = &remoteImageURL

	maxNumberDescriptionCandidates := new(int32)
	*maxNumberDescriptionCandidates = 1

	remoteImageDescription, err := client.DescribeImage(
		computerVisionContext,
		remoteImage,
		maxNumberDescriptionCandidates,
		"") // language
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Captions from remote image: ")
	if len(*remoteImageDescription.Captions) == 0 {
		fmt.Println("No captions detected.")
	} else {
		for _, caption := range *remoteImageDescription.Captions {
			fmt.Printf("'%v' with confidence %.2f%%\n", *caption.Text, *caption.Confidence*100)
			caption_value = *caption.Text
			caption_confidence = *caption.Confidence
		}
	}
	fmt.Println()

	statement, _ := database.Prepare("INSERT INTO item_description (item_id, value, score) VALUES (?, ?, ?)")
	result, err := statement.Exec(item_id, caption_value, caption_confidence)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Entity - Last Insert ID")
	iid, err := result.LastInsertId()
	fmt.Println(iid)

}
func CategorizeRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("CATEGORIZE IMAGE - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesCategories}
	imageAnalysis, err := client.AnalyzeImage(
		computerVisionContext,
		remoteImage,
		features,
		[]computervision.Details{},
		"")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Categories from remote image: ")
	if len(*imageAnalysis.Categories) == 0 {
		fmt.Println("No categories detected.")
	} else {
		for _, category := range *imageAnalysis.Categories {
			fmt.Printf("'%v' with confidence %.2f%%\n", *category.Name, *category.Score*100)

			statement, _ := database.Prepare("INSERT INTO item_category (item_id, value, score) VALUES (?, ?, ?)")
			result, err := statement.Exec(item_id, *category.Name, *category.Score)
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

func TagRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("TAG IMAGE - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	remoteImageTags, err := client.TagImage(
		computerVisionContext,
		remoteImage,
		"")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tags in the remote image: ")
	if len(*remoteImageTags.Tags) == 0 {
		fmt.Println("No tags detected.")
	} else {
		for _, tag := range *remoteImageTags.Tags {
			fmt.Printf("'%v' with confidence %.2f%%\n", *tag.Name, *tag.Confidence*100)

			statement, _ := database.Prepare("INSERT INTO item_tag (item_id, value, score) VALUES (?, ?, ?)")
			result, err := statement.Exec(item_id, *tag.Name, *tag.Confidence)
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

func DetectObjectsRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT OBJECTS - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	imageAnalysis, err := client.DetectObjects(
		computerVisionContext,
		remoteImage,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Detecting objects in remote image: ")
	if len(*imageAnalysis.Objects) == 0 {
		fmt.Println("No objects detected.")
	} else {
		// Print the objects found with confidence level and bounding box locations.
		for _, object := range *imageAnalysis.Objects {
			fmt.Printf("'%v' with confidence %.2f%% at location (%v, %v), (%v, %v)\n",
				*object.Object, *object.Confidence*100,
				*object.Rectangle.X, *object.Rectangle.X+*object.Rectangle.W,
				*object.Rectangle.Y, *object.Rectangle.Y+*object.Rectangle.H)

			statement, _ := database.Prepare("INSERT INTO item_object (item_id, value, x, y, width, height, score) VALUES (?, ?, ?, ?, ?, ?, ?)")
			result, err := statement.Exec(item_id, *object.Object, *object.Rectangle.X, *object.Rectangle.Y, *object.Rectangle.W, *object.Rectangle.H, *object.Confidence)
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

func DetectBrandsRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT BRANDS - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	// Define the kinds of features you want returned.
	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesBrands}

	imageAnalysis, err := client.AnalyzeImage(
		computerVisionContext,
		remoteImage,
		features,
		[]computervision.Details{},
		"en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Detecting brands in remote image: ")
	if len(*imageAnalysis.Brands) == 0 {
		fmt.Println("No brands detected.")
	} else {
		// Get bounding box around the brand and confidence level it's correctly identified.
		for _, brand := range *imageAnalysis.Brands {
			fmt.Printf("'%v' with confidence %.2f%% at location (%v, %v), (%v, %v)\n",
				*brand.Name, *brand.Confidence*100,
				*brand.Rectangle.X, *brand.Rectangle.X+*brand.Rectangle.W,
				*brand.Rectangle.Y, *brand.Rectangle.Y+*brand.Rectangle.H)

			statement, _ := database.Prepare("INSERT INTO item_brand (item_id, value, x, y, width, height, score) VALUES (?, ?, ?, ?, ?, ?, ?)")
			result, err := statement.Exec(item_id, *brand.Name, *brand.Rectangle.X, *brand.Rectangle.Y, *brand.Rectangle.W, *brand.Rectangle.H, *brand.Confidence)
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

func DetectFacesRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT FACES - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	// Define the features you want returned with the API call.
	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesFaces}
	imageAnalysis, err := client.AnalyzeImage(
		computerVisionContext,
		remoteImage,
		features,
		[]computervision.Details{},
		"")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Detecting faces in a remote image ...")
	if len(*imageAnalysis.Faces) == 0 {
		fmt.Println("No faces detected.")
	} else {
		// Print the bounding box locations of the found faces.
		for _, face := range *imageAnalysis.Faces {
			fmt.Printf("'%v' of age %v at location (%v, %v), (%v, %v)\n",
				face.Gender, *face.Age,
				*face.FaceRectangle.Left, *face.FaceRectangle.Top,
				*face.FaceRectangle.Left+*face.FaceRectangle.Width,
				*face.FaceRectangle.Top+*face.FaceRectangle.Height)

			statement, _ := database.Prepare("INSERT INTO item_face (item_id, gender, age, left, top, width, height) VALUES (?, ?, ?, ?, ?, ?, ?)")
			result, err := statement.Exec(item_id, face.Gender, *face.Age, *face.FaceRectangle.Left, *face.FaceRectangle.Top, *face.FaceRectangle.Width, *face.FaceRectangle.Height)
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

func DetectAdultOrRacyContentRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT ADULT OR RACY CONTENT - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	// Define the features you want returned from the API call.
	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesAdult}
	imageAnalysis, err := client.AnalyzeImage(
		computerVisionContext,
		remoteImage,
		features,
		[]computervision.Details{},
		"") // language, English is default
	if err != nil {
		log.Fatal(err)
	}

	// Print whether or not there is questionable content.
	// Confidence levels: low means content is OK, high means it's not.
	fmt.Println("Analyzing remote image for adult or racy content: ")
	fmt.Printf("Is adult content: %v with confidence %.2f%%\n", *imageAnalysis.Adult.IsAdultContent, *imageAnalysis.Adult.AdultScore*100)

	statement, _ := database.Prepare("INSERT INTO item_adult (item_id, value, score) VALUES (?, ?, ?)")
	result, err := statement.Exec(item_id, *imageAnalysis.Adult.IsAdultContent, *imageAnalysis.Adult.AdultScore)
	fmt.Println("Entity - Last Insert ID")
	iid, err := result.LastInsertId()
	fmt.Println(iid)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Has racy content: %v with confidence %.2f%%\n", *imageAnalysis.Adult.IsRacyContent, *imageAnalysis.Adult.RacyScore*100)

	statement, _ = database.Prepare("INSERT INTO item_racy (item_id, value, score) VALUES (?, ?, ?)")
	result, err = statement.Exec(item_id, *imageAnalysis.Adult.IsRacyContent, *imageAnalysis.Adult.RacyScore)
	fmt.Println("Entity - Last Insert ID")
	iid, err = result.LastInsertId()
	fmt.Println(iid)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
}

func DetectColorSchemeRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT COLOR SCHEME - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	// Define the features you'd like returned with the result.
	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesColor}
	imageAnalysis, err := client.AnalyzeImage(
		computerVisionContext,
		remoteImage,
		features,
		[]computervision.Details{},
		"") // language, English is default
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Color scheme of the remote image: ")
	fmt.Printf("Is black and white: %v\n", *imageAnalysis.Color.IsBWImg)
	fmt.Printf("Accent color: 0x%v\n", *imageAnalysis.Color.AccentColor)
	fmt.Printf("Dominant background color: %v\n", *imageAnalysis.Color.DominantColorBackground)
	fmt.Printf("Dominant foreground color: %v\n", *imageAnalysis.Color.DominantColorForeground)
	fmt.Printf("Dominant colors: %v\n", strings.Join(*imageAnalysis.Color.DominantColors, ", "))
	fmt.Println()

	statement, _ := database.Prepare("INSERT INTO item_color (item_id, black_and_white, accent_color, dominant_color_background, dominant_color_foreground, dominant_colors) VALUES (?, ?, ?, ?, ?, ?)")
	result, err := statement.Exec(item_id, *imageAnalysis.Color.IsBWImg, *imageAnalysis.Color.AccentColor, *imageAnalysis.Color.DominantColorBackground, *imageAnalysis.Color.DominantColorForeground, strings.Join(*imageAnalysis.Color.DominantColors, ", "))
	fmt.Println("Entity - Last Insert ID")
	iid, err := result.LastInsertId()
	fmt.Println(iid)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DetectDomainSpecificContentRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT DOMAIN-SPECIFIC CONTENT - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	fmt.Println("Detecting domain-specific content in the local image ...")

	// Check if there are any celebrities in the image.
	celebrities, err := client.AnalyzeImageByDomain(
		computerVisionContext,
		"celebrities",
		remoteImage,
		"") // language, English is default
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nCelebrities: ")

	// Marshal the output from AnalyzeImageByDomain into JSON.
	data, err := json.MarshalIndent(celebrities.Result, "", "\t")
      fmt.Println(string(data))

	// Define structs for which to unmarshal the JSON.
	type Celebrities struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	}

	type CelebrityResult struct {
		Celebrities []Celebrities `json:"celebrities"`
	}

	var celebrityResult CelebrityResult

	// Unmarshal the data.
	err = json.Unmarshal(data, &celebrityResult)
	if err != nil {
		log.Fatal(err)
	}

	//	Check if any celebrities detected.
	if len(celebrityResult.Celebrities) == 0 {
		fmt.Println("No celebrities detected.")
	} else {
		for _, celebrity := range celebrityResult.Celebrities {
			fmt.Printf("name: %v\n", celebrity.Name)
			fmt.Printf("confidence: %.2f%%\n", celebrity.Confidence)

			statement, _ := database.Prepare("INSERT INTO item_celebrity (item_id, value, score ) VALUES (?, ?, ?)")
			result, err := statement.Exec(item_id, celebrity.Name, celebrity.Confidence)
			fmt.Println("Entity - Last Insert ID")
			iid, err := result.LastInsertId()
			fmt.Println(iid)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	fmt.Println("\nLandmarks: ")

	// Check if there are any landmarks in the image.
	landmarks, err := client.AnalyzeImageByDomain(
		computerVisionContext,
		"landmarks",
		remoteImage,
		"")
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the output from AnalyzeImageByDomain into JSON.
	data, err = json.MarshalIndent(landmarks.Result, "", "\t")

	// Define structs for which to unmarshal the JSON.
	type Landmarks struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	}

	type LandmarkResult struct {
		Landmarks []Landmarks `json:"landmarks"`
	}

	var landmarkResult LandmarkResult

	// Unmarshal the data.
	err = json.Unmarshal(data, &landmarkResult)
	if err != nil {
		log.Fatal(err)
	}

	// Check if any celebrities detected.
	if len(landmarkResult.Landmarks) == 0 {
		fmt.Println("No landmarks detected.")
	} else {
		for _, landmark := range landmarkResult.Landmarks {
			fmt.Printf("name: %v\n", landmark.Name)

			statement, _ := database.Prepare("INSERT INTO item_landmark (item_id, value, score ) VALUES (?, ?, ?)")
			result, err := statement.Exec(item_id, landmark.Name, landmark.Confidence)
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

func BatchReadFileRemoteImage(client computervision.BaseClient, remoteImageURL string, item_id int64) {
	fmt.Println("-----------------------------------------")
	fmt.Println("BATCH READ FILE - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	var text_value string
	remoteImage.URL = &remoteImageURL

	// The response contains a field called "Operation-Location",
	// which is a URL with an ID that you'll use for GetReadOperationResult to access OCR results.
	textHeaders, err := client.BatchReadFile(computerVisionContext, remoteImage)
	if err != nil {
		log.Fatal(err)
	}

	// Use ExtractHeader from the autorest library to get the Operation-Location URL
	operationLocation := autorest.ExtractHeaderValue("Operation-Location", textHeaders.Response)

	numberOfCharsInOperationId := 36
	operationId := string(operationLocation[len(operationLocation)-numberOfCharsInOperationId : len(operationLocation)])
	readOperationResult, err := client.GetReadOperationResult(computerVisionContext, operationId)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the operation to complete.
	i := 0
	maxRetries := 10

	fmt.Println("Recognizing text in a remote image with the batch Read API ...")
	for readOperationResult.Status != computervision.Failed &&
		readOperationResult.Status != computervision.Succeeded {
		if i >= maxRetries {
			break
		}
		i++

		fmt.Printf("Server status: %v, waiting %v seconds...\n", readOperationResult.Status, i)
		time.Sleep(1 * time.Second)

		readOperationResult, err = client.GetReadOperationResult(computerVisionContext, operationId)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Display the results.
	fmt.Println()
	for _, recResult := range *(readOperationResult.RecognitionResults) {
		for _, line := range *recResult.Lines {
			fmt.Println(*line.Text)
			text_value += *line.Text
		}
	}
	statement, _ := database.Prepare("INSERT INTO item_text (item_id, value, score) VALUES (?, ?, ?)")
	result, err := statement.Exec(item_id, text_value, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Entity - Last Insert ID")
	iid, err := result.LastInsertId()
	fmt.Println(iid)
}
