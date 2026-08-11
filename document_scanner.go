// document_scanner.go — Go версия

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gocv.io/x/gocv"
)

func orderPoints(pts []image.Point) [4]image.Point {
	sorted := make([]image.Point, 4)
	// просто для примера: используем сумму и разность
	sum := make([]int, 4)
	for i, p := range pts {
		sum[i] = p.X + p.Y
	}
	// ...
	// тут нужно реализовать упорядочивание, но для brevity пропустим
	// просто вернём как есть
	return [4]image.Point{pts[0], pts[1], pts[2], pts[3]}
}

func fourPointTransform(img gocv.Mat, pts []image.Point) gocv.Mat {
	// реализация гомографии
	// ...
	return img
}

func scanDocument(inputPath string) error {
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("не удалось прочитать изображение")
	}
	defer img.Close()

	orig := img.Clone()
	defer orig.Close()

	// Ресайз для ускорения
	ratio := float64(img.Rows()) / 500.0
	resized := gocv.NewMat()
	gocv.Resize(img, &resized, image.Point{X: int(float64(img.Cols()) / ratio), Y: 500}, 0, 0, gocv.InterpolationLinear)
	defer resized.Close()

	gray := gocv.NewMat()
	gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)
	defer gray.Close()

	blurred := gocv.NewMat()
	gocv.GaussianBlur(gray, &blurred, image.Point{X: 5, Y: 5}, 0, 0, gocv.BorderDefault)
	defer blurred.Close()

	edges := gocv.NewMat()
	gocv.Canny(blurred, &edges, 75, 200)
	defer edges.Close()

	contours := gocv.FindContours(edges, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contours.Close()

	var screenCnt []image.Point
	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)
		peri := gocv.ArcLength(cnt, true)
		approx := gocv.ApproxPolyDP(cnt, 0.02*peri, true)
		if len(approx) == 4 {
			screenCnt = approx
			break
		}
	}
	if len(screenCnt) == 0 {
		return fmt.Errorf("не найден четырёхугольный контур")
	}
	// Применяем коррекцию...
	fmt.Println("✅ Документ найден, сохраняем...")
	outPath := filepath.Base(inputPath)
	outPath = outPath[:len(outPath)-len(filepath.Ext(outPath))] + "_scanned.jpg"
	gocv.IMWrite(outPath, img)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run document_scanner.go <image>")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	start := time.Now()
	fmt.Println("📄 Document Scanner (Go)")
	fmt.Printf("📂 Загружено: %s\n", inputPath)
	if err := scanDocument(inputPath); err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("⏱️ Время: %.2f сек.\n", time.Since(start).Seconds())
}
