package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/verbumby/verbum/backend/pkg/app"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use:   "rvblr",
		Short: "rvblr short",
		Long:  "rvblr long",
	}
	rootCmd.AddCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
