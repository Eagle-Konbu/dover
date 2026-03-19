package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.SetFlags(0)
	lambda.Start(handler)
}
