package main

import (
        "encoding/base64"
        "flag"
        "fmt"
        "io/ioutil"
        "os"
        "path/filepath"
        "golang.org/x/net/context"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/vision/v1"
)

func main() {
        flag.Usage = func() {
                fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
        }

        flag.Parse()

        args := flag.Args()
        if len(args) == 0 {
                flag.Usage()
                os.Exit(1)
        }
        
        if err := run(args[0]); err != nil {
                // Comes here if run() returns an error
                fmt.Fprintf(os.Stderr, "%s\n", err.Error())
                os.Exit(1)
        }
       
}

func run(file string) error {
        ctx := context.Background()

        // Authenticate to generate a vision service
        client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
        if err != nil {
                return err
        }

        service, err := vision.New(client)
        if err != nil {
                return err
        }
        // We now have a Vision API service with which we can make API calls.

        b, err := ioutil.ReadFile(file)
        if err != nil {
                return err
        }

        // Construct a text request, encoding the image in base64.
        req := &vision.AnnotateImageRequest{
                // Apply image which is encoded by base64
                Image: &vision.Image{
                        Content: base64.StdEncoding.EncodeToString(b),
                },
                // Apply features to indicate what type of image detection
                Features: []*vision.Feature{
                        {
                                Type: "TEXT_DETECTION", 
                        },
                },
        }

        batch := &vision.BatchAnnotateImagesRequest{
                Requests: []*vision.AnnotateImageRequest{req},
        }
        
        res, err := service.Images.Annotate(batch).Do()
        if err != nil {
                return err
        }
        // A POST request has been made

        // Parse annotations from responses
        if annotations := res.Responses[0].TextAnnotations; len(annotations) > 0 {
                text := annotations[0].Description
                fmt.Printf("Found text: %s\n", text)
                return nil
        }
        fmt.Printf("Not found text in: %s\n", file)

        return nil
}
