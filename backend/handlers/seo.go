package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/config"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"

	"github.com/verbumby/verbum/backend/chttp"
)

// RobotsTXT responds to /robots.txt request
func RobotsTXT(w http.ResponseWriter, rctx *chttp.Context) error {
	tmpl := `User-agent: *
Sitemap: %s/sitemap-index.xml
`
	body := fmt.Sprintf(tmpl, config.HTTPSCanonicalAddr())
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(body)); err != nil {
		return fmt.Errorf("write response body: %w", err)
	}
	return nil
}

// SitemapIndex handles sitemap index request
func SitemapIndex(w http.ResponseWriter, rctx *chttp.Context) error {
	dicts := dictionary.GetAllListed()

	type Sitemap struct {
		Loc string `xml:"loc"`
	}
	type sitemapindex struct {
		XMLNS   string    `xml:"xmlns,attr"`
		Sitemap []Sitemap `xml:"sitemap"`
	}

	result := sitemapindex{
		XMLNS:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemap: []Sitemap{},
	}

	for _, d := range dicts {
		countresp := struct {
			Count uint64 `json:"count"`
		}{}
		url := fmt.Sprintf("/dict-%s/_count", d.IndexID())
		if err := storage.Get(url, &countresp); err != nil {
			return fmt.Errorf("storage get %s docs count: %w", d.ID(), err)
		}

		for i := uint64(0); i <= countresp.Count/10000; i++ {
			result.Sitemap = append(result.Sitemap, Sitemap{
				Loc: fmt.Sprintf("%s/%s/sitemap-%d.xml", config.HTTPSCanonicalAddr(), d.ID(), i),
			})
		}
		// fmt.Println(d.ID)
	}

	w.Header().Set("Content-Type", "text/xml")
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("write xml header: %w", err)
	}

	if err := xml.NewEncoder(w).Encode(result); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}

// SitemapOfDictionary handles dictionary sitemap request
func SitemapOfDictionary(w http.ResponseWriter, rctx *chttp.Context) error {
	d := dictionary.Get(chi.URLParam(rctx.R, "dictionary"))
	if d == nil {
		return APINotFound(w, rctx)
	}

	nstr := chi.URLParam(rctx.R, "n")
	n, _ := strconv.ParseUint(nstr, 10, 64)

	reqbody := map[string]interface{}{
		"from":    n * 10000,
		"size":    10000,
		"sort":    []string{"_doc"},
		"_source": "ModifiedAt",
	}
	respbody := struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}{}
	surl := fmt.Sprintf("/dict-%s/_search", d.IndexID())
	if err := storage.Post(surl, reqbody, &respbody); err != nil {
		return fmt.Errorf("sotrage post: %w", err)
	}

	type urlt struct {
		Loc        string `xml:"loc"`
		Changefreq string `xml:"changefreq"`
	}
	type urlset struct {
		XMLNS string `xml:"xmlns,attr"`
		URL   []urlt `xml:"url"`
	}

	result := urlset{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL:   []urlt{},
	}

	for _, a := range respbody.Hits.Hits {
		result.URL = append(result.URL, urlt{
			Loc:        fmt.Sprintf("%s/%s/%s", config.HTTPSCanonicalAddr(), d.ID(), url.PathEscape(a.ID)),
			Changefreq: "yearly",
		})
	}

	w.Header().Set("Content-Type", "text/xml")
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("write xml header: %w", err)
	}
	if err := xml.NewEncoder(w).Encode(result); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
