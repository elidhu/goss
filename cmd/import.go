package cmd

import (
	"fmt"

	"github.com/kevinglasson/goss/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	kdotenv "github.com/knadh/koanf/parsers/dotenv"
	kjson "github.com/knadh/koanf/parsers/json"
	ktoml "github.com/knadh/koanf/parsers/toml"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
)

type unmarshaller interface {
	Unmarshal(b []byte) (map[string]interface{}, error)
}

var (
	// For flags
	file         string
	importType   string
	importPath   string
	importOWrite bool
	format       string

	// importCmd represents the import command
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import parameters from a file",
		Run: func(cmd *cobra.Command, args []string) {

			// Define the map of string to parsers
			parsers := map[string]unmarshaller{
				"json":   kjson.Parser(),
				"toml":   ktoml.Parser(),
				"yaml":   kyaml.Parser(),
				"dotenv": kdotenv.Parser(),
			}

			// Check the format arg is for a supported format
			err := isValidFormat(format)
			if err != nil {
				utils.PrintErrorAndExit(err)
			}

			parser := parsers[format]

			err = importParameters(file, importType, importPath, importOWrite, parser)
			if err != nil {
				utils.PrintErrorAndExit(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(
		&file, "file", "f", "", "path to the file to import",
	)
	importCmd.MarkFlagRequired("file")

	importCmd.Flags().StringVarP(
		&importType, "type", "t", "", "aws parameter type to use when importing",
	)
	importCmd.MarkFlagRequired("type")

	importCmd.Flags().StringVarP(
		&importPath, "path", "p", "", "base path to import the parameters to",
	)
	importCmd.MarkFlagRequired("path")

	importCmd.Flags().BoolVarP(
		&importOWrite, "overwrite", "o", false, "overwrite parameters if they exist, i.e. update them",
	)

	importCmd.Flags().StringVarP(
		&format, "format", "", "dotenv", "file format. [json toml yaml dotenv]",
	)
}

func importParameters(file string, typ string, path string, overwrite bool, parser unmarshaller) error {

	provider := kfile.Provider(file)

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
				Name:      aws.String(fmt.Sprintf("%s/%s", path, k)),
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

func isValidFormat(format string) error {
	// Supported formats
	supported := []string{"json", "toml", "yaml", "dotenv"}

	for _, v := range supported {
		if v == format {
			return nil
		}
	}
	return fmt.Errorf(
		"Unsupported input file format `%s`. Supported formats are: %v", format, supported,
	)
}
