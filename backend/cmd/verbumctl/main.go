package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/ctl"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(fmt.Errorf("app bootstrap: %w", err))
	}

	rootCmd := &cobra.Command{
		Use:   "verbumctl",
		Short: "verbumctl",
		Long:  "verbumctl",
	}
	rootCmd.AddCommand(
		dictimport.Command(),
		ctl.MigrateSlugs(),
		ctl.MigrateRvblrWrongHeadwords(),
		ctl.MigrateStardictFixPrefixCase(),
		ctl.MigrateModifiedAt(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
