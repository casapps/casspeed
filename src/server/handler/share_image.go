package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/casapps/casspeed/src/server/store"
	"github.com/go-chi/chi/v5"
)

type ShareImageHandler struct {
	store store.Store
}

func NewShareImageHandler(st store.Store) *ShareImageHandler {
	return &ShareImageHandler{store: st}
}

func (h *ShareImageHandler) GetSharePNG(w http.ResponseWriter, r *http.Request) {
	shareCode := chi.URLParam(r, "code")

	test, err := h.store.GetSpeedTestByShareCode(r.Context(), shareCode)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if test == nil {
		http.Error(w, "Share not found", http.StatusNotFound)
		return
	}

	h.store.IncrementShareViews(r.Context(), shareCode)

	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	bgColor := color.RGBA{15, 15, 35, 255}
	for y := 0; y < 630; y++ {
		for x := 0; x < 1200; x++ {
			img.Set(x, y, bgColor)
		}
	}

	drawText(img, 100, 100, "casspeed", color.White)
	drawText(img, 100, 200, fmt.Sprintf("Download: %.1f Mbps", test.DownloadMbps), color.White)
	drawText(img, 100, 280, fmt.Sprintf("Upload: %.1f Mbps", test.UploadMbps), color.White)
	drawText(img, 100, 360, fmt.Sprintf("Ping: %.1f ms", test.PingMs), color.White)
	drawText(img, 100, 440, test.Timestamp.Format("2006-01-02 15:04:05"), color.RGBA{136, 136, 136, 255})

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		http.Error(w, "Image generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(buf.Bytes())
}

func (h *ShareImageHandler) GetShareSVG(w http.ResponseWriter, r *http.Request) {
	shareCode := chi.URLParam(r, "code")

	test, err := h.store.GetSpeedTestByShareCode(r.Context(), shareCode)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if test == nil {
		http.Error(w, "Share not found", http.StatusNotFound)
		return
	}

	h.store.IncrementShareViews(r.Context(), shareCode)

	svg := fmt.Sprintf(`<svg width="1200" height="630" xmlns="http://www.w3.org/2000/svg">
  <rect width="1200" height="630" fill="#0f0f23"/>
  <text x="100" y="100" font-family="Arial" font-size="48" fill="white">casspeed</text>
  <text x="100" y="200" font-family="Arial" font-size="32" fill="white">Download: %.1f Mbps</text>
  <text x="100" y="280" font-family="Arial" font-size="32" fill="white">Upload: %.1f Mbps</text>
  <text x="100" y="360" font-family="Arial" font-size="32" fill="white">Ping: %.1f ms</text>
  <text x="100" y="440" font-family="Arial" font-size="24" fill="#888">%s</text>
</svg>`, test.DownloadMbps, test.UploadMbps, test.PingMs, test.Timestamp.Format("2006-01-02 15:04:05"))

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	fmt.Fprint(w, svg)
}

func drawText(img *image.RGBA, x, y int, text string, col color.Color) {
	for i, ch := range text {
		drawChar(img, x+(i*12), y, ch, col)
	}
}

func drawChar(img *image.RGBA, x, y int, ch rune, col color.Color) {
	for dy := 0; dy < 16; dy++ {
		for dx := 0; dx < 10; dx++ {
			if (dy+dx)%3 == 0 {
				img.Set(x+dx, y+dy, col)
			}
		}
	}
}
