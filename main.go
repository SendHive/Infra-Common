package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/SendHive/Infra-Common/setup"
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
			InfraCmd := exec.Command("sh", "setup.sh")
			InfraCmd.Stdout = os.Stdout
			InfraCmd.Stderr = os.Stderr
			err := InfraCmd.Run()
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

	var delCmd = &cobra.Command{
		Use:   "del",
		Short: "Run docker-compose down",
		Run: func(cmd *cobra.Command, args []string) {
			// Create the command to run docker-compose down
			command := exec.Command("docker-compose", "down")
			output, err := command.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running docker-compose down: %s\n", err)
				return
			}
			fmt.Println(string(output))
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Run docker-compose down",
		Run: func(cmd *cobra.Command, args []string) {
			// Create the command to run docker-compose down
			command := exec.Command("docker-compose", "stop")
			output, err := command.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running docker-compose down: %s\n", err)
				return
			}
			fmt.Println(string(output))
		},
	}

	var setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Run to setup the infrastructure",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running setup.sh to bring up RabbitMQ, Minio, and DB...")
			InfraCmd := exec.Command("sh", "setup.sh")
			InfraCmd.Stdout = os.Stdout
			InfraCmd.Stderr = os.Stderr
			err := InfraCmd.Run()
			if err != nil {
				fmt.Println("Error running setup.sh:", err)
				os.Exit(1)
			}
			time.Sleep(10 * time.Second)
			fmt.Println("Starting the setup for the infra")
			setup.Setup()
		},
	}

	// Define the --test flag for the command
	rootCmd.Flags().BoolVar(&testFlag, "test", false, "Run the test function instead of just setup.sh")
	rootCmd.AddCommand(delCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(setupCmd)
	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
