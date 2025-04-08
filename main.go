package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ahmedsat/middleware/handlers"
	"github.com/ahmedsat/middleware/helpers"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var local bool = false

func init() {

	// dataFile, err := os.Open("kobo-data.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer dataFile.Close()

	// fa := internals.FarmApplication{}
	// err = fa.Scan(dataFile)
	// if err != nil {
	// 	panic(err)
	// }
	// err = fa.Validate()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// os.Exit(0)
	// res, err := ERPRequest("GET", "/api/resource/User Permission", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.StatusCode)
	// f, err := os.Create("res.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// io.Copy(f, res.Body)

	// os.Exit(0)
	// errList := CreateFarmApplicationFromKoboData("in.json")
	// fmt.Println(errList)
	// if len(errList) > 0 {
	// 	fmt.Println(errList)
	// 	os.Exit(1)
	// }

	// out, err := os.Open("out.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer out.Close()

	// res, err := ERPRequest("POST", "/api/resource/Farm Application", out)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(res.StatusCode)

	// f, err := os.Create("template.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// io.Copy(f, res.Body)

	// os.Exit(0)
}

func main() {
	// Setup handlers once
	mux := http.NewServeMux()
	mux.HandleFunc("/", helpers.ChainFuncs(handlers.SaveRequest, handlers.RootHandler))
	mux.HandleFunc("/favicon.ico", handlers.FaviconHandler)

	if local {
		log.Println("Serving on http://localhost:8080")
		err := http.ListenAndServe("localhost:8080", mux)
		if err != nil {
			log.Fatalf("Local server error: %v", err)
		}
	} else {
		// Start ngrok tunnel
		listener, err := ngrok.Listen(context.Background(),
			config.HTTPEndpoint(
				config.WithDomain(os.Getenv("DOMAIN")),
			),
			ngrok.WithAuthtokenFromEnv(),
		)
		if err != nil {
			log.Fatalf("Ngrok listener error: %v", err)
		}

		log.Println("App URL (ngrok):", listener.URL())
		err = http.Serve(listener, mux)
		if err != nil {
			log.Fatalf("Ngrok server error: %v", err)
		}
	}
}
