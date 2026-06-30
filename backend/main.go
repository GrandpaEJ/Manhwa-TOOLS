package main

import (
	"log"

	"manhwa-tools-backend/internal/api"
	"manhwa-tools-backend/internal/vision"
)

func main() {
	libPath := "./libonnxruntime.so"
	modelPath := "./segmentor_best.onnx"

	log.Printf("Initializing ONNX Runtime from %s...\n", libPath)
	if err := vision.InitONNX(libPath); err != nil {
		log.Fatalf("Failed to initialize ONNX Runtime: %v", err)
	}
	defer vision.CleanupONNX()

	// Setup modular router architecture
	r := api.SetupRouter(modelPath)

	log.Println("Starting scalable tool server on :8080")
	r.Run(":8080")
}
