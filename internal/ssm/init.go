package ssm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type EC2Instance struct {
	InstanceID string `json:"InstanceId"`
	State      struct {
		Name string `json:"Name"`
	} `json:"State"`
	Tags []struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	} `json:"Tags"`
}

type EC2Response struct {
	Reservations []struct {
		Instances []EC2Instance `json:"Instances"`
	} `json:"Reservations"`
}

func InitConfig(path string) error {
	instances, err := fetchEC2Instances()
	if err != nil {
		return fmt.Errorf("fetch ec2 instances: %w", err)
	}

	if len(instances) == 0 {
		return fmt.Errorf("no ec2 instances found")
	}

	cfg := &Config{
		Global:    GlobalConfig{User: "ec2-user"},
		Instances: make(map[string]InstanceConfig),
	}

	for _, inst := range instances {
		if inst.State.Name != "running" {
			continue
		}

		name := getNameTag(inst)
		if name == "" {
			name = inst.InstanceID
		}

		// Use only global user, no per-instance override
		cfg.Instances[name] = InstanceConfig{
			Target: inst.InstanceID,
		}
	}

	if len(cfg.Instances) == 0 {
		return fmt.Errorf("no running ec2 instances found")
	}

	// Marshal with compact formatting
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetDefaultFlowStyle(false)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func fetchEC2Instances() ([]EC2Instance, error) {
	cmd := exec.Command(
		"aws", "ec2", "describe-instances",
		"--output", "json",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("aws cli: %w", err)
	}

	var resp EC2Response
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	var instances []EC2Instance
	for _, res := range resp.Reservations {
		instances = append(instances, res.Instances...)
	}

	return instances, nil
}

func getNameTag(inst EC2Instance) string {
	for _, tag := range inst.Tags {
		if tag.Key == "Name" {
			return tag.Value
		}
	}
	return ""
}
