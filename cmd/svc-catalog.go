package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	scTypes "github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	view "github.com/pennywisdom/gerard/view"
	"github.com/spf13/cobra"
)

var (
	PACKAGE_REPO_NONPROD_ID          = os.Getenv("PACKAGE_REPO_PROD_ID")          //          = "prod-wkmlzxp6ly6b6"
	CONTAINER_IMAGES_REPO_NONPROD_ID = os.Getenv("CONTAINER_IMAGES_REPO_PROD_ID") //= "prod-2jaq3qdizlarc"
)

func svcCatExecute(ctx context.Context) *cobra.Command {
	svcCatCmd := &cobra.Command{
		Use:   "svc-catalog",
		Short: "Your personal assistant for AWS Service Catalog",
		Long:  `svc-catalog [command]\n\nYour personal assistant for AWS Service Catalog`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}
	svcCatCmd.AddCommand(svcCatProvisionProduct(ctx))

	return svcCatCmd
}

func svcCatProvisionProduct(ctx context.Context) *cobra.Command {
	input := new(svcCatProvisionProductInput)

	cmd := &cobra.Command{
		Use:          "provision-product",
		Short:        "Provision a product in AWS Service Catalog",
		Long:         `provision-product [flags]\n\nProvision a product in AWS Service Catalog`,
		RunE:         uiProvisionProduct(ctx, input),
		SilenceUsage: true,
	}

	cmd.Flags().StringToStringVarP(&input.vars, "var", "v", map[string]string{}, "Parameters for the product")
	cmd.Flags().StringVarP(&input.productId, "product-id", "p", "", "Product ID")
	//nolint:errcheck
	cmd.MarkFlagRequired("product-id")

	return cmd
}

func uiProvisionProduct(ctx context.Context, input *svcCatProvisionProductInput) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		model, err := view.NewModel()
		if err != nil {
			return err
		}
		// model.F = provisionProduct
		model.F = func() error {
			return provisionProduct(ctx, input)
		}
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	}
}

func provisionProduct(ctx context.Context, input *svcCatProvisionProductInput) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := servicecatalog.NewFromConfig(cfg)

	idempotencyToken, err := uuid.NewV7()
	if err != nil {
		return err
	}

	suffix, err := generateRandom10Char()
	if err != nil {
		return err
	}

	// log.Printf("inputs: %v", input)

	params := []scTypes.ProvisioningParameter{}
	for k, v := range input.vars {
		params = append(params, scTypes.ProvisioningParameter{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	output, err := client.ProvisionProduct(ctx, &servicecatalog.ProvisionProductInput{
		ProductId:                aws.String(input.productId),
		ProvisionedProductName:   aws.String(fmt.Sprintf("%s-%s", input.productId, suffix)),
		ProvisioningArtifactName: aws.String("v1.0.1"),
		ProvisioningParameters:   params,
		ProvisionToken:           aws.String(idempotencyToken.String()),
	})
	if err != nil {
		return err
	}

	status := output.RecordDetail.Status
	// log.Println("status: ", status)

	for status == scTypes.RecordStatusInProgress || status == scTypes.RecordStatusCreated {
		fmt.Printf(">")
		record, err := client.DescribeRecord(ctx, &servicecatalog.DescribeRecordInput{
			Id: output.RecordDetail.RecordId,
		})
		if err != nil {
			return err
		}
		status = record.RecordDetail.Status
		time.Sleep(5 * time.Second)
		// log.Println("status: ", status) ,
	}
	// log.Println("status: ", status)

	return nil

}
