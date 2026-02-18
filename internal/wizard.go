package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RunConfigWizard() {
	scanner := bufio.NewScanner(os.Stdin)

	name := promptRequired(scanner, "Config name or path (saves to configs/<name>.json if no path given): ")

	var savePath string
	if strings.Contains(name, "/") || strings.HasSuffix(name, ".json") {
		if !filepath.IsAbs(name) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Could not determine working directory: %v", err)
			}
			savePath = filepath.Join(cwd, name)
		} else {
			savePath = name
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not determine working directory: %v", err)
		}
		savePath = filepath.Join(cwd, "configs", name+".json")
	}

	branch := promptRequired(scanner, "Branch: ")

	fmt.Println("Add repo (owner/repo), empty line to finish:")
	var repos []string
	for {
		fmt.Print("  > ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(repos) == 0 {
				fmt.Println("  (at least one repo required)")
				continue
			}
			break
		}
		repos = append(repos, line)
	}

	probe := promptBool(scanner, "Probe mode? (y/N): ")
	merge := promptBool(scanner, "Auto-merge? (y/N): ")

	cfg := Config{
		Branch: branch,
		Repos:  repos,
		Probe:  probe,
		Merge:  merge,
	}

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(savePath, data, 0644); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	fmt.Println("\nSaved to " + savePath)
}

func promptRequired(scanner *bufio.Scanner, message string) string {
	for {
		fmt.Print(message)
		if !scanner.Scan() {
			log.Fatal("unexpected end of input")
		}
		val := strings.TrimSpace(scanner.Text())
		if val != "" {
			return val
		}
		fmt.Println("  (required, please enter a value)")
	}
}

func promptBool(scanner *bufio.Scanner, message string) bool {
	fmt.Print(message)
	if !scanner.Scan() {
		return false
	}
	val := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return val == "y" || val == "yes"
}
