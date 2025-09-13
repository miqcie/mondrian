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

package policy

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type CheckResult struct {
	RuleName    string                 `json:"rule_name"`
	Status      string                 `json:"status"` // "pass", "fail", "warn"
	Message     string                 `json:"message"`
	File        string                 `json:"file,omitempty"`
	Line        int                    `json:"line,omitempty"`
	Remediation string                 `json:"remediation,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type PolicyEngine struct {
	Rules []PolicyRule
}

type PolicyRule interface {
	Name() string
	Description() string
	Check(files map[string]string) []CheckResult
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		Rules: []PolicyRule{
			&S3PublicBucketRule{},
			&SecurityGroupOpenRule{},
			&MissingOIDCRule{},
		},
	}
}

func (pe *PolicyEngine) RunChecks(files map[string]string) []CheckResult {
	var results []CheckResult
	
	for _, rule := range pe.Rules {
		ruleResults := rule.Check(files)
		results = append(results, ruleResults...)
	}
	
	return results
}

// S3PublicBucketRule checks for public S3 buckets in Terraform
type S3PublicBucketRule struct{}

func (r *S3PublicBucketRule) Name() string {
	return "s3-no-public-buckets"
}

func (r *S3PublicBucketRule) Description() string {
	return "S3 buckets should not allow public read access"
}

func (r *S3PublicBucketRule) Check(files map[string]string) []CheckResult {
	var results []CheckResult
	
	for filename, content := range files {
		if !isTerraformFile(filename) {
			continue
		}
		
		// Check for public bucket configurations
		publicPatterns := []string{
			`public_read_write`,
			`public-read-write`,
			`public_read`,
			`public-read`,
			`"*".*s3:GetObject`,
		}
		
		lines := strings.Split(content, "\n")
		for lineNum, line := range lines {
			for _, pattern := range publicPatterns {
				matched, _ := regexp.MatchString(pattern, line)
				if matched {
					results = append(results, CheckResult{
						RuleName:    r.Name(),
						Status:      "fail",
						Message:     "S3 bucket configured with public access",
						File:        filename,
						Line:        lineNum + 1,
						Remediation: "Remove public access configuration or use specific IAM policies instead",
						Metadata: map[string]interface{}{
							"pattern": pattern,
							"line_content": strings.TrimSpace(line),
						},
					})
				}
			}
		}
	}
	
	if len(results) == 0 {
		results = append(results, CheckResult{
			RuleName: r.Name(),
			Status:   "pass",
			Message:  "No public S3 buckets detected",
		})
	}
	
	return results
}

// SecurityGroupOpenRule checks for overly permissive security groups
type SecurityGroupOpenRule struct{}

func (r *SecurityGroupOpenRule) Name() string {
	return "sg-no-open-ingress"
}

func (r *SecurityGroupOpenRule) Description() string {
	return "Security groups should not allow ingress from 0.0.0.0/0 on sensitive ports"
}

func (r *SecurityGroupOpenRule) Check(files map[string]string) []CheckResult {
	var results []CheckResult
	
	for filename, content := range files {
		if !isTerraformFile(filename) {
			continue
		}
		
		// Look for security group rules with 0.0.0.0/0
		openCIDRPattern := `0\.0\.0\.0/0`
		
		lines := strings.Split(content, "\n")
		inIngressBlock := false
		
		for lineNum, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			
			// Track if we're in an ingress block
			if strings.Contains(trimmedLine, "ingress {") {
				inIngressBlock = true
			} else if strings.Contains(trimmedLine, "egress {") {
				inIngressBlock = false
			} else if strings.Contains(trimmedLine, "}") && inIngressBlock {
				inIngressBlock = false
			}
			
			// Only flag 0.0.0.0/0 in ingress blocks
			matched, _ := regexp.MatchString(openCIDRPattern, line)
			if matched && inIngressBlock && strings.Contains(line, "cidr_blocks") {
				results = append(results, CheckResult{
					RuleName:    r.Name(),
					Status:      "fail",
					Message:     "Security group allows ingress from 0.0.0.0/0",
					File:        filename,
					Line:        lineNum + 1,
					Remediation: "Restrict ingress to specific CIDR blocks or security groups",
					Metadata: map[string]interface{}{
						"line_content": strings.TrimSpace(line),
					},
				})
			}
		}
	}
	
	if len(results) == 0 {
		results = append(results, CheckResult{
			RuleName: r.Name(),
			Status:   "pass",
			Message:  "No overly permissive security groups detected",
		})
	}
	
	return results
}

// MissingOIDCRule checks for proper OIDC configuration in CI/CD
type MissingOIDCRule struct{}

func (r *MissingOIDCRule) Name() string {
	return "deploy-require-oidc"
}

func (r *MissingOIDCRule) Description() string {
	return "Deployment workflows should use OIDC workload identity instead of long-lived credentials"
}

func (r *MissingOIDCRule) Check(files map[string]string) []CheckResult {
	var results []CheckResult
	foundOIDC := false
	foundSecrets := false
	
	for filename, content := range files {
		if !isGitHubActionFile(filename) {
			continue
		}
		
		// Check for OIDC configuration
		oidcPatterns := []string{
			`id-token:.*write`,
			`aws-actions/configure-aws-credentials`,
			`role-to-assume`,
		}
		
		for _, pattern := range oidcPatterns {
			matched, _ := regexp.MatchString(pattern, content)
			if matched {
				foundOIDC = true
				break
			}
		}
		
		// Check for problematic secrets usage
		secretPatterns := []string{
			`AWS_ACCESS_KEY_ID`,
			`AWS_SECRET_ACCESS_KEY`,
			`secrets\.AWS_ACCESS_KEY`,
		}
		
		lines := strings.Split(content, "\n")
		for lineNum, line := range lines {
			for _, pattern := range secretPatterns {
				matched, _ := regexp.MatchString(pattern, line)
				if matched {
					foundSecrets = true
					results = append(results, CheckResult{
						RuleName:    r.Name(),
						Status:      "fail",
						Message:     "GitHub Action uses long-lived credentials instead of OIDC",
						File:        filename,
						Line:        lineNum + 1,
						Remediation: "Replace with OIDC workload identity using role-to-assume",
						Metadata: map[string]interface{}{
							"pattern": pattern,
						},
					})
				}
			}
		}
	}
	
	if foundOIDC && !foundSecrets {
		results = append(results, CheckResult{
			RuleName: r.Name(),
			Status:   "pass",
			Message:  "OIDC workload identity properly configured",
		})
	} else if !foundOIDC && !foundSecrets {
		results = append(results, CheckResult{
			RuleName: r.Name(),
			Status:   "warn",
			Message:  "No deployment workflows detected",
		})
	}
	
	return results
}

// Helper functions
func isTerraformFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".tf" || ext == ".tfvars"
}

func isGitHubActionFile(filename string) bool {
	return strings.Contains(filename, ".github/workflows/") && filepath.Ext(filename) == ".yml" || filepath.Ext(filename) == ".yaml"
}

// FormatResults formats check results for display
func FormatResults(results []CheckResult) string {
	var output strings.Builder
	
	passCount := 0
	failCount := 0
	warnCount := 0
	
	for _, result := range results {
		switch result.Status {
		case "pass":
			passCount++
			fmt.Fprintf(&output, "✅ %s: %s\n", result.RuleName, result.Message)
		case "fail":
			failCount++
			fmt.Fprintf(&output, "❌ %s: %s\n", result.RuleName, result.Message)
			if result.File != "" {
				fmt.Fprintf(&output, "   📁 %s:%d\n", result.File, result.Line)
			}
			if result.Remediation != "" {
				fmt.Fprintf(&output, "   💡 %s\n", result.Remediation)
			}
		case "warn":
			warnCount++
			fmt.Fprintf(&output, "⚠️  %s: %s\n", result.RuleName, result.Message)
		}
	}
	
	fmt.Fprintf(&output, "\n📊 Summary: %d passed, %d failed, %d warnings\n", passCount, failCount, warnCount)
	
	if failCount > 0 {
		fmt.Fprintf(&output, "\n🚫 Policy check failed - %d violations found\n", failCount)
	} else {
		fmt.Fprintf(&output, "\n✅ All policy checks passed!\n")
	}
	
	return output.String()
}

// ToJSON converts results to JSON for attestation
func ToJSON(results []CheckResult) ([]byte, error) {
	return json.MarshalIndent(results, "", "  ")
}