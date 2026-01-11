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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/miqcie/mondrian/internal/evidence"
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
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	
	// Run policy checks first to get results
	scanner := policy.NewFileScanner(wd)
	files, err := scanner.ScanRelevantFiles()
	if err != nil {
		fmt.Printf("❌ Error scanning files: %v\n", err)
		os.Exit(1)
	}
	
	if len(files) == 0 {
		fmt.Println("ℹ️  No relevant files found for attestation")
		return
	}
	
	fmt.Printf("📝 Generating attestation for %d files...\n", len(files))
	
	// Run policy checks
	engine := policy.NewPolicyEngine()
	results := engine.RunChecks(files)
	
	// Create evidence directory
	evidenceDir := filepath.Join(wd, ".mondrian", "attestations")
	
	// Initialize chain manager
	chainManager := evidence.NewChainManager(evidenceDir)
	chain, err := chainManager.LoadOrCreateChain()
	if err != nil {
		fmt.Printf("❌ Error loading evidence chain: %v\n", err)
		os.Exit(1)
	}
	
	// Get file list
	var fileList []string
	for filename := range files {
		fileList = append(fileList, filename)
	}
	
	// Get rule names
	ruleNames := make([]string, len(engine.Rules))
	for i, rule := range engine.Rules {
		ruleNames[i] = rule.Name()
	}
	
	// Create attestation metadata
	metadata := evidence.AttestationMetadata{
		Repository:   getRepositoryName(wd),
		Branch:       getBranchName(),
		Commit:       getCommitHash(),
		Workflow:     getWorkflowContext(),
		FilesScanned: fileList,
		RulesUsed:    ruleNames,
		ParentHash:   chain.Head, // Will be updated by chain manager
	}
	
	// Create attestation
	attestation := evidence.NewAttestation(results, metadata)
	
	// Create signer
	signer, err := evidence.NewSigner()
	if err != nil {
		fmt.Printf("❌ Error creating signer: %v\n", err)
		os.Exit(1)
	}
	
	// Sign attestation
	signed, err := signer.SignAttestation(attestation)
	if err != nil {
		fmt.Printf("❌ Error signing attestation: %v\n", err)
		os.Exit(1)
	}
	
	// Save signed attestation
	timestamp := signed.Metadata.Timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("attestation-%s-%s.json", timestamp, signed.Metadata.KeyID[:8])
	filePath := filename
	
	if err := evidence.SaveSignedAttestation(signed, evidenceDir); err != nil {
		fmt.Printf("❌ Error saving attestation: %v\n", err)
		os.Exit(1)
	}
	
	// Add to evidence chain
	if err := chainManager.AddAttestation(chain, attestation, filePath); err != nil {
		fmt.Printf("❌ Error adding to evidence chain: %v\n", err)
		os.Exit(1)
	}
	
	// Display results
	fmt.Printf("✅ Attestation generated and signed\n")
	fmt.Printf("📁 Evidence directory: %s\n", evidenceDir)
	fmt.Printf("🔑 Key ID: %s\n", signed.Metadata.KeyID[:16])
	fmt.Printf("🔗 Chain length: %d attestations\n", chain.Length+1)
	fmt.Printf("📊 Status: %s (%d checks)\n", attestation.Predicate.Summary.OverallStatus, attestation.Predicate.Summary.TotalChecks)
}

func verifyEvidence() {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	
	// Evidence directory
	evidenceDir := filepath.Join(wd, ".mondrian", "attestations")
	
	// Check if evidence directory exists
	if _, err := os.Stat(evidenceDir); os.IsNotExist(err) {
		fmt.Println("❌ No evidence directory found")
		fmt.Printf("💡 Run 'mondrian attest' to generate attestations first\n")
		os.Exit(1)
	}
	
	fmt.Println("🔍 Verifying evidence chain...")
	
	// Initialize chain manager
	chainManager := evidence.NewChainManager(evidenceDir)
	
	// Load existing chain
	chain, err := chainManager.LoadOrCreateChain()
	if err != nil {
		fmt.Printf("❌ Error loading evidence chain: %v\n", err)
		os.Exit(1)
	}
	
	if chain.Length == 0 {
		fmt.Println("ℹ️  No attestations found in evidence chain")
		return
	}
	
	// Verify chain integrity
	fmt.Printf("🔗 Verifying chain integrity (%d attestations)...\n", chain.Length)
	if err := chainManager.VerifyChain(chain); err != nil {
		fmt.Printf("❌ Chain verification failed: %v\n", err)
		os.Exit(1)
	}
	
	// Display chain summary
	fmt.Println("✅ Evidence chain verification passed!")
	fmt.Println()
	fmt.Println("📊 Chain Summary:")
	fmt.Println(chain.GetChainSummary())
	fmt.Println()
	fmt.Printf("🔑 Chain ID: %s\n", chain.ChainID)
	fmt.Printf("🏷️  Genesis Hash: %s\n", chain.Genesis[:16]+"...")
	fmt.Printf("🔝 Head Hash: %s\n", chain.Head[:16]+"...")
	
	// Show recent attestations
	fmt.Println()
	fmt.Println("📜 Recent Attestations:")
	start := 0
	if chain.Length > 5 {
		start = chain.Length - 5
		fmt.Printf("   (showing last 5 of %d)\n", chain.Length)
	}
	
	for i := start; i < chain.Length; i++ {
		entry := chain.Attestations[i]
		status := "✅"
		if entry.Status == "fail" {
			status = "❌"
		} else if entry.Status == "warn" {
			status = "⚠️"
		}
		
		fmt.Printf("   %s %s [%s] %s\n", 
			status, 
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Status,
			entry.RunID[:8]+"...")
	}
	
	fmt.Println()
	fmt.Printf("🎯 Verification complete - evidence chain is valid and tamper-evident\n")
}

func initializeProject() {
	fmt.Println("⚠️  Project initialization implementation coming soon...")
}

func startServer() {
	fmt.Println("⚠️  Evidence viewer implementation coming soon...")
}

// Helper functions for gathering context information

func getRepositoryName(wd string) string {
	// Try to get from git remote
	if output, err := runCommand("git", "remote", "get-url", "origin"); err == nil {
		return string(output)
	}
	
	// Fallback to directory name
	return filepath.Base(wd)
}

func getBranchName() string {
	if output, err := runCommand("git", "branch", "--show-current"); err == nil {
		branch := strings.TrimSpace(string(output))
		if branch != "" {
			return branch
		}
	}
	return "unknown"
}

func getCommitHash() string {
	if output, err := runCommand("git", "rev-parse", "HEAD"); err == nil {
		hash := strings.TrimSpace(string(output))
		if len(hash) >= 12 {
			return hash[:12] // Short hash
		}
	}
	return "unknown"
}

func getWorkflowContext() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		workflow := os.Getenv("GITHUB_WORKFLOW")
		runNumber := os.Getenv("GITHUB_RUN_NUMBER")
		if workflow != "" && runNumber != "" {
			return fmt.Sprintf("%s #%s", workflow, runNumber)
		}
		return "github-actions"
	}
	return "local"
}

func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}