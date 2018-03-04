package article

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	reform "gopkg.in/reform.v1"
)

// RecordSaveHandler record create handler
type RecordSaveHandler struct {
	DB *reform.DB
}

func (h *RecordSaveHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
	record := ArticleTable.NewRecord()
	if err := json.NewDecoder(ctx.R.Body).Decode(record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	var err error
	if record.HasPK() {
		err = h.DB.Update(record)
		if err == reform.ErrNoRows {
			err = nil
		}
	} else {
		err = h.DB.Insert(record)
	}

	if err != nil {
		return errors.Wrap(err, "save record")
	}

	if err := index(record); err != nil {
		return errors.Wrap(err, "save record")
	}

	return nil
}

// index updates sphinx index
func index(record reform.Struct) error {
	p := parser{
		a: record.(*Article),
	}
	if err := p.parse(); err != nil {
		return errors.Wrap(err, "parse article")
	}
	return nil
}
