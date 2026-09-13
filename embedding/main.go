package main

import (
	"fmt"

	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	ort.SetSharedLibraryPath("./onnxruntime/lib/libonnxruntime.so")

	err := ort.InitializeEnvironment()
	if err != nil {
		panic(err)
	}
	defer ort.DestroyEnvironment()

	fmt.Println("ONNX Runtime initialized successfully")
}
