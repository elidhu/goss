package cmd

import (
	"fmt"
	"goss/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	kdotenv "github.com/knadh/koanf/parsers/dotenv"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
)

var (
	// For flags
	file         string
	importType   string
	importOWrite bool

	// importCmd represents the import command
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a file into SSM at the given path",
		Run: func(cmd *cobra.Command, args []string) {
			err := importParameters(file, importType, importOWrite)
			if err != nil {
				utils.PrintErrorAndExit(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(
		&file, "file", "f", "", "absolute path to the file",
	)
	importCmd.MarkFlagRequired("file")
	importCmd.Flags().StringVarP(
		&importType, "type", "t", "", "type of the parameter",
	)
	importCmd.MarkFlagRequired("type")
	importCmd.Flags().BoolVarP(
		&importOWrite, "overwrite", "o", false, "overwrite parameters if they exist",
	)
}

func importParameters(file string, typ string, overwrite bool) error {

	provider := kfile.Provider(file)
	parser := kdotenv.Parser()

	// Get the raw bytest from the provider
	b, err := provider.ReadBytes()
	if err != nil {
		return err
	}

	// Parse the raw bytes
	mp, err := parser.Unmarshal(b)
	if err != nil {
		return err
	}

	// Create Session
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("Session error: %w", err)
	}

	// Create SSM service
	svc := ssm.New(sess)

	// utils.OutputAsJSON(mp)
	count := 0

	for k, v := range mp {
		// Retrieve parameters
		_, err = svc.PutParameter(
			&ssm.PutParameterInput{
				Name:      aws.String(k),
				Value:     aws.String(v.(string)),
				Type:      aws.String(typ),
				Overwrite: aws.Bool(overwrite),
			},
		)
		if err != nil {
			return fmt.Errorf(
				"Failed to put `%s` -> `%s`. SSM request error: %w", k, v, err,
			)
		}
		count++
	}

	fmt.Printf("Successfully imported %d / %d parameters", count, len(mp))

	return nil
}
