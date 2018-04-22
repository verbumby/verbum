package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/verbumby/verbum/backend/pkg/app"
)

// DictData dict data
type DictData struct {
	Items ItemsData `xml:"items"`
}

// ItemsData items data
type ItemsData struct {
	Item []ItemData `xml:"item"`
}

// ItemData item data
type ItemData struct {
	Title string `xml:"title"`
	Meta  string `xml:"meta"`
	Def   string `xml:"definition"`
}

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use:   "rvblr",
		Short: "rvblr short",
		Long:  "rvblr long",
	}
	rootCmd.AddCommand(rvblrImportCmd, rvblrNumlistsFixCmd, regexReplaceCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
