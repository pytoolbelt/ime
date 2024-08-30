package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

// Define the structs to match the updated YAML structure
type Config struct {
	GlobalPrefix string             `mapstructure:"global_prefix"`
	Projects     map[string]Project `mapstructure:"projects"`
}

type Project struct {
	Prefix       string                 `mapstructure:"prefix"`
	Environments map[string]Environment `mapstructure:"environments"`
}

type Environment struct {
	Prefix    string `mapstructure:"prefix"`
	LocalPath string `mapstructure:"local_path"`
}

func (e *Environment) GetResolvedLocalPath() string {
	return os.ExpandEnv(e.LocalPath)
}

// Method to get the environment path
func (c *Config) GetEnvironment(projectName, environmentName string) (*Environment, error) {
	project, exists := c.Projects[projectName]
	if !exists {
		return nil, fmt.Errorf("project %s not found", projectName)
	}

	env, exists := project.Environments[environmentName]
	if !exists {
		return nil, fmt.Errorf("environment %s not found in project %s", environmentName, projectName)
	}

	return &env, nil
}

func (c *Config) GetProject(projectName string) (*Project, error) {
	project, exists := c.Projects[projectName]
	if !exists {
		return nil, fmt.Errorf("project %s not found", projectName)
	}

	return &project, nil
}

// Function to validate the config
func (c *Config) ValidateConfig() error {
	if !strings.HasPrefix(c.GlobalPrefix, "/") {
		return fmt.Errorf("global_prefix must start with '/' got %s", c.GlobalPrefix)
	}

	for projectName, project := range c.Projects {
		if !strings.HasPrefix(project.Prefix, "/") {
			return fmt.Errorf("prefix for project %s must start with '/'", projectName)
		}

		for envName, env := range project.Environments {
			if !strings.HasPrefix(env.Prefix, "/") {
				return fmt.Errorf("prefix for environment %s in project %s must start with '/'", envName, projectName)
			}
		}
	}

	return nil
}

func (c *Config) FormatParameterStorePath(projectName, environmentName string) (string, error) {

	env, err := c.GetEnvironment(projectName, environmentName)
	if err != nil {
		return "", err
	}

	prj, err := c.GetProject(projectName)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s%s%s", c.GlobalPrefix, prj.Prefix, env.Prefix)
	return path, nil
}

// Method to print the config as a table
func (c *Config) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Project", "Environment", "Prefix", "Local Path"})

	for projectName, project := range c.Projects {
		for envName, env := range project.Environments {
			table.Append([]string{projectName, envName, env.Prefix, env.LocalPath})
		}
	}

	table.Render() // Send output
}

// Standalone function to load the configuration from a file
func LoadConfig() (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := config.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return &config, nil
}

func InitializeConfig() {
	viper.SetConfigName("ime")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.ime")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
