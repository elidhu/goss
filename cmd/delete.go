package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var (
	// For flags
	names []string

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete parameters from SSM",
		Run: func(cmd *cobra.Command, args []string) {
			deleteParameters(&names)
		},
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&names, "names", "n", []string{}, "")
	deleteCmd.MarkFlagRequired("names")
}

func deleteParameters(names *[]string) error {
	// Create Session
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("Session error: %w", err)
	}

	// We need to convert *[]string -> []*string
	var sliceNames []*string
	for _, v := range *names {
		sliceNames = append(sliceNames, &v)
	}

	// Create SSM service
	svc := ssm.New(sess)

	// Retrieve parameters
	_, err = svc.DeleteParameters(
		&ssm.DeleteParametersInput{
			Names: sliceNames,
		},
	)
	if err != nil {
		return fmt.Errorf("SSM request error: %w", err)
	}

	return nil
}
