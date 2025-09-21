package storage

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BulkResponseItem struct {
	ID    string          `json:"_id"`
	Error json.RawMessage `json:"error"`
}

type BulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Create *BulkResponseItem `json:"create"`
		Index  *BulkResponseItem `json:"index"`
		Delete *BulkResponseItem `json:"delete"`
		Update *BulkResponseItem `json:"update"`
	} `json:"items"`
}

func (resp BulkResponse) Error() error {
	if !resp.Errors {
		return nil
	}

	errors := []string{}
	for _, item := range resp.Items {
		for _, op := range []*BulkResponseItem{item.Create, item.Index, item.Delete, item.Update} {
			if op != nil && op.Error != nil {
				errors = append(errors, fmt.Sprintf("id `%s`: %s", op.ID, string(op.Error)))
			}
		}
	}
	return fmt.Errorf("bulk post failed: %s", strings.Join(errors, "; "))
}
