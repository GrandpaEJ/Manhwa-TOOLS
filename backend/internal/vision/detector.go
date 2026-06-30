package vision

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"golang.org/x/image/draw"
	"manhwa-tools-backend/internal/models"
)

func loadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func getSliceHeight(imgW int) int {
	sliceHeight := imgW
	if sliceHeight > 1000 {
		sliceHeight = 1000
	}
	return sliceHeight
}

func getStep(sliceHeight int) int {
	step := sliceHeight - (sliceHeight / 4)
	if step <= 0 {
		step = sliceHeight
	}
	return step
}

// SegmentTextBubbles performs native ONNX inference on an image path and returns segmented bubbles
func SegmentTextBubbles(imagePath string, modelPath string) (*models.VisionResponse, error) {
	img, err := loadImage(imagePath)
	if err != nil {
		return nil, err
	}

	session, inputTensor, out0T, out1T, err := initSessionTensors(modelPath)
	if err != nil {
		return nil, err
	}
	defer session.Destroy()
	defer inputTensor.Destroy()
	defer out0T.Destroy()
	defer out1T.Destroy()

	var allBlocks []models.TranslationBlock
	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	sliceHeight := getSliceHeight(imgW)
	step := getStep(sliceHeight)

	for yOffset := 0; yOffset < imgH; yOffset += step {
		y2 := yOffset + sliceHeight
		if y2 > imgH {
			y2 = imgH
			yOffset = y2 - sliceHeight
			if yOffset < 0 {
				yOffset = 0
			}
		}

		sliceImg := image.NewRGBA(image.Rect(0, 0, imgW, y2-yOffset))
		draw.Draw(sliceImg, sliceImg.Bounds(), img, image.Point{0, yOffset}, draw.Src)

		inputData, scale, offX, offY := preprocessImageForYOLO(sliceImg)
		
		tensorData := inputTensor.GetData()
		copy(tensorData, inputData)

		if err := session.Run(); err != nil {
			return nil, err
		}

		out0, out1 := out0T.GetData(), out1T.GetData()
		
		boxes := parseYoloOutput(out0)
		survivors := applyNonMaxSuppression(boxes, 0.45)
		
		blocks := extractPolygonMasks(survivors, out1, scale, offX, offY, yOffset, float32(imgW), float32(imgH))
		allBlocks = append(allBlocks, blocks...)

		if y2 >= imgH {
			break
		}
	}

	finalBlocks := deduplicateBlocks(allBlocks)
	return &models.VisionResponse{Blocks: finalBlocks}, nil
}
