package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
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

func setDefaultPlan(_ *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set-default",
		Short: `Set your default plan`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := getLogger(hermes.DeRef(cfg))
			if err != nil {
				return err
			}

			planReference, err := cmd.Flags().GetString("reference")
			if err != nil {
				logger.Error("please provide the reference flag", zap.Error(err))
				return err
			}

			logger = logger.With(
				zap.String("plan_reference", planReference))

			logger.Debug("setting plan as default")

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			defer db.Close()

			planRepository := postgres.NewPlanRepository(db)

			plan, err := planRepository.Get(context.Background(), &malak.FetchPlanOptions{
				Reference: planReference,
			})
			if err != nil {
				logger.Error("could not list plan from the database",
					zap.Error(err))
				return err
			}

			if err := planRepository.SetDefault(context.Background(), plan); err != nil {
				logger.Error("could not set default plan",
					zap.Error(err))
				return err
			}

			logger.Debug("successfully set plan as default")
			return nil
		},
	}

	cmd.Flags().String("reference", "",
		`reference of the plan you want to make the default for all users`)

	return cmd
}
