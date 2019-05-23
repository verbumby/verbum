package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/storage"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/chttp"
)

// RobotsTXT responds to /robots.txt request
func RobotsTXT(w http.ResponseWriter, rctx *chttp.Context) error {
	tmpl := `User-agent: *
Sitemap: %s/sitemap-index.xml
`
	body := fmt.Sprintf(tmpl, viper.GetString("https.canonicalAddr"))
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(body))
	return errors.Wrap(err, "write response body")
}

// SitemapIndex handles sitemap index request
func SitemapIndex(w http.ResponseWriter, rctx *chttp.Context) error {
	dicts, err := dictionary.GetAll()
	if err != nil {
		return errors.Wrap(err, "get all dictionaries")
	}

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
		url := fmt.Sprintf("/dict-%s/_count", d.ID)
		if err := storage.Get(url, &countresp); err != nil {
			return errors.Wrapf(err, "storage get %s docs count", d.ID)
		}

		for i := uint64(0); i <= countresp.Count/10000; i++ {
			result.Sitemap = append(result.Sitemap, Sitemap{
				Loc: fmt.Sprintf("%s/%s/sitemap-%d.xml", viper.GetString("https.canonicalAddr"), d.ID, i),
			})
		}
		// fmt.Println(d.ID)
	}

	w.Header().Set("Content-Type", "text/xml")
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return errors.Wrap(err, "write xml header")
	}
	return errors.Wrap(xml.NewEncoder(w).Encode(result), "encode response")
}

// SitemapOfDictionary handles dictionary sitemap request
func SitemapOfDictionary(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)
	dictID := vars["dictionary"]
	nstr := vars["n"]
	n, _ := strconv.ParseUint(nstr, 10, 64)

	reqbody := map[string]interface{}{
		"from":    n * 10000,
		"size":    10000,
		"sort":    []string{"_doc"},
		"_source": false,
	}
	respbody := struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}{}
	url := fmt.Sprintf("/dict-%s/_search", dictID)
	if err := storage.Post(url, reqbody, &respbody); err != nil {
		return errors.Wrap(err, "sotrage post")
	}

	type urlt struct {
		Loc string `xml:"loc"`
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
			Loc: fmt.Sprintf("%s/%s/%s", viper.GetString("https.canonicalAddr"), dictID, a.ID),
		})
	}

	w.Header().Set("Content-Type", "text/xml")
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return errors.Wrap(err, "write xml header")
	}
	return errors.Wrap(xml.NewEncoder(w).Encode(result), "encode response")
}
