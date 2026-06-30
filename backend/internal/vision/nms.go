package vision

import (
	"math"
	"sort"

	"manhwa-tools-backend/internal/models"
)

type YoloDetection struct {
	x, y, w, h float32
	score      float32
	class      int
	index      int
	coeffs     []float32
}

func calculateYoloIoU(b1, b2 YoloDetection) float32 {
	x1 := float32(math.Max(float64(b1.x-b1.w/2), float64(b2.x-b2.w/2)))
	y1 := float32(math.Max(float64(b1.y-b1.h/2), float64(b2.y-b2.h/2)))
	x2 := float32(math.Min(float64(b1.x+b1.w/2), float64(b2.x+b2.w/2)))
	y2 := float32(math.Min(float64(b1.y+b1.h/2), float64(b2.y+b2.h/2)))
	if x2 < x1 || y2 < y1 {
		return 0.0
	}
	intersection := (x2 - x1) * (y2 - y1)
	area1 := b1.w * b1.h
	area2 := b2.w * b2.h
	return intersection / (area1 + area2 - intersection)
}

func applyNonMaxSuppression(boxes []YoloDetection, iouThresh float32) []YoloDetection {
	sort.Slice(boxes, func(i, j int) bool {
		return boxes[i].score > boxes[j].score
	})
	var result []YoloDetection
	for len(boxes) > 0 {
		current := boxes[0]
		result = append(result, current)
		boxes = boxes[1:]
		var remaining []YoloDetection
		for _, b := range boxes {
			if calculateYoloIoU(current, b) < iouThresh {
				remaining = append(remaining, b)
			}
		}
		boxes = remaining
	}
	return result
}

func deduplicateBlocks(allBlocks []models.TranslationBlock) []models.TranslationBlock {
	var finalBlocks []models.TranslationBlock
	for _, b := range allBlocks {
		isDup := false
		for i, fb := range finalBlocks {
			if calculateBlockIoU(b.BoundingBox, fb.BoundingBox) > 0.3 {
				areaB := (b.BoundingBox[2] - b.BoundingBox[0]) * (b.BoundingBox[3] - b.BoundingBox[1])
				areaFB := (fb.BoundingBox[2] - fb.BoundingBox[0]) * (fb.BoundingBox[3] - fb.BoundingBox[1])
				if areaB > areaFB {
					finalBlocks[i] = b
				}
				isDup = true
				break
			}
		}
		if !isDup {
			finalBlocks = append(finalBlocks, b)
		}
	}
	return finalBlocks
}

func calculateBlockIoU(b1, b2 [4]int) float64 {
	y1 := math.Max(float64(b1[0]), float64(b2[0]))
	x1 := math.Max(float64(b1[1]), float64(b2[1]))
	y2 := math.Min(float64(b1[2]), float64(b2[2]))
	x2 := math.Min(float64(b1[3]), float64(b2[3]))
	if x2 < x1 || y2 < y1 {
		return 0.0
	}
	intersection := (x2 - x1) * (y2 - y1)
	area1 := float64((b1[2] - b1[0]) * (b1[3] - b1[1]))
	area2 := float64((b2[2] - b2[0]) * (b2[3] - b2[1]))
	return intersection / (area1 + area2 - intersection)
}
