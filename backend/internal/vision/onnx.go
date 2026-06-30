package vision

import (
	"fmt"

	"github.com/yalue/onnxruntime_go"
)

var isInitialized bool

// InitONNX sets the ONNX shared library path and initializes the environment.
func InitONNX(libPath string) error {
	if isInitialized {
		return nil
	}
	onnxruntime_go.SetSharedLibraryPath(libPath)
	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		return fmt.Errorf("failed to init onnx environment: %v", err)
	}
	isInitialized = true
	return nil
}

func CleanupONNX() {
	if isInitialized {
		onnxruntime_go.DestroyEnvironment()
		isInitialized = false
	}
}

func initSessionTensors(modelPath string) (*onnxruntime_go.Session[float32], *onnxruntime_go.Tensor[float32], *onnxruntime_go.Tensor[float32], *onnxruntime_go.Tensor[float32], error) {
	dummyInput := make([]float32, 3*640*640)
	inputTensor, _ := onnxruntime_go.NewTensor([]int64{1, 3, 640, 640}, dummyInput)
	out0T, _ := onnxruntime_go.NewEmptyTensor[float32]([]int64{1, 37, 8400})
	out1T, _ := onnxruntime_go.NewEmptyTensor[float32]([]int64{1, 32, 160, 160})

	session, err := onnxruntime_go.NewSession[float32](modelPath, []string{"images"}, []string{"output0", "output1"}, []*onnxruntime_go.Tensor[float32]{inputTensor}, []*onnxruntime_go.Tensor[float32]{out0T, out1T})
	if err != nil {
		inputTensor.Destroy()
		out0T.Destroy()
		out1T.Destroy()
		return nil, nil, nil, nil, err
	}
	return session, inputTensor, out0T, out1T, nil
}
