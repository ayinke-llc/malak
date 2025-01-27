package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

func listPlans(_ *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: `List all available plans in the system`,
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, err := getLogger(hermes.DeRef(cfg))
			if err != nil {
				return err
			}

			logger.Debug("listing all plans")

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			defer db.Close()

			planRepository := postgres.NewPlanRepository(db)

			plans, err := planRepository.List(context.Background())
			if err != nil {
				logger.Error("could not list plans from the database",
					zap.Error(err))
				return err
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Amount", "Is Default", "Reference"})

			for _, plan := range plans {
				table.Append([]string{
					plan.ID.String(),
					plan.PlanName,
					fmt.Sprintf("$%d", plan.Amount/100),
					fmt.Sprintf("%v", plan.IsDefault),
					plan.Reference,
				})
			}

			table.Render()

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
