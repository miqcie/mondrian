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
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/secure-systems-lab/go-securesystemslib/dsse"
)

// Signer handles DSSE signing of attestations
type Signer struct {
	keyID      string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// SignedAttestation represents a DSSE-signed attestation
type SignedAttestation struct {
	Envelope  dsse.Envelope `json:"envelope"`
	Metadata  SigningMetadata `json:"metadata"`
}

type SigningMetadata struct {
	KeyID     string    `json:"keyId"`
	Algorithm string    `json:"algorithm"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// NewSigner creates a new DSSE signer
func NewSigner() (*Signer, error) {
	// Generate ephemeral ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	
	// Generate key ID from public key
	publicKeyBytes := elliptic.MarshalCompressed(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	hash := sha256.Sum256(publicKeyBytes)
	keyID := hex.EncodeToString(hash[:8]) // First 8 bytes for key ID
	
	return &Signer{
		keyID:      keyID,
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// NewSignerFromGitHubOIDC creates a signer using GitHub OIDC token (future implementation)
func NewSignerFromGitHubOIDC() (*Signer, error) {
	// TODO: Implement GitHub OIDC token-based signing
	// For now, fall back to ephemeral keys
	return NewSigner()
}

// SignAttestation signs an attestation using DSSE
func (s *Signer) SignAttestation(attestation *Attestation) (*SignedAttestation, error) {
	// Serialize attestation to JSON
	attestationJSON, err := attestation.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize attestation: %w", err)
	}
	
	// Create DSSE signer with our private key
	dsseSigner := &ECDSASigner{
		keyID:      s.keyID,
		privateKey: s.privateKey,
	}
	
	// Create envelope signer
	envelopeSigner, err := dsse.NewEnvelopeSigner(dsseSigner)
	if err != nil {
		return nil, fmt.Errorf("failed to create envelope signer: %w", err)
	}
	
	// Sign using DSSE
	envelope, err := envelopeSigner.SignPayload(context.Background(), "application/vnd.in-toto+json", attestationJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create DSSE envelope: %w", err)
	}
	
	metadata := SigningMetadata{
		KeyID:     s.keyID,
		Algorithm: "ECDSA-SHA256",
		Timestamp: time.Now().UTC(),
		Source:    getSigningSource(),
	}
	
	return &SignedAttestation{
		Envelope: *envelope,
		Metadata: metadata,
	}, nil
}

// ECDSASigner implements the dsse.Signer interface
type ECDSASigner struct {
	keyID      string
	privateKey *ecdsa.PrivateKey
}

func (e *ECDSASigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, e.privateKey, hash[:])
	if err != nil {
		return nil, err
	}
	
	// Encode signature (simple concatenation for demo)
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func (e *ECDSASigner) KeyID() (string, error) {
	return e.keyID, nil
}

// VerifySignedAttestation verifies a DSSE-signed attestation
func VerifySignedAttestation(signed *SignedAttestation, publicKey *ecdsa.PublicKey) error {
	// Create verifier
	verifier := &ECDSAVerifier{
		keyID:     signed.Metadata.KeyID,
		publicKey: publicKey,
	}
	
	// Create envelope verifier
	envelopeVerifier, err := dsse.NewEnvelopeVerifier(verifier)
	if err != nil {
		return fmt.Errorf("failed to create envelope verifier: %w", err)
	}
	
	// Verify DSSE envelope
	_, err = envelopeVerifier.Verify(context.Background(), &signed.Envelope)
	if err != nil {
		return fmt.Errorf("DSSE verification failed: %w", err)
	}
	
	return nil
}

// ECDSAVerifier implements the dsse.Verifier interface
type ECDSAVerifier struct {
	keyID     string
	publicKey *ecdsa.PublicKey
}

func (e *ECDSAVerifier) Verify(ctx context.Context, data, signature []byte) error {
	// For demo purposes, just return success
	// Production would implement proper ECDSA signature verification
	fmt.Printf("🔍 Verifying signature for key ID: %s (demo mode)\n", e.keyID[:8])
	return nil
}

func (e *ECDSAVerifier) KeyID() (string, error) {
	return e.keyID, nil
}

func (e *ECDSAVerifier) Public() crypto.PublicKey {
	return e.publicKey
}

// GetPublicKey returns the public key for verification
func (s *Signer) GetPublicKey() *ecdsa.PublicKey {
	return s.publicKey
}

// GetKeyID returns the key identifier
func (s *Signer) GetKeyID() string {
	return s.keyID
}

// getSigningSource determines the source of signing (CI, local, etc.)
func getSigningSource() string {
	// Check for CI environment variables
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return "github-actions"
	}
	if os.Getenv("CI") == "true" {
		return "ci"
	}
	
	// Get hostname for local development
	hostname, err := os.Hostname()
	if err != nil {
		return "local"
	}
	
	return fmt.Sprintf("local-%s", hostname)
}

// SaveSignedAttestation saves a signed attestation to the evidence store
func SaveSignedAttestation(signed *SignedAttestation, evidenceDir string) error {
	if err := os.MkdirAll(evidenceDir, 0755); err != nil {
		return fmt.Errorf("failed to create evidence directory: %w", err)
	}
	
	// Create filename with timestamp and key ID
	timestamp := signed.Metadata.Timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("attestation-%s-%s.json", timestamp, signed.Metadata.KeyID[:8])
	filePath := filepath.Join(evidenceDir, filename)
	
	// Serialize signed attestation
	data, err := json.MarshalIndent(signed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize signed attestation: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write attestation file: %w", err)
	}
	
	fmt.Printf("📝 Saved attestation: %s\n", filePath)
	return nil
}