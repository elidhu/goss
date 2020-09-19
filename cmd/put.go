package cmd

import (
	"fmt"

	"github.com/kevinglasson/goss/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/cobra"
)

var (
	// For flags
	name      string
	value     string
	putType   string
	putOWrite bool

	putCmd = &cobra.Command{
		Use:   "put",
		Short: "Put a single parameter into SSM",
		Run: func(cmd *cobra.Command, args []string) {

			// Run and report errors.
			err := putParameter(name, value, putType, putOWrite)
			if err != nil {
				utils.PrintErrorAndExit(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.Flags().StringVarP(&name, "name", "n", "", "parameter name including path")
	putCmd.MarkFlagRequired("name")
	putCmd.Flags().StringVarP(&value, "value", "v", "", "parameter value")
	putCmd.MarkFlagRequired("value")
	putCmd.Flags().StringVarP(&putType, "type", "t", "", "parameter type")
	putCmd.MarkFlagRequired("type")
	putCmd.Flags().BoolVarP(
		&putOWrite, "overwrite", "o", false, "overwrite parameter if it exists",
	)
}

func putParameter(name string, value string, typ string, overwrite bool) error {
	// Create Session
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("Session error: %w", err)
	}

	// Create SSM service
	svc := ssm.New(sess)

	// Retrieve parameters
	_, err = svc.PutParameter(
		&ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String(value),
			Type:      aws.String(typ),
			Overwrite: aws.Bool(overwrite),
		},
	)
	if err != nil {
		return fmt.Errorf("SSM request error: %w", err)
	}

	return nil
}
