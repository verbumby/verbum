package serverrender

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"

	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/frontend"
	"rogchap.com/v8go"
)

type serverRenderer struct {
	handler   http.Handler
	v8ctx     *v8go.Context
	mu        *sync.Mutex
	indexHTML string
}

func New(handler http.Handler) (*serverRenderer, error) {
	result := &serverRenderer{
		handler: handler,
	}

	if err := result.prepareIndexHTML(); err != nil {
		return nil, fmt.Errorf("prepare index.html: %w", err)
	}

	var err error
	result.v8ctx, err = result.newV8Ctx()
	if err != nil {
		return nil, fmt.Errorf("create v8 ctx: %w", err)
	}

	result.mu = &sync.Mutex{}

	return result, nil
}

func (r *serverRenderer) prepareIndexHTML() error {
	indexHTMLFilename := "dist/index.html"
	indexHtmlBytes, err := frontend.Dist.ReadFile(indexHTMLFilename)
	if err != nil {
		return fmt.Errorf("read %s: %w", indexHTMLFilename, err)
	}
	r.indexHTML = string(indexHtmlBytes)

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
	r.indexHTML = strings.ReplaceAll(r.indexHTML, "CSS_BUNDLES_PLACEHOLDER", strings.Join(csss, "\n"))
	r.indexHTML = strings.ReplaceAll(r.indexHTML, "JS_BUNDLES_PLACEHOLDER", strings.Join(jss, "\n"))
	return nil
}

func (r *serverRenderer) newV8Ctx() (*v8go.Context, error) {
	serverRenderFilename := "dist/server.js"
	serverRender, err := frontend.Dist.ReadFile(serverRenderFilename)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", serverRenderFilename, err)
	}

	v8ctx := v8go.NewContext()

	vbridge := v8go.NewFunctionTemplate(v8ctx.Isolate(), func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		// TODO: common err prefix

		rawUrl := info.Args()[0].String()
		u, err := url.Parse(rawUrl)
		if err != nil {
			log.Printf("parse url: %v", err)
			return v8go.Null(v8ctx.Isolate())
		}

		w := httptest.NewRecorder()
		rctx := &http.Request{URL: u}
		r.handler.ServeHTTP(w, rctx)

		if w.Code == http.StatusNotFound {
			return v8go.Null(v8ctx.Isolate())
		}

		result, err := v8go.JSONParse(v8ctx, w.Body.String())
		if err != nil {
			log.Printf("parse as js from json: %v", err)
			return v8go.Null(v8ctx.Isolate())
		}

		return result
	})

	if err := v8ctx.Global().Set("verbumV8Bridge", vbridge.GetFunction(v8ctx)); err != nil {
		return nil, fmt.Errorf("set v8 bridge function: %w", err)
	}

	_, err = v8ctx.RunScript(string(serverRender), serverRenderFilename)
	if err != nil {
		return nil, fmt.Errorf("v8 run %s: %w", serverRenderFilename, err)
	}

	return v8ctx, nil
}

func (r *serverRenderer) ServeHTTP(w http.ResponseWriter, rctx *chttp.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var render *v8go.Function
	if f, err := r.v8ctx.Global().Get("render"); err != nil {
		return fmt.Errorf("get render function from renderer context: %w", err)
	} else {
		render, err = f.AsFunction()
		if err != nil {
			return fmt.Errorf("cast render key to a function: %w", err)
		}
	}

	url := "https://" + rctx.R.Host + rctx.R.RequestURI
	vurl, err := v8go.NewValue(r.v8ctx.Isolate(), url)
	if err != nil {
		return fmt.Errorf("new url value: %w", err)
	}

	vres, err := promiseResolver(r.v8ctx)(render.Call(v8go.Undefined(r.v8ctx.Isolate()), vurl))
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	res, err := vres.AsObject()
	if err != nil {
		return fmt.Errorf("render result convert to object: %w", err)
	}

	str, err := v8go.JSONStringify(r.v8ctx, res)
	if err != nil {
		return fmt.Errorf("stringify render result: %w", err)
	}

	result := RenderResult{}
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		return fmt.Errorf("unmarshal render result: %w", err)
	}

	if result.Location != "" {
		http.Redirect(w, rctx.R, result.Location, http.StatusMovedPermanently)
		return nil
	}

	if result.StatusCode > 0 {
		w.WriteHeader(result.StatusCode)
	}

	body := r.indexHTML
	body = strings.ReplaceAll(body, "HEAD_TITLE_PLACEHOLDER", result.Title)
	body = strings.ReplaceAll(body, "HEAD_META_PLACEHOLDER", result.Meta)
	body = strings.ReplaceAll(body, "PRELOADED_STATE_PLACEHOLDER", result.State)
	body = strings.ReplaceAll(body, "BODY_PLACEHOLDER", result.Body)

	fmt.Fprint(w, body)

	return nil
}

func promiseResolver(v8ctx *v8go.Context) func(*v8go.Value, error) (*v8go.Value, error) {
	return func(val *v8go.Value, err error) (*v8go.Value, error) {
		if err != nil || !val.IsPromise() {
			return val, err
		}
		for {
			switch p, _ := val.AsPromise(); p.State() {
			case v8go.Fulfilled:
				return p.Result(), nil
			case v8go.Rejected:
				return nil, errors.New(p.Result().DetailString())
			case v8go.Pending:
				v8ctx.PerformMicrotaskCheckpoint() // run VM to make progress on the promise
				fmt.Printf(".")
				// go round the loop again...
			default:
				return nil, fmt.Errorf("illegal v8go.Promise state %d", p) // unreachable
			}
		}
	}
}

type RenderResult struct {
	StatusCode int
	Location   string
	Title      string
	Meta       string
	State      string
	Body       string
}
