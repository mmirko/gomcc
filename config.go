package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// AppType represents the type of application
type AppType string

const (
	TypeCheck      AppType = "check"
	TypeExecutable AppType = "executable"
)

// DependencyAction defines what to execute based on dependency result
type DependencyAction struct {
	OnSuccess string `json:"on_success,omitempty"`
	OnFailure string `json:"on_failure,omitempty"`
}

// App represents an application configuration
type App struct {
	Name         string                      `json:"name"`
	Type         AppType                     `json:"type"`
	Command      string                      `json:"command"`
	Args         []string                    `json:"args,omitempty"`
	Tags         []string                    `json:"tags,omitempty"`
	Dependencies map[string]DependencyAction `json:"dependencies,omitempty"`
}

// Config represents the entire configuration file
type Config struct {
	Apps []App `json:"apps"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig performs validation on the loaded configuration
func validateConfig(config *Config) error {
	names := make(map[string]bool)

	for _, app := range config.Apps {
		// Check for unique names
		if names[app.Name] {
			return fmt.Errorf("duplicate app name: %s", app.Name)
		}
		names[app.Name] = true

		// Validate app type
		if app.Type != TypeCheck && app.Type != TypeExecutable {
			return fmt.Errorf("invalid app type '%s' for app '%s'", app.Type, app.Name)
		}

		// Validate command is not empty
		if app.Command == "" {
			return fmt.Errorf("app '%s' has empty command", app.Name)
		}
	}

	// Validate dependencies exist
	for _, app := range config.Apps {
		for depName := range app.Dependencies {
			if !names[depName] {
				return fmt.Errorf("app '%s' has dependency on non-existent app '%s'", app.Name, depName)
			}
		}
	}

	return nil
}

// GetApp returns an app by name
func (c *Config) GetApp(name string) *App {
	for i := range c.Apps {
		if c.Apps[i].Name == name {
			return &c.Apps[i]
		}
	}
	return nil
}

// GetAppsByTag returns all apps with a given tag
func (c *Config) GetAppsByTag(tag string) []App {
	var result []App
	for _, app := range c.Apps {
		for _, appTag := range app.Tags {
			if appTag == tag {
				result = append(result, app)
				break
			}
		}
	}
	return result
}

// GetAppsByTags returns all apps that have at least one of the given tags
func (c *Config) GetAppsByTags(tags []string) []App {
	if len(tags) == 0 {
		return c.Apps
	}

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag] = true
	}

	var result []App
	for _, app := range c.Apps {
		for _, appTag := range app.Tags {
			if tagMap[appTag] {
				result = append(result, app)
				break
			}
		}
	}
	return result
}
