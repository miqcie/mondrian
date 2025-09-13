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
	"os"
	"path/filepath"
	"strings"
)

type FileScanner struct {
	rootDir string
}

func NewFileScanner(rootDir string) *FileScanner {
	return &FileScanner{rootDir: rootDir}
}

func (fs *FileScanner) ScanRelevantFiles() (map[string]string, error) {
	files := make(map[string]string)
	
	err := filepath.Walk(fs.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip hidden directories and common non-relevant dirs
		if info.IsDir() {
			dirName := info.Name()
			if strings.HasPrefix(dirName, ".") && dirName != ".github" {
				return filepath.SkipDir
			}
			if dirName == "node_modules" || dirName == "vendor" || dirName == ".terraform" {
				return filepath.SkipDir
			}
		}
		
		if !info.IsDir() && fs.isRelevantFile(path) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			// Use relative path as key
			relPath, err := filepath.Rel(fs.rootDir, path)
			if err != nil {
				relPath = path
			}
			
			files[relPath] = string(content)
		}
		
		return nil
	})
	
	return files, err
}

func (fs *FileScanner) isRelevantFile(path string) bool {
	// Check file extensions we care about
	ext := filepath.Ext(path)
	
	relevantExts := map[string]bool{
		".tf":     true, // Terraform
		".tfvars": true, // Terraform variables
		".yml":    true, // YAML (GitHub Actions, etc.)
		".yaml":   true, // YAML
		".json":   true, // JSON configs
		".hcl":    true, // HashiCorp Configuration Language
	}
	
	if relevantExts[ext] {
		return true
	}
	
	// Check specific file patterns
	fileName := filepath.Base(path)
	relevantFiles := []string{
		"Dockerfile",
		"docker-compose.yml",
		"docker-compose.yaml",
	}
	
	for _, relevant := range relevantFiles {
		if fileName == relevant {
			return true
		}
	}
	
	// Check if it's in .github/workflows/
	if strings.Contains(path, ".github/workflows/") && (ext == ".yml" || ext == ".yaml") {
		return true
	}
	
	return false
}

func (fs *FileScanner) GetFileContent(relativePath string) (string, error) {
	fullPath := filepath.Join(fs.rootDir, relativePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}