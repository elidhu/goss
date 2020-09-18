package cmd

import (
	"encoding/json"
	"fmt"
	"goss/internal/utils"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var (
	// For flags.
	recursive bool
	decrypt   bool

	listCmd = &cobra.Command{
		Use:       "list PATH",
		Short:     "List parameters in SSM by path",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"path"},
		Run: func(cmd *cobra.Command, args []string) {
			// Guaranteed to work because of ExactArgs.
			path := args[0]

			// Global flags.
			asJSON, err := cmd.Flags().GetBool("json")
			if err != nil {
				utils.PrintErrorAndExit(err)
			}

			// Run and report errors.
			err = listParameters(path, recursive, decrypt, asJSON)
			if err != nil {
				utils.PrintErrorAndExit(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "recurse into parameter path")
	listCmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "decrypt values")
}

func listParameters(
	path string,
	recursive bool,
	decrypt bool,
	asJSON bool,
) error {
	// Create Session
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("Session error: %w", err)
	}

	// Create SSM service
	svc := ssm.New(sess)

	// Retrieve parameters
	res, err := svc.GetParametersByPath(
		&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(recursive),
			WithDecryption: aws.Bool(decrypt),
		},
	)
	if err != nil {
		return fmt.Errorf("SSM request error: %w", err)
	}

	// Output
	if asJSON {
		out, err := json.MarshalIndent(res.Parameters, "", "  ")
		if err != nil {
			return fmt.Errorf("Marshalling response error: %w", err)
		}
		fmt.Println(string(out))
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Value"})

		for _, v := range res.Parameters {
			table.Append(
				[]string{
					utils.TruncateString(*v.Name, 35),
					utils.TruncateString(*v.Value, 35),
				},
			)
		}
		table.Render()
	}

	return nil
}
