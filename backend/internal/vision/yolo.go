package vision

import (
	"image"
	"math"

	"golang.org/x/image/draw"
	"manhwa-tools-backend/internal/models"
)

func sigmoid(x float32) float32 {
	return float32(1.0 / (1.0 + math.Exp(float64(-x))))
}

func preprocessImageForYOLO(img image.Image) ([]float32, float32, int, int) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	scale := 640.0 / math.Max(float64(w), float64(h))
	nw, nh := int(float64(w)*scale), int(float64(h)*scale)
	resized := image.NewRGBA(image.Rect(0, 0, 640, 640))
	for i := 0; i < len(resized.Pix); i += 4 {
		resized.Pix[i], resized.Pix[i+1], resized.Pix[i+2], resized.Pix[i+3] = 114, 114, 114, 255
	}
	offX, offY := (640-nw)/2, (640-nh)/2
	draw.BiLinear.Scale(resized, image.Rect(offX, offY, offX+nw, offY+nh), img, bounds, draw.Over, nil)
	input := make([]float32, 3*640*640)
	for y := 0; y < 640; y++ {
		for x := 0; x < 640; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			input[0*640*640+y*640+x] = float32(r>>8) / 255.0
			input[1*640*640+y*640+x] = float32(g>>8) / 255.0
			input[2*640*640+y*640+x] = float32(b>>8) / 255.0
		}
	}
	return input, float32(scale), offX, offY
}

func parseYoloOutput(out0 []float32) []YoloDetection {
	var boxes []YoloDetection
	for i := 0; i < 8400; i++ {
		score := out0[4*8400+i]
		if score > 0.25 {
			b := YoloDetection{x: out0[0*8400+i], y: out0[1*8400+i], w: out0[2*8400+i], h: out0[3*8400+i], score: score, index: i, coeffs: make([]float32, 32)}
			for k := 0; k < 32; k++ {
				b.coeffs[k] = out0[(5+k)*8400+i]
			}
			boxes = append(boxes, b)
		}
	}
	return boxes
}

func extractPolygonMasks(survivors []YoloDetection, out1 []float32, scale float32, offX, offY, yOffset int, fImgW, fImgH float32) []models.TranslationBlock {
	var blocks []models.TranslationBlock
	for _, s := range survivors {
		realX, realY := (s.x-float32(offX))/scale, (s.y-float32(offY))/scale
		realW, realH := s.w/scale, s.h/scale
		globalY := realY + float32(yOffset)

		mXmin, mYmin := int(math.Max(0, math.Floor(float64(s.x-s.w/2)/4))), int(math.Max(0, math.Floor(float64(s.y-s.h/2)/4)))
		mXmax, mYmax := int(math.Min(159, math.Ceil(float64(s.x+s.w/2)/4))), int(math.Min(159, math.Ceil(float64(s.y+s.h/2)/4)))
		mask := make([][]bool, 160)
		for j := range mask {
			mask[j] = make([]bool, 160)
		}
		for y := mYmin; y <= mYmax; y++ {
			for x := mXmin; x <= mXmax; x++ {
				var val float32
				for k := 0; k < 32; k++ {
					val += s.coeffs[k] * out1[k*160*160+y*160+x]
				}
				if sigmoid(val) > 0.5 {
					mask[y][x] = true
				}
			}
		}
		poly160 := traceMaskContour(mask, mXmin, mYmin, mXmax, mYmax)
		if len(poly160) == 0 {
			continue
		}

		var normalizedPoly [][2]int
		for _, pt := range poly160 {
			pxReal, pyReal := (float32(pt[0]*4)-float32(offX))/scale, (float32(pt[1]*4)-float32(offY))/scale
			globalPyReal := pyReal + float32(yOffset)
			normalizedPoly = append(normalizedPoly, [2]int{int((pxReal / fImgW) * 1000), int((globalPyReal / fImgH) * 1000)})
		}

		ymin := int(((globalY - realH/2) / fImgH) * 1000)
		xmin := int(((realX - realW/2) / fImgW) * 1000)
		ymax := int(((globalY + realH/2) / fImgH) * 1000)
		xmax := int(((realX + realW/2) / fImgW) * 1000)

		blocks = append(blocks, models.TranslationBlock{
			BoundingBox: [4]int{ymin, xmin, ymax, xmax},
			Polygon:     normalizedPoly,
			Type:        "dialogue",
		})
	}
	return blocks
}
