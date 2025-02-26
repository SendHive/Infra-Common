package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	testsuite "github.com/SendHive/Infra-Common/test-suite"
	"github.com/spf13/cobra"
)

func runTestFunction() {
	fmt.Println("Entered the Test....")
	time.Sleep(15 * time.Second)
	fmt.Println("Running test function...")
	testsuite.TestSuite()
}

func main() {
	var testFlag bool

	rootCmd := &cobra.Command{
		Use:   "infra-common",
		Short: "A CLI that runs a test function or executes a shell script",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running setup.sh to bring up RabbitMQ, Minio, and DB...")
			setupCmd := exec.Command("sh", "setup.sh")
			setupCmd.Stdout = os.Stdout
			setupCmd.Stderr = os.Stderr
			err := setupCmd.Run()
			if err != nil {
				fmt.Println("Error running setup.sh:", err)
				os.Exit(1)
			}

			// After setup, check if testFlag is set
			if testFlag {
				// Run the test function if the flag is set
				runTestFunction()
			}
		},
	}

	// Define the --test flag for the command
	rootCmd.Flags().BoolVar(&testFlag, "test", false, "Run the test function instead of just setup.sh")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
