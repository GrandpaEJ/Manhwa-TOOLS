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

// RegisterEraserTool registers the /clean endpoint for the eraser tool.
func RegisterEraserTool(router *gin.RouterGroup, modelPath string) {
	router.POST("/clean", func(c *gin.Context) {
		handleCleanRequest(c, modelPath)
	})
}

func handleCleanRequest(c *gin.Context, modelPath string) {
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

	newImg := dc.Image()

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 95})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode image"})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
}
