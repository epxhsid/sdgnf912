package main

import (
	"fmt"
	"log"
	"math"

	"github.com/sugarme/tokenizer/pretrained"
	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	ort.SetSharedLibraryPath("./onnxruntime/lib/libonnxruntime.so")
	if err := ort.InitializeEnvironment(); err != nil {
		log.Fatal(err)
	}
	defer ort.DestroyEnvironment()

	tk, err := pretrained.FromFile(
		"./models/bge-small-en-v1.5/tokenizer.json",
	)
	if err != nil {
		log.Fatal(err)
	}

	session, err := ort.NewDynamicAdvancedSession(
		"./models/bge-small-en-v1.5/model.onnx",
		[]string{
			"input_ids",
			"attention_mask",
			"token_type_ids",
		},
		[]string{
			"last_hidden_state",
		},
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Destroy()

	embedding, err := embed(
		tk,
		session,
		"hello world",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Embedding dimensions: %d\n", len(embedding))
	fmt.Printf("First 10 values: %v\n", embedding[:10])

	var sum float32
	for _, v := range embedding {
		sum += v * v
	}

	fmt.Printf("L2 norm: %f\n", math.Sqrt(float64(sum)))
}
