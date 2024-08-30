package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGetEnvironment(t *testing.T) {
	config := &Config{
		GlobalPrefix: "/global",
		Projects: map[string]Project{
			"project1": {
				Prefix: "/project1",
				Environments: map[string]Environment{
					"dev": {
						Prefix:    "/dev",
						LocalPath: "/local/dev",
					},
					"prod": {
						Prefix:    "/prod",
						LocalPath: "/local/prod",
					},
				},
			},
		},
	}

	tests := []struct {
		projectName     string
		environmentName string
		expectedPrefix  string
		expectedPath    string
		expectError     bool
	}{
		{"project1", "dev", "/dev", "/local/dev", false},
		{"project1", "prod", "/prod", "/local/prod", false},
		{"project1", "staging", "", "", true},
		{"project2", "dev", "", "", true},
	}

	for _, tt := range tests {
		env, err := config.GetEnvironment(tt.projectName, tt.environmentName)
		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for project %s and environment %s, but got none", tt.projectName, tt.environmentName)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for project %s and environment %s: %v", tt.projectName, tt.environmentName, err)
			} else {
				if env.Prefix != tt.expectedPrefix {
					t.Errorf("expected prefix %s, but got %s", tt.expectedPrefix, env.Prefix)
				}
				if env.LocalPath != tt.expectedPath {
					t.Errorf("expected local path %s, but got %s", tt.expectedPath, env.LocalPath)
				}
			}
		}
	}
}

func TestGetProject(t *testing.T) {
	config := &Config{
		GlobalPrefix: "/global",
		Projects: map[string]Project{
			"project1": {
				Prefix: "/project1",
				Environments: map[string]Environment{
					"dev": {
						Prefix:    "/dev",
						LocalPath: "/local/dev",
					},
				},
			},
		},
	}

	tests := []struct {
		projectName    string
		expectedPrefix string
		expectError    bool
	}{
		{"project1", "/project1", false},
		{"project2", "", true},
	}

	for _, tt := range tests {
		project, err := config.GetProject(tt.projectName)
		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for project %s, but got none", tt.projectName)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for project %s: %v", tt.projectName, err)
			} else {
				if project.Prefix != tt.expectedPrefix {
					t.Errorf("expected prefix %s, but got %s", tt.expectedPrefix, project.Prefix)
				}
			}
		}
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		config      Config
		expectError bool
	}{
		{
			Config{
				GlobalPrefix: "/global",
				Projects: map[string]Project{
					"project1": {
						Prefix: "/project1",
						Environments: map[string]Environment{
							"dev": {
								Prefix:    "/dev",
								LocalPath: "/local/dev",
							},
						},
					},
				},
			},
			false,
		},
		{
			Config{
				GlobalPrefix: "global",
				Projects: map[string]Project{
					"project1": {
						Prefix: "/project1",
						Environments: map[string]Environment{
							"dev": {
								Prefix:    "/dev",
								LocalPath: "/local/dev",
							},
						},
					},
				},
			},
			true,
		},
		{
			Config{
				GlobalPrefix: "/global",
				Projects: map[string]Project{
					"project1": {
						Prefix: "project1",
						Environments: map[string]Environment{
							"dev": {
								Prefix:    "/dev",
								LocalPath: "/local/dev",
							},
						},
					},
				},
			},
			true,
		},
		{
			Config{
				GlobalPrefix: "/global",
				Projects: map[string]Project{
					"project1": {
						Prefix: "/project1",
						Environments: map[string]Environment{
							"dev": {
								Prefix:    "dev",
								LocalPath: "/local/dev",
							},
						},
					},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		err := tt.config.ValidateConfig()
		if tt.expectError {
			if err == nil {
				t.Errorf("expected error, but got none")
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestFormatParameterStorePath(t *testing.T) {
	config := &Config{
		GlobalPrefix: "/global",
		Projects: map[string]Project{
			"project1": {
				Prefix: "/project1",
				Environments: map[string]Environment{
					"dev": {
						Prefix:    "/dev",
						LocalPath: "/local/dev",
					},
				},
			},
		},
	}

	tests := []struct {
		projectName     string
		environmentName string
		expectedPath    string
		expectError     bool
	}{
		{"project1", "dev", "/global/project1/dev", false},
		{"project1", "prod", "", true},
		{"project2", "dev", "", true},
	}

	for _, tt := range tests {
		path, err := config.FormatParameterStorePath(tt.projectName, tt.environmentName)
		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for project %s and environment %s, but got none", tt.projectName, tt.environmentName)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for project %s and environment %s: %v", tt.projectName, tt.environmentName, err)
			} else {
				if path != tt.expectedPath {
					t.Errorf("expected path %s, but got %s", tt.expectedPath, path)
				}
			}
		}
	}
}

func TestGetResolvedLocalPath(t *testing.T) {
	os.Setenv("TEST_PATH", "/test/path")
	env := &Environment{
		LocalPath: "$TEST_PATH",
	}

	expectedPath := "/test/path"
	resolvedPath := env.GetResolvedLocalPath()
	if resolvedPath != expectedPath {
		t.Errorf("expected resolved path %s, but got %s", expectedPath, resolvedPath)
	}
}

func TestLoadConfig(t *testing.T) {
	// Set up a temporary config file for testing
	configContent := `
global_prefix: /global
projects:
  project1:
    prefix: /project1
    environments:
      dev:
        prefix: /dev
        local_path: /local/dev
`
	configFile := "test_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}
	defer os.Remove(configFile)

	viper.SetConfigFile(configFile)
	viper.ReadInConfig()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	expectedGlobalPrefix := "/global"
	if config.GlobalPrefix != expectedGlobalPrefix {
		t.Errorf("expected global prefix %s, but got %s", expectedGlobalPrefix, config.GlobalPrefix)
	}

	expectedProjectPrefix := "/project1"
	project, err := config.GetProject("project1")
	if err != nil {
		t.Fatalf("failed to get project: %v", err)
	}
	if project.Prefix != expectedProjectPrefix {
		t.Errorf("expected project prefix %s, but got %s", expectedProjectPrefix, project.Prefix)
	}

	expectedEnvPrefix := "/dev"
	env, err := config.GetEnvironment("project1", "dev")
	if err != nil {
		t.Fatalf("failed to get environment: %v", err)
	}
	if env.Prefix != expectedEnvPrefix {
		t.Errorf("expected environment prefix %s, but got %s", expectedEnvPrefix, env.Prefix)
	}
}
