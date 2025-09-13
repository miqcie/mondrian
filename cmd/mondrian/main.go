/*
Copyright 2025 Chris McConnell

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/miqcie/mondrian/internal/policy"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mondrian",
	Short: "Evidence-first Zero Trust runtime for startups",
	Long: `Mondrian blocks risky changes before merge or deploy and proves controls run with tamper-evident records.

Complete documentation is available at https://github.com/miqcie/mondrian`,
	Version: "v0.1.0",
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run policy checks against current environment",
	Long:  `Check runs all configured policies against the current repository, infrastructure, and environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Running Mondrian policy checks...")
		runPolicyChecks()
	},
}

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "Generate signed attestation for current state",
	Long:  `Attest creates a signed attestation documenting the current state and policy check results.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📝 Generating attestation...")
		generateAttestation()
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify attestation chain and print proof",
	Long:  `Verify validates the attestation chain and prints a human-readable proof bundle.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("✅ Verifying evidence chain...")
		verifyEvidence()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Mondrian in current repository",
	Long:  `Init sets up Mondrian configuration and creates necessary directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Initializing Mondrian...")
		initializeProject()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start web server for evidence viewer",
	Long:  `Serve starts a local web server to view attestation chains and evidence bundles.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌐 Starting evidence viewer...")
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(attestCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(serveCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runPolicyChecks() {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	
	// Scan for relevant files
	scanner := policy.NewFileScanner(wd)
	files, err := scanner.ScanRelevantFiles()
	if err != nil {
		fmt.Printf("❌ Error scanning files: %v\n", err)
		os.Exit(1)
	}
	
	if len(files) == 0 {
		fmt.Println("ℹ️  No relevant files found (looking for .tf, .yml, .yaml files)")
		return
	}
	
	fmt.Printf("🔍 Scanning %d files for policy violations...\n", len(files))
	
	// Run policy checks
	engine := policy.NewPolicyEngine()
	results := engine.RunChecks(files)
	
	// Display results
	output := policy.FormatResults(results)
	fmt.Print(output)
	
	// Exit with error code if there are failures
	for _, result := range results {
		if result.Status == "fail" {
			os.Exit(1)
		}
	}
}

func generateAttestation() {
	fmt.Println("⚠️  Attestation generation implementation coming soon...")
}

func verifyEvidence() {
	fmt.Println("⚠️  Evidence verification implementation coming soon...")
}

func initializeProject() {
	fmt.Println("⚠️  Project initialization implementation coming soon...")
}

func startServer() {
	fmt.Println("⚠️  Evidence viewer implementation coming soon...")
}