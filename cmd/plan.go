package main

import (
	"github.com/ayinke-llc/malak/config"
	"github.com/spf13/cobra"
)

func addPlanCommand(c *cobra.Command, cfg *config.Config) {

	cmd := &cobra.Command{
		Use: "plan",
	}

	cmd.AddCommand(listPlans(c, cfg))
	cmd.AddCommand(createPlan(c, cfg))
	cmd.AddCommand(setDefaultPlan(c, cfg))

	c.AddCommand(cmd)
}

func listPlans(c *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: `List all available plans in the system`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func createPlan(c *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: `Create a new plan`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func setDefaultPlan(c *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set-default",
		Short: `Set your default plan`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
