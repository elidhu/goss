/*
The MIT License (MIT)

Copyright (c) 2020 Kevin Glasson
Copyright (c) 2015 99designs (substantial copy of 'aws-vault exec')

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("env called")
	},
}

func init() {
	// rootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func execEnvironment(input ExecCommandInput, config *vault.Config, creds *credentials.Credentials) error {
// 	val, err := creds.Get()
// 	if err != nil {
// 		return fmt.Errorf("Failed to get credentials for %s: %w", input.ProfileName, err)
// 	}

// 	env := environ(os.Environ())
// 	env = updateEnvForAwsVault(env, input.ProfileName, config.Region)

// 	log.Println("Setting subprocess env: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY")
// 	env.Set("AWS_ACCESS_KEY_ID", val.AccessKeyID)
// 	env.Set("AWS_SECRET_ACCESS_KEY", val.SecretAccessKey)

// 	if val.SessionToken != "" {
// 		log.Println("Setting subprocess env: AWS_SESSION_TOKEN, AWS_SECURITY_TOKEN")
// 		env.Set("AWS_SESSION_TOKEN", val.SessionToken)
// 		env.Set("AWS_SECURITY_TOKEN", val.SessionToken)
// 	}
// 	if expiration, err := creds.ExpiresAt(); err == nil {
// 		log.Println("Setting subprocess env: AWS_SESSION_EXPIRATION")
// 		env.Set("AWS_SESSION_EXPIRATION", iso8601.Format(expiration))
// 	}

// 	if !supportsExecSyscall() {
// 		return execCmd(input.Command, input.Args, env)
// 	}

// 	return execSyscall(input.Command, input.Args, env)
// }

// func execCmd(command string, args []string, env []string) error {
// 	log.Printf("Starting child process: %s %s", command, strings.Join(args, " "))

// 	cmd := exec.Command(command, args...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.Env = env

// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan)

// 	if err := cmd.Start(); err != nil {
// 		return err
// 	}

// 	go func() {
// 		for {
// 			sig := <-sigChan
// 			cmd.Process.Signal(sig)
// 		}
// 	}()

// 	if err := cmd.Wait(); err != nil {
// 		cmd.Process.Signal(os.Kill)
// 		return fmt.Errorf("Failed to wait for command termination: %v", err)
// 	}

// 	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
// 	os.Exit(waitStatus.ExitStatus())
// 	return nil
// }

// func execSyscall(command string, args []string, env []string) error {
// 	log.Printf("Exec command %s %s", command, strings.Join(args, " "))

// 	argv0, err := exec.LookPath(command)
// 	if err != nil {
// 		return fmt.Errorf("Couldn't find the executable '%s': %w", command, err)
// 	}

// 	log.Printf("Found executable %s", argv0)

// 	argv := make([]string, 0, 1+len(args))
// 	argv = append(argv, command)
// 	argv = append(argv, args...)

// 	return syscall.Exec(argv0, argv, env)
// }
