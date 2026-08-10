package ssm

import "testing"

func TestGetNameTag(t *testing.T) {
	tests := []struct {
		name     string
		inst     EC2Instance
		wantName string
	}{
		{
			name: "with name tag",
			inst: EC2Instance{
				InstanceID: "i-123",
				Tags: []struct {
					Key   string `json:"Key"`
					Value string `json:"Value"`
				}{
					{Key: "Name", Value: "my-server"},
					{Key: "Environment", Value: "prod"},
				},
			},
			wantName: "my-server",
		},
		{
			name: "no name tag",
			inst: EC2Instance{
				InstanceID: "i-123",
				Tags: []struct {
					Key   string `json:"Key"`
					Value string `json:"Value"`
				}{
					{Key: "Environment", Value: "prod"},
				},
			},
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNameTag(tt.inst)
			if got != tt.wantName {
				t.Fatalf("got %q, want %q", got, tt.wantName)
			}
		})
	}
}
