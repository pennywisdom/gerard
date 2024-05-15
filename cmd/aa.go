package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	scTypes "github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

const (
	PACKAGE_REPO_NONPROD_ID          = "prod-wkmlzxp6ly6b6"
	CONTAINER_IMAGES_REPO_NONPROD_ID = "prod-2jaq3qdizlarc"
)

func aaExecute() *cobra.Command {
	aacmd := &cobra.Command{
		Use:   "aa",
		Short: "aa [command] [args]\n\nAvailable Commands:\n  help        Help about any command\n  version     Print the version number of gerard",
		Long:  `aa [command] [args]\n\nAvailable Commands:\n  help        Help about any command\n  version     Print the version number of gerard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}
	aacmd.AddCommand(aaCreateRepo())

	return aacmd
}

func aaCreateRepo() *cobra.Command {
	input := new(aaCreateRepoInput)
	validRepoTypes := []string{"both", "dev", "prod"}
	cmd := &cobra.Command{
		Use:   "create-repo",
		Short: "Create a new repository",
		Long:  `create-repo [args]\n\nCreate a new repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			isValidRepoType := false
			for _, repoType := range validRepoTypes {
				if repoType == input.repoType {
					isValidRepoType = true
					break
				}
			}
			if !isValidRepoType {
				log.Panicf("invalid repo type %q, must be one of %v", input.repoType, validRepoTypes)
			}

			cfg, err := config.LoadDefaultConfig(cmd.Context())
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

			log.Printf("inputs: %v", input)

			output, err := client.ProvisionProduct(cmd.Context(), &servicecatalog.ProvisionProductInput{
				ProductId:                aws.String(PACKAGE_REPO_NONPROD_ID),
				ProvisionedProductName:   aws.String(fmt.Sprintf("%s-%s", input.repoType, suffix)),
				ProvisioningArtifactName: aws.String("v1.0.1"),
				ProvisioningParameters: []scTypes.ProvisioningParameter{
					{
						Key:   aws.String("CreateRepositoryType"),
						Value: aws.String(input.repoType),
					},
					{
						Key:   aws.String("Product"),
						Value: aws.String(input.product),
					},
					{
						Key:   aws.String("Bu"),
						Value: aws.String(input.businessUnit),
					},
					{
						Key:   aws.String("Div"),
						Value: aws.String(input.division),
					},
					{
						Key:   aws.String("Proj"),
						Value: aws.String(input.project),
					},
				},
				ProvisionToken: aws.String(idempotencyToken.String()),
			})
			if err != nil {
				return err
			}

			status := output.RecordDetail.Status
			log.Println("status: ", status)

			for status == scTypes.RecordStatusInProgress || status == scTypes.RecordStatusCreated {
				fmt.Printf(">")
				record, err := client.DescribeRecord(cmd.Context(), &servicecatalog.DescribeRecordInput{
					Id: output.RecordDetail.RecordId,
				})
				if err != nil {
					return err
				}
				status = record.RecordDetail.Status
				time.Sleep(5 * time.Second)
				// log.Println("status: ", status)
			}

			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&input.repoType, "repo-type", "r", "dev", fmt.Sprintf("Repository type (one of - %v)", validRepoTypes))
	cmd.Flags().StringVarP(&input.product, "product", "p", "", "Product name")
	cmd.Flags().StringVarP(&input.businessUnit, "business-unit", "b", "", "Business unit")
	cmd.Flags().StringVarP(&input.division, "division", "d", "", "Division")
	cmd.Flags().StringVarP(&input.project, "project", "j", "", "Project")

	//nolint:errcheck
	cmd.MarkFlagRequired("product")
	//nolint:errcheck
	cmd.MarkFlagRequired("business-unit")
	//nolint:errcheck
	cmd.MarkFlagRequired("division")
	//nolint:errcheck
	cmd.MarkFlagRequired("project")

	return cmd
}
