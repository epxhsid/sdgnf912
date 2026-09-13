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

	vector, err := embedder.Embed("hello world")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Embedding dimensions: %d\n", len(vector))
	fmt.Printf("First 10 values: %v\n", vector[:10])

	var sumSquares float64
	for _, value := range vector {
		sumSquares += float64(value) * float64(value)
	}

	fmt.Printf("L2 norm: %.6f\n", math.Sqrt(sumSquares))
}
