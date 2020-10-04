package chttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/verbumby/verbum/backend/pkg/storage"
)

type accessLogMsg struct {
	TS        time.Time
	IP        string
	Path      string
	Query     string
	UserAgent string
	Referer   string
}

var accessLog struct {
	ch   chan accessLogMsg
	buff []accessLogMsg
}

// InitAccessLog initializes the middleware that writes access log to elastic
func InitAccessLog() {
	accessLog.ch = make(chan accessLogMsg, 15)
	accessLog.buff = make([]accessLogMsg, 0, 100)
	go elasticAccessLogWriter()
}

func elasticAccessLogMiddleware(f HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, ctx *Context) error {
		query, _ := url.QueryUnescape(ctx.R.URL.RawQuery)
		accessLog.ch <- accessLogMsg{
			TS:        time.Now().UTC(),
			IP:        ctx.R.RemoteAddr,
			Path:      ctx.R.URL.Path,
			Query:     query,
			UserAgent: ctx.R.UserAgent(),
			Referer:   ctx.R.Referer(),
		}
		result := f(w, ctx)
		return result
	}
}

func elasticAccessLogWriter() {
	for {
		select {
		case m := <-accessLog.ch:
			accessLog.buff = append(accessLog.buff, m)
			if len(accessLog.buff) > 100 {
				if err := elasticAccessLogWriterFlush(); err != nil {
					log.Println(fmt.Errorf("flush access log buffer: %w", err))
				} else {
					accessLog.buff = make([]accessLogMsg, 0, 100)
				}
			}
		case <-time.After(5 * time.Second):
			if len(accessLog.buff) > 0 {
				if err := elasticAccessLogWriterFlush(); err != nil {
					log.Println(fmt.Errorf("flush access log buffer: %w", err))
				} else {
					accessLog.buff = make([]accessLogMsg, 0, 100)
				}
			}
		}

	}
}

func elasticAccessLogWriterFlush() error {
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	for _, doc := range accessLog.buff {
		index := fmt.Sprintf("access-log-%d-%d", doc.TS.Year(), doc.TS.Month())
		if err := enc.Encode(map[string]interface{}{
			"index": map[string]interface{}{"_index": index, "_type": "_doc"},
		}); err != nil {
			return fmt.Errorf("encode action and meta data to json: %w", err)
		}

		if err := enc.Encode(doc); err != nil {
			return fmt.Errorf("encode doc to json: %w", err)
		}
	}

	if err := storage.Post("/_bulk", buff, nil); err != nil {
		return fmt.Errorf("bulk post to storage: %w", err)
	}
	return nil
}
