package main

import (
	"flag"
	"log"

	"github.com/Kevin-Jii/tower-go/apidocs"
)

func main() {
	input := flag.String("input", "docs/swagger.json", "path to the swag-generated Swagger 2 JSON")
	output := flag.String("output", "docs/openapi.json", "path for the validated OpenAPI 3 JSON")
	flag.Parse()

	if err := apidocs.ConvertFile(*input, *output); err != nil {
		log.Fatal(err)
	}
	log.Printf("generated validated OpenAPI %s document at %s", apidocs.Version, *output)
}
