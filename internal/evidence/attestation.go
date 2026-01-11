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

package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/miqcie/mondrian/internal/policy"
)

// Attestation represents a signed statement about policy check results
type Attestation struct {
	// SLSA/in-toto compatible structure
	PredicateType string                 `json:"predicateType"`
	Subject       []Subject              `json:"subject"`
	Predicate     PolicyCheckPredicate   `json:"predicate"`
	
	// Mondrian-specific metadata
	Timestamp     time.Time              `json:"timestamp"`
	RunID         string                 `json:"runId"`
	ParentHash    string                 `json:"parentHash,omitempty"`
	Hash          string                 `json:"hash"`
}

type Subject struct {
	Name   string            `json:"name"`
	Digest map[string]string `json:"digest"`
}

type PolicyCheckPredicate struct {
	// Core policy check information
	Results       []policy.CheckResult `json:"results"`
	Summary       Summary              `json:"summary"`
	
	// Environment context
	Repository    string               `json:"repository,omitempty"`
	Branch        string               `json:"branch,omitempty"`
	Commit        string               `json:"commit,omitempty"`
	Workflow      string               `json:"workflow,omitempty"`
	
	// Execution metadata
	Scanner       ScannerInfo          `json:"scanner"`
	FilesScanned  []string             `json:"filesScanned"`
}

type Summary struct {
	TotalChecks   int `json:"totalChecks"`
	Passed        int `json:"passed"`
	Failed        int `json:"failed"`
	Warnings      int `json:"warnings"`
	OverallStatus string `json:"overallStatus"` // "pass", "fail", "warn"
}

type ScannerInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	RulesUsed []string `json:"rulesUsed"`
}

// NewAttestation creates a new attestation from policy check results
func NewAttestation(results []policy.CheckResult, metadata AttestationMetadata) *Attestation {
	summary := calculateSummary(results)
	
	// Generate subjects from scanned files
	subjects := make([]Subject, len(metadata.FilesScanned))
	for i, file := range metadata.FilesScanned {
		hash := sha256.Sum256([]byte(file)) // Simple file path hash for now
		subjects[i] = Subject{
			Name: file,
			Digest: map[string]string{
				"sha256": hex.EncodeToString(hash[:]),
			},
		}
	}
	
	predicate := PolicyCheckPredicate{
		Results:      results,
		Summary:      summary,
		Repository:   metadata.Repository,
		Branch:       metadata.Branch,
		Commit:       metadata.Commit,
		Workflow:     metadata.Workflow,
		Scanner: ScannerInfo{
			Name:      "mondrian",
			Version:   "v0.1.0",
			RulesUsed: metadata.RulesUsed,
		},
		FilesScanned: metadata.FilesScanned,
	}
	
	attestation := &Attestation{
		PredicateType: "https://mondrian.dev/policy-check/v0.1",
		Subject:       subjects,
		Predicate:     predicate,
		Timestamp:     time.Now().UTC(),
		RunID:         generateRunID(),
		ParentHash:    metadata.ParentHash,
	}
	
	// Calculate hash of the complete attestation
	attestation.Hash = attestation.calculateHash()
	
	return attestation
}

type AttestationMetadata struct {
	Repository   string
	Branch       string
	Commit       string
	Workflow     string
	FilesScanned []string
	RulesUsed    []string
	ParentHash   string
}

// calculateSummary generates summary statistics from policy check results
func calculateSummary(results []policy.CheckResult) Summary {
	summary := Summary{
		TotalChecks: len(results),
	}
	
	for _, result := range results {
		switch result.Status {
		case "pass":
			summary.Passed++
		case "fail":
			summary.Failed++
		case "warn":
			summary.Warnings++
		}
	}
	
	// Determine overall status
	if summary.Failed > 0 {
		summary.OverallStatus = "fail"
	} else if summary.Warnings > 0 {
		summary.OverallStatus = "warn"
	} else {
		summary.OverallStatus = "pass"
	}
	
	return summary
}

// calculateHash generates a SHA-256 hash of the attestation content
func (a *Attestation) calculateHash() string {
	// Create a copy without the hash field for hashing
	temp := *a
	temp.Hash = ""
	
	data, err := json.Marshal(temp)
	if err != nil {
		// Fallback to timestamp-based hash
		data = []byte(fmt.Sprintf("%s-%s", a.RunID, a.Timestamp.Format(time.RFC3339)))
	}
	
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// ToJSON serializes the attestation to JSON
func (a *Attestation) ToJSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

// generateRunID creates a unique run identifier
func generateRunID() string {
	timestamp := time.Now().Unix()
	hash := sha256.Sum256([]byte(fmt.Sprintf("mondrian-%d", timestamp)))
	return hex.EncodeToString(hash[:8]) // First 8 bytes for shorter ID
}