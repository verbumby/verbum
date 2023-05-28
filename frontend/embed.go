package frontend

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed dist
var Dist embed.FS

var DistPublic fs.FS

func init() {
	var err error
	DistPublic, err = fs.Sub(Dist, "dist/public")
	if err != nil {
		panic(fmt.Errorf("initalize dist public sub fs: %w", err))
	}
}
