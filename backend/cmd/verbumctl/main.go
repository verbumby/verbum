package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/ctl"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use:   "verbumctl",
		Short: "verbumctl short",
		Long:  "verbumctl long",
	}
	rootCmd.AddCommand(ctl.Slugs())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
