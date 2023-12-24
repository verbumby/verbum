package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/frontend"
)

func IndexHTML(w http.ResponseWriter, ctx *chttp.Context) error {
	indexHTMLFilename := "dist/index.html"
	indexHtmlBytes, err := frontend.Dist.ReadFile(indexHTMLFilename)
	if err != nil {
		return fmt.Errorf("read %s: %w", indexHTMLFilename, err)
	}
	indexHTML := string(indexHtmlBytes)

	publicDirname := "dist/public"
	publicAssets, err := frontend.Dist.ReadDir(publicDirname)
	if err != nil {
		return fmt.Errorf("read list files %s: %w", publicDirname, err)
	}

	jss := []string{}
	csss := []string{}
	for _, p := range publicAssets {
		if strings.HasSuffix(p.Name(), ".js") {
			jss = append(jss, fmt.Sprintf(`<script defer="defer" src="/statics/%s"></script>`, p.Name()))
		}
		if strings.HasSuffix(p.Name(), ".css") {
			csss = append(csss, fmt.Sprintf(`<link href="/statics/%s" rel="stylesheet">`, p.Name()))
		}
	}
	indexHTML = strings.ReplaceAll(indexHTML, "CSS_BUNDLES_PLACEHOLDER", strings.Join(csss, "\n"))
	indexHTML = strings.ReplaceAll(indexHTML, "JS_BUNDLES_PLACEHOLDER", strings.Join(jss, "\n"))

	if _, err := w.Write([]byte(indexHTML)); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}
