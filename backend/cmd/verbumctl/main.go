package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/ctl"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(errors.Wrap(err, "app bootstrap"))
	}

	rootCmd := &cobra.Command{
		Use:   "verbumctl",
		Short: "verbumctl short",
		Long:  "verbumctl long",
	}
	rootCmd.AddCommand(
		ctl.Slugs(),
		ctl.RvblrWrongHeadwords(),
		dictimport.Command(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
