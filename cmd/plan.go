package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/charmbracelet/huh"
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

func createPlan(_ *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: `Create a new plan`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				amount        string
				name          string
				reference     string
				defaultPlanID string
				isDefaultPlan bool
				teamCount     int
			)

			logger, err := getLogger(hermes.DeRef(cfg))
			if err != nil {
				return err
			}

			logger.Debug("creating a new plan")

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("What’s the name of your plan?").
						Value(&name).
						Validate(func(str string) error {
							if hermes.IsStringEmpty(str) {
								return errors.New("please provide the name of your plan")
							}

							if len(str) > 20 {
								return errors.New("plan name cannot be more than 20")
							}

							return nil
						}),

					huh.NewInput().
						Title("What is your reference for this plan").
						Description(`If you use stripe, it starts from prod_. If you do not use Stripe, just leave it empty.
						We will generate one for you`).
						CharLimit(100).
						Value(&reference),

					huh.NewInput().
						Title("What is your default plan id for this plan").
						Description(`If you use stripe, it starts from price_, we will automatically move the user to this 
						pricing structure 
						if they select the plan.

						If you do not use Stripe, just leave it empty. We will generate one for you`).
						CharLimit(100).
						Value(&defaultPlanID),

					huh.NewConfirm().
						Title("Default plan").
						Description("should this plan be the default one for new users?").
						Value(&isDefaultPlan),
				),

				huh.NewGroup(
					huh.NewInput().
						Title("how much does this plan cost monthly?").
						Description("Input in dollars not cents").
						Value(&amount),

					huh.NewSelect[int]().
						Title("Number of team mates that can be in this workspace").
						Options(
							huh.NewOption("1", 1),
							huh.NewOption("3", 3),
							huh.NewOption("5", 5),
							huh.NewOption("10", 10),
							huh.NewOption("15", 15),
							huh.NewOption("20", 20),
						).
						Value(&teamCount),
				),
			)

			refGen := malak.NewReferenceGenerator()

			ref := reference
			if hermes.IsStringEmpty(ref) {
				ref = refGen.Generate(malak.EntityTypePlan).String()
			}

			var defaultPlanIDRef = defaultPlanID
			if hermes.IsStringEmpty(defaultPlanIDRef) {
				defaultPlanIDRef = refGen.Generate(malak.EntityTypePrice).String()
			}

			amountinCents, err := strconv.Atoi(amount)
			if err != nil {
				logger.Error("could not generate amount", zap.Error(err), zap.String("amount", amount))
				return err
			}

			defer db.Close()

			if err := form.Run(); err != nil {
				logger.Error("could not take in data from user",
					zap.Error(err))
				return err
			}

			plan := &malak.Plan{
				PlanName:       name,
				Reference:      reference,
				DefaultPriceID: defaultPlanIDRef,
				Amount:         int64(amountinCents) * 100,
				IsDefault:      isDefaultPlan,
				Metadata: malak.PlanMetadata{
					Team: struct {
						Size    malak.Counter "json:\"size,omitempty\""
						Enabled bool          "json:\"enabled,omitempty\""
					}{
						Size:    malak.Counter(teamCount),
						Enabled: true,
					},
				},
			}

			planRepo := postgres.NewPlanRepository(db)

			if err := planRepo.Create(context.Background(), plan); err != nil {
				logger.Error("could not create plan", zap.Error(err))
				return err
			}

			if isDefaultPlan {
				if err := planRepo.SetDefault(context.Background(), plan); err != nil {
					logger.Error("could not set plan as default", zap.Error(err))
					return err
				}
			}

			return nil
		},
	}

	return cmd
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
