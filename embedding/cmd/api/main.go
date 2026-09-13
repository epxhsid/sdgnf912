package main

import (
	"fmt"
	"log"
	"math"

	"github.com/epxhsid/embedding"
	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	ort.SetSharedLibraryPath("./onnxruntime/lib/libonnxruntime.so")

	if err := ort.InitializeEnvironment(); err != nil {
		log.Fatal(err)
	}
	defer ort.DestroyEnvironment()

	embedder, err := embedding.New(
		"./models/bge-small-en-v1.5/model.onnx",
		"./models/bge-small-en-v1.5/tokenizer.json",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer embedder.Close()

	queryVector, err := embedder.Embed(
		"When was the company founded?",
		embedding.Query,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Query embedding dimensions: %d\n", len(queryVector))
	fmt.Printf("Query first 10 values: %v\n", queryVector[:10])

	printNorm(queryVector)

	documentVector, err := embedder.Embed(
		"The company was founded in 1998.",
		embedding.Document,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nDocument embedding dimensions: %d\n", len(documentVector))
	fmt.Printf("Document first 10 values: %v\n", documentVector[:10])

	printNorm(documentVector)
}

func printNorm(vector []float32) {
	var sumSquares float64

	for _, value := range vector {
		sumSquares += float64(value) * float64(value)
	}

	fmt.Printf("L2 norm: %.6f\n", math.Sqrt(sumSquares))
}
