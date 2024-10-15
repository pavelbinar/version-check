package check

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// ToolConfig represents the configuration for each tool with a name.
type ToolConfig struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	Expect  string `yaml:"expect"`
}

// Config holds an array of tools from the YAML configuration.
type Config struct {
	Tools []ToolConfig `yaml:"tools"`
}

// readConfig reads and parses the config.yaml file.
func readConfig(cfgFile string) (*Config, error) {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// extractVersion extracts the version from the command output.
func extractVersion(command string, output string) string {
	if strings.Contains(command, "rsync") {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "rsync version") {
				re := regexp.MustCompile(`v?(\d+(\.\d+){0,2})`)
				match := re.FindString(line)
				return strings.TrimPrefix(match, "v")
			}
		}
		return ""
	}

	re := regexp.MustCompile(`v?(\d+(\.\d+){0,2})`)
	match := re.FindString(output)
	return strings.TrimPrefix(match, "v")
}

// compareVersions compares two version strings.
func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] != parts2[i] {
			return false
		}
	}

	return len(parts1) == len(parts2)
}

// checkVersion checks if the output of the command contains the expected version.
func checkVersion(config ToolConfig) bool {
	cmd := exec.Command("sh", "-c", config.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command for %s: %s\n", config.Name, err)
		return false
	}
	actualOutput := strings.TrimSpace(string(output))
	extractedVersion := extractVersion(config.Command, actualOutput)
	if extractedVersion == "" {
		fmt.Printf("%s: No version found in output\n", config.Name)
		fmt.Printf("Command output: %s\n", actualOutput)
		return false
	}
	if compareVersions(extractedVersion, config.Expect) {
		return true
	} else {
		fmt.Printf("%s version mismatch: Expected '%s', got '%s'\n", config.Name, config.Expect, extractedVersion)
		fmt.Printf("Command output: %s\n", actualOutput)
		return false
	}
}

func RunVersionCheck(cmd *cobra.Command, args []string, cfgFile string) {
	config, err := readConfig(cfgFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %s\n", err)
		os.Exit(1)
	}

	allMatch := true
	for _, toolConfig := range config.Tools {
		if !checkVersion(toolConfig) {
			allMatch = false
		}
	}

	if allMatch {
		fmt.Println("Versions OK")
	} else {
		os.Exit(1)
	}
}
