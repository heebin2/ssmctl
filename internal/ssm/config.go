package ssm

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Global    GlobalConfig              `yaml:"global"`
	Instances map[string]InstanceConfig `yaml:"instances"`
}

type GlobalConfig struct {
	User string `yaml:"user"`
}

type InstanceConfig struct {
	Target string `yaml:"target"`
}

type instance struct {
	name   string
	target string
	user   string
}

var linuxUserPattern = regexp.MustCompile(`^[a-z_][a-z0-9_-]*[$]?$`)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("read: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	if len(cfg.Instances) == 0 {
		return nil, fmt.Errorf("no instances defined in config")
	}

	for name, inst := range cfg.Instances {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("instance name cannot be empty")
		}
		if strings.TrimSpace(inst.Target) == "" {
			return nil, fmt.Errorf("instance %q: target is required", name)
		}
	}

	return &cfg, nil
}

func (c *Config) resolve(name string) (instance, error) {
	inst, ok := c.Instances[name]
	if !ok {
		available := make([]string, 0, len(c.Instances))
		for n := range c.Instances {
			available = append(available, n)
		}
		sort.Strings(available)
		return instance{}, fmt.Errorf("unknown instance %q (available: %v)", name, available)
	}

	user := strings.TrimSpace(c.Global.User)
	if user == "" {
		return instance{}, fmt.Errorf("instance %q: user not configured (set global.user)", name)
	}
	if !linuxUserPattern.MatchString(user) {
		return instance{}, fmt.Errorf("instance %q: invalid user %q (must start with letter/underscore)", name, user)
	}

	return instance{name: name, target: inst.Target, user: user}, nil
}

func (c *Config) listInstances() []instance {
	insts := make([]instance, 0, len(c.Instances))

	globalUser := strings.TrimSpace(c.Global.User)
	if globalUser == "" {
		globalUser = "-"
	}

	for name, inst := range c.Instances {
		insts = append(insts, instance{name: name, target: inst.Target, user: globalUser})
	}

	sort.Slice(insts, func(i, j int) bool {
		return insts[i].name < insts[j].name
	})

	return insts
}

func PrintList(cfg *Config) {
	fmt.Printf("%-40s %-20s %s\n", "NAME", "TARGET", "USER")
	for _, inst := range cfg.listInstances() {
		fmt.Printf("%-40s %-20s %s\n", inst.name, inst.target, inst.user)
	}
}

func Start(cfg *Config, name string) error {
	inst, err := cfg.resolve(name)
	if err != nil {
		return err
	}
	return startSession(inst.target, inst.user)
}
