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
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/cobra"
)

var (
	// For flags.
	envPath string

	// envCmd represents the env command.
	envCmd = &cobra.Command{
		Use:   "env",
		Short: "Load parameters into the environment and run a command",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("env called")
			envRun(args[0], args[1:], envPath)
		},
	}
)

func init() {
	rootCmd.AddCommand(envCmd)

	envCmd.Flags().StringVarP(
		&envPath, "path", "p", "", "parameter path",
	)
	envCmd.MarkFlagRequired("path")

}

func envRun(command string, args []string, path string) error {
	fmt.Printf("Command: %s\nArgs: %s\nPath: %s\n", command, args, path)

	// Fetch the parameters.
	env := environ(os.Environ())
	mp, err := getParameters(path)
	if err != nil {
		return fmt.Errorf("Failed to get parameters: %w", err)
	}

	// Set the parameters in the environment.
	for k, v := range mp {
		env.Set(k, v)
	}

	if !supportsExecSyscall() {
		return execCmd(command, args, env)
	}

	return execSyscall(command, args, env)
}

// getParameters is a virtual copy of listParameters - it needs to be refactored
func getParameters(
	path string,
) (map[string]string, error) {
	// Create Session.
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Session error: %w", err)
	}

	// Create SSM service.
	svc := ssm.New(sess)

	// Retrieve parameters.
	res, err := svc.GetParametersByPath(
		&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("SSM request error: %w", err)
	}

	// Put the variables into a k-v map.
	mp := make(map[string]string)
	for _, v := range res.Parameters {
		ss := strings.Split(*v.Name, "/")
		key := ss[len(ss)-1]
		if key != "" {
			mp[key] = *v.Value
		}
	}

	return mp, nil
}

func supportsExecSyscall() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd"
}

func execCmd(command string, args []string, env []string) error {
	log.Printf("Starting child process: %s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		for {
			sig := <-sigChan
			cmd.Process.Signal(sig)
		}
	}()

	if err := cmd.Wait(); err != nil {
		cmd.Process.Signal(os.Kill)
		return fmt.Errorf("Failed to wait for command termination: %v", err)
	}

	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	os.Exit(waitStatus.ExitStatus())
	return nil
}

func execSyscall(command string, args []string, env []string) error {
	argv0, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("Couldn't find the executable '%s': %w", command, err)
	}

	argv := make([]string, 0, 1+len(args))
	argv = append(argv, command)
	argv = append(argv, args...)

	return syscall.Exec(argv0, argv, env)
}

// environ is a slice of strings representing the environment, in the form "key=value".
type environ []string

// Unset an environment variable by key
func (e *environ) Unset(key string) {
	for i := range *e {
		// If we found the key
		if strings.HasPrefix((*e)[i], key+"=") {
			// Move the last value to replace the key
			(*e)[i] = (*e)[len(*e)-1]
			// Slice of the last value as we moved it to 'i'
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

// Set adds an environment variable, replacing any existing ones of the same key
func (e *environ) Set(key, val string) {
	e.Unset(key)
	*e = append(*e, key+"="+val)
}
