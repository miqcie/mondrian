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
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// EvidenceChain represents a hash-chained sequence of attestations
type EvidenceChain struct {
	ChainID     string              `json:"chainId"`
	StartTime   time.Time           `json:"startTime"`
	LastUpdated time.Time           `json:"lastUpdated"`
	Length      int                 `json:"length"`
	Head        string              `json:"head"`        // Hash of most recent attestation
	Genesis     string              `json:"genesis"`     // Hash of first attestation
	Attestations []ChainEntry       `json:"attestations"`
}

type ChainEntry struct {
	Hash       string    `json:"hash"`
	ParentHash string    `json:"parentHash"`
	Timestamp  time.Time `json:"timestamp"`
	RunID      string    `json:"runId"`
	Status     string    `json:"status"`     // "pass", "fail", "warn"
	FilePath   string    `json:"filePath"`   // Path to attestation file
}

// ChainManager handles evidence chain operations
type ChainManager struct {
	evidenceDir string
	chainPath   string
}

// NewChainManager creates a new chain manager
func NewChainManager(evidenceDir string) *ChainManager {
	return &ChainManager{
		evidenceDir: evidenceDir,
		chainPath:   filepath.Join(evidenceDir, "chain.json"),
	}
}

// LoadOrCreateChain loads existing chain or creates a new one
func (cm *ChainManager) LoadOrCreateChain() (*EvidenceChain, error) {
	if _, err := os.Stat(cm.chainPath); os.IsNotExist(err) {
		// Create new chain
		chain := &EvidenceChain{
			ChainID:      generateChainID(),
			StartTime:    time.Now().UTC(),
			LastUpdated:  time.Now().UTC(),
			Length:       0,
			Attestations: make([]ChainEntry, 0),
		}
		return chain, cm.SaveChain(chain)
	}
	
	// Load existing chain
	data, err := os.ReadFile(cm.chainPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chain file: %w", err)
	}
	
	var chain EvidenceChain
	if err := json.Unmarshal(data, &chain); err != nil {
		return nil, fmt.Errorf("failed to parse chain file: %w", err)
	}
	
	return &chain, nil
}

// SaveChain saves the evidence chain to disk
func (cm *ChainManager) SaveChain(chain *EvidenceChain) error {
	if err := os.MkdirAll(cm.evidenceDir, 0755); err != nil {
		return fmt.Errorf("failed to create evidence directory: %w", err)
	}
	
	data, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize chain: %w", err)
	}
	
	if err := os.WriteFile(cm.chainPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write chain file: %w", err)
	}
	
	return nil
}

// AddAttestation adds a new attestation to the chain
func (cm *ChainManager) AddAttestation(chain *EvidenceChain, attestation *Attestation, filePath string) error {
	// Set parent hash to current head (empty for genesis)
	parentHash := chain.Head
	if chain.Length == 0 {
		parentHash = ""
		chain.Genesis = attestation.Hash
	}
	
	// Update attestation with parent hash
	attestation.ParentHash = parentHash
	
	// Recalculate hash with parent hash included
	attestation.Hash = attestation.calculateHash()
	
	// Create chain entry
	entry := ChainEntry{
		Hash:       attestation.Hash,
		ParentHash: parentHash,
		Timestamp:  attestation.Timestamp,
		RunID:      attestation.RunID,
		Status:     attestation.Predicate.Summary.OverallStatus,
		FilePath:   filePath,
	}
	
	// Add to chain
	chain.Attestations = append(chain.Attestations, entry)
	chain.Length++
	chain.Head = attestation.Hash
	chain.LastUpdated = time.Now().UTC()
	
	return cm.SaveChain(chain)
}

// VerifyChain verifies the integrity of the evidence chain
func (cm *ChainManager) VerifyChain(chain *EvidenceChain) error {
	if len(chain.Attestations) == 0 {
		return nil // Empty chain is valid
	}
	
	// Verify genesis attestation has no parent
	if chain.Attestations[0].ParentHash != "" {
		return fmt.Errorf("genesis attestation must have empty parent hash")
	}
	
	// Verify hash chain links
	for i, entry := range chain.Attestations {
		// Load and verify attestation file
		if _, err := os.Stat(filepath.Join(cm.evidenceDir, entry.FilePath)); os.IsNotExist(err) {
			return fmt.Errorf("attestation file missing: %s", entry.FilePath)
		}
		
		// Verify parent hash linkage
		if i > 0 {
			expectedParent := chain.Attestations[i-1].Hash
			if entry.ParentHash != expectedParent {
				return fmt.Errorf("broken chain at position %d: parent hash mismatch", i)
			}
		}
	}
	
	// Verify head hash
	if len(chain.Attestations) > 0 {
		lastEntry := chain.Attestations[len(chain.Attestations)-1]
		if chain.Head != lastEntry.Hash {
			return fmt.Errorf("head hash mismatch: expected %s, got %s", lastEntry.Hash, chain.Head)
		}
	}
	
	return nil
}

