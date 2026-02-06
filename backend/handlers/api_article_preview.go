package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
)

func APIArticlePreview(w http.ResponseWriter, rctx *chttp.Context) error {
	d := dictionary.Get(chi.URLParam(rctx.R, "dictionary"))
	if d == nil {
		return APINotFound(w, rctx)
	}

	aIDRaw := chi.URLParam(rctx.R, "article")
	var err error
	aID, err := url.QueryUnescape(aIDRaw)
	if err != nil {
		return fmt.Errorf("unescape aID: %w", err)
	}

	a, err := article.Get(d, aID)
	if err != nil {
		return fmt.Errorf("get article: %w", err)
	}

	if a.ID == "" {
		return APINotFound(w, rctx)
	}

	tmpFile, err := os.CreateTemp("", "og_preview_*.png")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	screenshotPath := tmpFile.Name()
	tmpFile.Close()

	defer os.Remove(screenshotPath)

	targetURL := fmt.Sprintf("http://127.0.0.1:8080/%s/%s?og_preview=true", d.ID(), aIDRaw)
	cmd := exec.Command("firefox-esr",
		"--headless",
		"--window-size", "800,420",
		"--screenshot", screenshotPath,
		targetURL,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("firefox failed: %s %w", output, err)
	}

	imgData, err := os.ReadFile(screenshotPath)
	if err != nil {
		return fmt.Errorf("read img data: %w", err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(imgData)

	return nil
}
