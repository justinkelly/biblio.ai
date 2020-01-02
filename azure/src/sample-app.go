package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	"log"
        "os"
	"strings"
	"time"
)

// Declare global so don't have to pass it to all of the tasks.
var computerVisionContext context.Context

func main() {
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
	printedImageURL := "https://i.imgur.com/I9r02n7.png"
	landmarkImageURL := printedImageURL
	brandsImageURL := printedImageURL
	facesImageURL := printedImageURL
	objectsImageURL := printedImageURL
	adultRacyImageURL := printedImageURL
	detectTypeImageURL := printedImageURL
	// Analyze text in an image, remote
	BatchReadFileRemoteImage(computerVisionClient, printedImageURL)

	// Analyze features of an image, remote
	DescribeRemoteImage(computerVisionClient, landmarkImageURL)
	CategorizeRemoteImage(computerVisionClient, landmarkImageURL)
	TagRemoteImage(computerVisionClient, landmarkImageURL)
	DetectFacesRemoteImage(computerVisionClient, facesImageURL)
	DetectObjectsRemoteImage(computerVisionClient, objectsImageURL)
	DetectBrandsRemoteImage(computerVisionClient, brandsImageURL)
	DetectAdultOrRacyContentRemoteImage(computerVisionClient, adultRacyImageURL)
	DetectColorSchemeRemoteImage(computerVisionClient, brandsImageURL)
	DetectDomainSpecificContentRemoteImage(computerVisionClient, landmarkImageURL)
	DetectImageTypesRemoteImage(computerVisionClient, detectTypeImageURL)

}

func DescribeRemoteImage(client computervision.BaseClient, remoteImageURL string) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DESCRIBE IMAGE - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
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
		}
	}
	fmt.Println()
}

func CategorizeRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
		}
	}
	fmt.Println()
}

func TagRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
		}
	}
	fmt.Println()
}

func DetectObjectsRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
		}
	}
	fmt.Println()
}

func DetectBrandsRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
		}
	}
	fmt.Println()
}

func DetectFacesRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
		}
	}
	fmt.Println()
}

func DetectAdultOrRacyContentRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
	fmt.Printf("Has racy content: %v with confidence %.2f%%\n", *imageAnalysis.Adult.IsRacyContent, *imageAnalysis.Adult.RacyScore*100)
	fmt.Println()
}

func DetectColorSchemeRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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
}

func DetectDomainSpecificContentRemoteImage(client computervision.BaseClient, remoteImageURL string) {
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

	// Define structs for which to unmarshal the JSON.
	type Celebrities struct {
		Name string `json:"name"`
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
		Name string `json:"name"`
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
		}
	}
	fmt.Println()
}

func DetectImageTypesRemoteImage(client computervision.BaseClient, remoteImageURL string) {
	fmt.Println("-----------------------------------------")
	fmt.Println("DETECT IMAGE TYPES - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesImageType}

	imageAnalysis, err := client.AnalyzeImage(
		computerVisionContext,
		remoteImage,
		features,
		[]computervision.Details{},
		"")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image type of remote image:")

	fmt.Println("\nClip art type: ")
	switch *imageAnalysis.ImageType.ClipArtType {
	case 0:
		fmt.Println("Image is not clip art.")
	case 1:
		fmt.Println("Image is ambiguously clip art.")
	case 2:
		fmt.Println("Image is normal clip art.")
	case 3:
		fmt.Println("Image is good clip art.")
	}

	fmt.Println("\nLine drawing type: ")
	if *imageAnalysis.ImageType.LineDrawingType == 1 {
		fmt.Println("Image is a line drawing.")
	} else {
		fmt.Println("Image is not a line drawing.")
	}
	fmt.Println()
}

func BatchReadFileRemoteImage(client computervision.BaseClient, remoteImageURL string) {
	fmt.Println("-----------------------------------------")
	fmt.Println("BATCH READ FILE - remote")
	fmt.Println()
	var remoteImage computervision.ImageURL
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
		}
	}
}