// ScanAndRepairChain scans the evidence directory and rebuilds the chain
func (cm *ChainManager) ScanAndRepairChain() (*EvidenceChain, error) {
	attestationFiles, err := cm.findAttestationFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to scan attestation files: %w", err)
	}
	
	if len(attestationFiles) == 0 {
		// No attestations found, create empty chain
		chain := &EvidenceChain{
			ChainID:      generateChainID(),
			StartTime:    time.Now().UTC(),
			LastUpdated:  time.Now().UTC(),
			Length:       0,
			Attestations: make([]ChainEntry, 0),
		}
		return chain, cm.SaveChain(chain)
	}
	
	// Load and sort attestations by timestamp
	var entries []ChainEntry
	for _, file := range attestationFiles {
		entry, err := cm.loadAttestationEntry(file)
		if err != nil {
			fmt.Printf("Warning: failed to load attestation %s: %v\n", file, err)
			continue
		}
		entries = append(entries, entry)
	}
	
	// Sort by timestamp
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	
	// Rebuild chain with proper parent hash links
	var rebuiltEntries []ChainEntry
	var previousHash string
	
	for _, entry := range entries {
		// Update parent hash
		entry.ParentHash = previousHash
		rebuiltEntries = append(rebuiltEntries, entry)
		previousHash = entry.Hash
	}
	
	// Create rebuilt chain
	chain := &EvidenceChain{
		ChainID:      generateChainID(),
		StartTime:    rebuiltEntries[0].Timestamp,
		LastUpdated:  time.Now().UTC(),
		Length:       len(rebuiltEntries),
		Genesis:      rebuiltEntries[0].Hash,
		Head:         rebuiltEntries[len(rebuiltEntries)-1].Hash,
		Attestations: rebuiltEntries,
	}
	
	return chain, cm.SaveChain(chain)
}

// findAttestationFiles finds all attestation JSON files in the evidence directory
func (cm *ChainManager) findAttestationFiles() ([]string, error) {
	var files []string
	
	err := filepath.WalkDir(cm.evidenceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() && strings.HasPrefix(d.Name(), "attestation-") && strings.HasSuffix(d.Name(), ".json") {
			relPath, err := filepath.Rel(cm.evidenceDir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		
		return nil
	})
	
	return files, err
}

// loadAttestationEntry loads basic information from an attestation file
func (cm *ChainManager) loadAttestationEntry(filePath string) (ChainEntry, error) {
	fullPath := filepath.Join(cm.evidenceDir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return ChainEntry{}, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Try to parse as SignedAttestation first
	var signed SignedAttestation
	if err := json.Unmarshal(data, &signed); err == nil {
		// Extract attestation from DSSE envelope
		var attestation Attestation
		if err := json.Unmarshal([]byte(signed.Envelope.Payload), &attestation); err == nil {
			return ChainEntry{
				Hash:      attestation.Hash,
				Timestamp: attestation.Timestamp,
				RunID:     attestation.RunID,
				Status:    attestation.Predicate.Summary.OverallStatus,
				FilePath:  filePath,
			}, nil
		}
	}
	
	// Try to parse as plain Attestation
	var attestation Attestation
	if err := json.Unmarshal(data, &attestation); err == nil {
		return ChainEntry{
			Hash:      attestation.Hash,
			Timestamp: attestation.Timestamp,
			RunID:     attestation.RunID,
			Status:    attestation.Predicate.Summary.OverallStatus,
			FilePath:  filePath,
		}, nil
	}
	
	return ChainEntry{}, fmt.Errorf("failed to parse attestation file")
}

// generateChainID creates a unique chain identifier
func generateChainID() string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("mondrian-chain-%d", timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

// GetChainSummary returns a human-readable summary of the chain
func (chain *EvidenceChain) GetChainSummary() string {
	if chain.Length == 0 {
		return "Empty evidence chain"
	}
	
	passCount := 0
	failCount := 0
	warnCount := 0
	
	for _, entry := range chain.Attestations {
		switch entry.Status {
		case "pass":
			passCount++
		case "fail":
			failCount++
		case "warn":
			warnCount++
		}
	}
	
	return fmt.Sprintf("Chain: %d attestations (%d passed, %d failed, %d warnings)\nSpan: %s to %s",
		chain.Length, passCount, failCount, warnCount,
		chain.StartTime.Format("2006-01-02 15:04:05"),
		chain.LastUpdated.Format("2006-01-02 15:04:05"))
}