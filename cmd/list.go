package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kevinglasson/goss/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var (
	// For flags.
	path      string
	recursive bool
	decrypt   bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List parameters SSM",
		Run: func(cmd *cobra.Command, args []string) {
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

	listCmd.Flags().StringVarP(
		&path, "path", "p", "", "parameter(s) path",
	)
	listCmd.MarkFlagRequired("path")
	listCmd.Flags().BoolVarP(
		&recursive, "recursive", "r", false, "recurse into parameter path",
	)
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
		table.SetHeader([]string{"Name", "Value", "Version", "Last Mod"})

		for _, v := range res.Parameters {
			table.Append(
				[]string{
					utils.TruncateString(*v.Name, 35),
					utils.TruncateString(*v.Value, 35),
					utils.TruncateString(strconv.FormatInt(*v.Version, 10), 35),
					utils.TruncateString(
						v.LastModifiedDate.Format(time.RFC3339), 35,
					),
				},
			)
		}
		table.Render()
	}

	return nil
}
