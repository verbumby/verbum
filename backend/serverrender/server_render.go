package serverrender

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/frontend"
)

type serverRenderer struct {
	handler   http.Handler
	vm        *goja.Runtime
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
	result.vm, err = result.newVM()
	if err != nil {
		return nil, fmt.Errorf("create goja vm: %w", err)
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

func (r *serverRenderer) newVM() (*goja.Runtime, error) {
	serverRenderFilename := "dist/server.js"
	serverRender, err := frontend.Dist.ReadFile(serverRenderFilename)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", serverRenderFilename, err)
	}

	vm := goja.New()

	bridge := func(call goja.FunctionCall) goja.Value {
		rawUrl := call.Argument(0).Export().(string)
		u, err := url.Parse(rawUrl)
		if err != nil {
			log.Printf("parse url: %v", err)
			return goja.Undefined()
		}

		w := httptest.NewRecorder()
		rctx := &http.Request{URL: u}
		r.handler.ServeHTTP(w, rctx)

		if w.Code == http.StatusNotFound {
			return goja.Null()
		}

		if w.Code != http.StatusOK {
			msg := fmt.Sprintf("unexpected %d status code when calling %s: %s", w.Code, rawUrl, w.Body.String())
			panic(vm.ToValue(msg))
		}

		var result interface{}

		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			log.Printf("parse as js from json: %v", err)
			return goja.Null()
		}

		return vm.ToValue(result)
	}
	if err := vm.Set("verbumV8Bridge", bridge); err != nil {
		return nil, fmt.Errorf("set bridge function: %w", err)
	}

	if _, err := vm.RunScript(serverRenderFilename, string(serverRender)); err != nil {
		return nil, fmt.Errorf("run %s: %w", serverRenderFilename, err)
	}

	return vm, nil
}

func (r *serverRenderer) ServeHTTP(w http.ResponseWriter, rctx *chttp.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	render, ok := goja.AssertFunction(r.vm.Get("render"))
	if !ok {
		return fmt.Errorf("failed to get render function from goja vm")
	}

	url := "https://" + rctx.R.Host + rctx.R.RequestURI

	pres, err := render(goja.Undefined(), r.vm.ToValue(url))
	if err != nil {
		return fmt.Errorf("goja render: %w", err)
	}

	res, err := resolvePromise(pres)
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	var result RenderResult
	if err := r.vm.ExportTo(res, &result); err != nil {
		return fmt.Errorf("parse render result: %w", err)
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

	if _, err := fmt.Fprint(w, body); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

func resolvePromise(v goja.Value) (goja.Value, error) {
	p, ok := v.Export().(*goja.Promise)
	if !ok {
		return nil, fmt.Errorf("not a promise: %v", v)
	}

	switch p.State() {
	case goja.PromiseStateFulfilled:
		return p.Result(), nil
	case goja.PromiseStateRejected:
		r := p.Result()
		if o, ok := r.(*goja.Object); ok {
			return nil, fmt.Errorf("rejected: %v", o.Get("stack"))
		} else {
			return nil, fmt.Errorf("rejected: %v", r.String())
		}
	default:
		return nil, fmt.Errorf("unexpected pending state of promise %v", p)
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
