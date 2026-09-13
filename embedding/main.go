package main

import (
	"fmt"
	"log"

	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	ort.SetSharedLibraryPath("./onnxruntime/lib/libonnxruntime.so")

	if err := ort.InitializeEnvironment(); err != nil {
		log.Fatal(err)
	}
	defer ort.DestroyEnvironment()

	session, err := ort.NewDynamicAdvancedSession(
		"./models/bge-small-en-v1.5/model.onnx",
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Destroy()

	fmt.Println("ONNX model loaded successfully")

	inputs, outputs, err := ort.GetInputOutputInfo(
		"./models/bge-small-en-v1.5/model.onnx",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("INPUTS:")
	for _, input := range inputs {
		fmt.Printf(
			"  name=%s dimensions=%v type=%v\n",
			input.Name,
			input.Dimensions,
			input.DataType,
		)
	}

	fmt.Println("OUTPUTS:")
	for _, output := range outputs {
		fmt.Printf(
			"  name=%s dimensions=%v type=%v\n",
			output.Name,
			output.Dimensions,
			output.DataType,
		)
	}
}
