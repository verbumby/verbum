package chttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
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
					log.Println(errors.Wrap(err, "flush access log buffer"))
				} else {
					accessLog.buff = make([]accessLogMsg, 0, 100)
				}
			}
		case <-time.After(5 * time.Second):
			if len(accessLog.buff) > 0 {
				if err := elasticAccessLogWriterFlush(); err != nil {
					log.Println(errors.Wrap(err, "flush access log buffer"))
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
			return errors.Wrap(err, "encode action and meta data to json")
		}

		if err := enc.Encode(doc); err != nil {
			return errors.Wrap(err, "encode doc to json")
		}
	}

	err := storage.Post("/_bulk", buff, nil)

	return errors.Wrap(err, "bulk post to storage")
}
