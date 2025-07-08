package storage

import (
	"encoding/json"
	"fmt"
	"log"
)

// ScrollCallback callback type to be called while scrolling an index
type ScrollCallback func(rawhits []json.RawMessage) error

// Scroll scrolls the specified index
func Scroll(index string, reqbody map[string]interface{}, cb ScrollCallback) error {
	if reqbody == nil {
		reqbody = map[string]interface{}{}
	}

	if _, ok := reqbody["sort"]; !ok {
		reqbody["sort"] = []string{"SortKey"}
	}

	if _, ok := reqbody["size"]; !ok {
		reqbody["size"] = 100
	}

	type scrollbodyt struct {
		ScrollID string `json:"_scroll_id"`
		Hits     struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []json.RawMessage `json:"hits"`
		} `json:"hits"`
	}

	respbody := &scrollbodyt{}
	path := fmt.Sprintf("/%s/_search?scroll=1m", index)
	if err := Post(path, reqbody, respbody); err != nil {
		return fmt.Errorf("start scroll: %w", err)
	}
	var pscrollID *string
	pscrollID = &respbody.ScrollID

	defer func() {
		reqbody := map[string]interface{}{"scroll_id": *pscrollID}
		if err := Delete("/_search/scroll", reqbody, nil); err != nil {
			log.Printf("failed to delete scroll id: %v", err)
		}
	}()

	for len(respbody.Hits.Hits) > 0 {
		if err := cb(respbody.Hits.Hits); err != nil {
			return fmt.Errorf("callback error: %w", err)
		}

		reqbody := map[string]interface{}{
			"scroll":    "1m",
			"scroll_id": respbody.ScrollID,
		}
		respbody = &scrollbodyt{}
		if err := Post("/_search/scroll", reqbody, respbody); err != nil {
			return fmt.Errorf("advance scroll: %w", err)
		}
		pscrollID = &respbody.ScrollID
	}

	return nil
}
