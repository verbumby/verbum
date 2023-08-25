package ctl

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
)

func CleanupCommand() *cobra.Command {
	c := &cleanupController{}
	result := &cobra.Command{
		Use: "cleanup",
		Run: c.Run,
	}

	result.PersistentFlags().BoolVar(&c.dryrun, "dryrun", true, "--dryrun=false to disable dryrun")

	return result
}

type cleanupController struct {
	dryrun bool
}

func (c *cleanupController) Run(cmd *cobra.Command, args []string) {
	if c.dryrun {
		log.Println("running in dryrun mode")
	}
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}

func (c *cleanupController) run() error {
	todelete := []string{}
	list := map[string]struct {
		Aliases map[string]any `json:"aliases"`
	}{}
	if err := storage.Get("/dict-*,sugg-*", &list); err != nil {
		return fmt.Errorf("list indices: %w", err)
	}
	for index, settings := range list {
		if dictionary.GetByIndex(index) == nil {
			todelete = append(todelete, index)
		}
		if len(settings.Aliases) == 0 {
			todelete = append(todelete, index)
		}
	}

	for _, index := range todelete {
		if !c.dryrun {
			if err := storage.Delete("/"+index, nil, nil); err != nil {
				return fmt.Errorf("deleting %s index: %w", index, err)
			}
		}
		log.Printf("%s index deleted", index)
	}

	if len(todelete) == 0 {
		log.Println("nothing to delete")
	}

	return nil
}
