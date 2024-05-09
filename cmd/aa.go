package cmd

import (
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
		Short: "aa [command] [args]\n\nAvailable Commands:\n  help        Help about any command\n  version     Print the version number of gerrard",
		Long:  `aa [command] [args]\n\nAvailable Commands:\n  help        Help about any command\n  version     Print the version number of gerrard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}

	aacmd.AddCommand(aaCreateRepo())

	return aacmd
}

func aaCreateRepo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-repo",
		Short: "create-repo [args]\n\nCreate a new repository",
		Long:  `create-repo [args]\n\nCreate a new repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadDefaultConfig(cmd.Context())
			if err != nil {
				return err
			}

			client := servicecatalog.NewFromConfig(cfg)

			idempotencyToken, err := uuid.NewV7()
			if err != nil {
				return err
			}

			_, err = client.ProvisionProduct(cmd.Context(), &servicecatalog.ProvisionProductInput{
				ProductId:                aws.String(PACKAGE_REPO_NONPROD_ID),
				ProvisionedProductName:   aws.String("test-repo"),
				ProvisioningArtifactName: aws.String("v1.0.1"),
				ProvisioningParameters: []scTypes.ProvisioningParameter{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-repo"),
					},
				},
				ProvisionToken: aws.String(idempotencyToken.String()),
			})
			if err != nil {
				return err
			}

			return nil
		},
		SilenceUsage: true,
	}

	return cmd
}
