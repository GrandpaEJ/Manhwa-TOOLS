package tools

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	gg "github.com/GrandpaEJ/advancegg"
	"manhwa-tools-backend/internal/vision"
)

// RegisterTranslatorTool registers the /translate endpoint.
func RegisterTranslatorTool(router *gin.RouterGroup, modelPath string) {
	router.POST("/translate", func(c *gin.Context) {
		handleTranslateRequest(c, modelPath)
	})
}

func handleTranslateRequest(c *gin.Context, modelPath string) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode image"})
		return
	}

	tempFile, err := os.CreateTemp("", "upload-*.jpg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error creating temp file"})
		return
	}
	defer os.Remove(tempFile.Name())

	jpeg.Encode(tempFile, img, nil)
	tempFile.Close()

	visionResp, err := vision.SegmentTextBubbles(tempFile.Name(), modelPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "detection failed: " + err.Error()})
		return
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dc := gg.NewContextForImage(img)

	// Step 1: Erase original text (white fill the polygons)
	dc.SetColor(color.White)
	for _, block := range visionResp.Blocks {
		if len(block.Polygon) == 0 {
			continue
		}
		for i, pt := range block.Polygon {
			px := float64(pt[0]) * float64(width) / 1000.0
			py := float64(pt[1]) * float64(height) / 1000.0
			if i == 0 {
				dc.MoveTo(px, py)
			} else {
				dc.LineTo(px, py)
			}
		}
		dc.ClosePath()
		dc.Fill()
	}

	// Step 2: Render Bengali Text
	if err := dc.LoadFontFace("NotoSansBengali.ttf", 28); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load font: " + err.Error()})
		return
	}
	dc.SetColor(color.Black)

	for _, block := range visionResp.Blocks {
		if len(block.Polygon) == 0 {
			continue
		}
		// BoundingBox: [ymin, xmin, ymax, xmax] (0-1000)
		yMin := float64(block.BoundingBox[0]) * float64(height) / 1000.0
		xMin := float64(block.BoundingBox[1]) * float64(width) / 1000.0
		yMax := float64(block.BoundingBox[2]) * float64(height) / 1000.0
		xMax := float64(block.BoundingBox[3]) * float64(width) / 1000.0

		cx := (xMin + xMax) / 2.0
		cy := (yMin + yMax) / 2.0
		w := xMax - xMin

		// Dummy Bengali text to test rendering
		text := "অসাধারণ! এটা কাজ করছে।" // "Awesome! This is working."
		
		// Draw text centered in the bounding box
		dc.DrawStringWrapped(text, cx, cy, 0.5, 0.5, w*0.8, 1.2, gg.AlignCenter)
	}

	newImg := dc.Image()
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 95})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode image"})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
}
