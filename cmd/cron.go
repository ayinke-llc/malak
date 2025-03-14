package main

import (
	"github.com/ayinke-llc/malak/config"
	"github.com/spf13/cobra"
)

func addCronCommand(c *cobra.Command, cfg *config.Config) {

	cmd := &cobra.Command{
		Use: "cron",
	}

	cmd.AddCommand(sendScheduledUpdates(c, cfg))
	cmd.AddCommand(processDeckAnalytics(c, cfg))

	c.AddCommand(cmd)
}
