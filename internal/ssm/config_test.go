package ssm

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigResolveInstancePreferGlobalUser(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		instName string
		wantUser string
		wantErr  bool
	}{
		{
			name: "use global user",
			cfg: &Config{
				Global: GlobalConfig{User: "ec2-user"},
				Instances: map[string]InstanceConfig{
					"prod": {Target: "i-123"},
				},
			},
			instName: "prod",
			wantUser: "ec2-user",
		},
		{
			name: "no user configured",
			cfg: &Config{
				Instances: map[string]InstanceConfig{
					"prod": {Target: "i-123"},
				},
			},
			instName: "prod",
			wantErr:  true,
		},
		{
			name: "unknown instance",
			cfg: &Config{
				Instances: map[string]InstanceConfig{
					"prod": {Target: "i-123"},
				},
			},
			instName: "dev",
			wantErr:  true,
		},
		{
			name: "invalid user",
			cfg: &Config{
				Global: GlobalConfig{User: "user-name!"},
				Instances: map[string]InstanceConfig{
					"prod": {Target: "i-123"},
				},
			},
			instName: "prod",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst, err := tt.cfg.resolve(tt.instName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolve: got error %v, want error %v", err, tt.wantErr)
			}
			if err == nil && inst.user != tt.wantUser {
				t.Fatalf("user: got %s, want %s", inst.user, tt.wantUser)
			}
		})
	}
}

func TestConfigListInstancesSorted(t *testing.T) {
	cfg := &Config{
		Global: GlobalConfig{User: "ec2-user"},
		Instances: map[string]InstanceConfig{
			"z-app": {Target: "i-z"},
			"a-app": {Target: "i-a"},
			"m-app": {Target: "i-m"},
		},
	}

	insts := cfg.listInstances()
	if len(insts) != 3 {
		t.Fatalf("got %d instances, want 3", len(insts))
	}

	want := []string{"a-app", "m-app", "z-app"}
	for i, exp := range want {
		if insts[i].name != exp {
			t.Fatalf("position %d: got %s, want %s", i, insts[i].name, exp)
		}
	}
}

func TestConfigInstanceNamesSorted(t *testing.T) {
	cfg := &Config{Instances: map[string]InstanceConfig{
		"z-app": {Target: "i-z"},
		"a-app": {Target: "i-a"},
		"m-app": {Target: "i-m"},
	}}

	got := cfg.InstanceNames()
	want := []string{"a-app", "m-app", "z-app"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InstanceNames() = %v, want %v", got, want)
	}
}

func TestLoadConfigParseYAML(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		wantErr   bool
		wantNames []string
	}{
		{
			name:      "valid config",
			yaml:      "global:\n  user: ec2-user\ninstances:\n  prod:\n    target: i-123\n",
			wantNames: []string{"prod"},
		},
		{
			name:    "no instances",
			yaml:    "global:\n  user: ec2-user\ninstances: {}\n",
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			yaml:    "global: {invalid",
			wantErr: true,
		},
		{
			name:    "missing target",
			yaml:    "instances:\n  prod:\n    user: ec2-user\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.yml")
			if err := writeFile(path, tt.yaml); err != nil {
				t.Fatalf("write config: %v", err)
			}

			cfg, err := LoadConfig(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadConfig: got error %v, want error %v", err, tt.wantErr)
			}

			if err == nil {
				if len(cfg.Instances) != len(tt.wantNames) {
					t.Fatalf("instances: got %d, want %d", len(cfg.Instances), len(tt.wantNames))
				}
			}
		})
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !contains(err.Error(), "file not found") {
		t.Fatalf("expected 'file not found' in error, got: %v", err)
	}
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
