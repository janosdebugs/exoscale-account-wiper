package terraform

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"os"
	"testing"
)

type TerraformTestScaffold struct {
	terraformOptions *terraform.Options
	t                *testing.T
	ExoscaleKey      string
	ExoscaleSecret   string
}

func New(t *testing.T, dir string) *TerraformTestScaffold {
	var exoscaleKey = os.Getenv("EXOSCALE_KEY")
	var exoscaleSecret = os.Getenv("EXOSCALE_SECRET")
	if exoscaleKey == "" {
		t.Skip("EXOSCALE_KEY not set")
		return nil
	}
	if exoscaleSecret == "" {
		t.Skip("EXOSCALE_SECRET not set")
		return nil
	}

	terraformOptions := &terraform.Options{
		TerraformDir: dir,
		Vars: map[string]interface{}{
			"exoscale_key":    exoscaleKey,
			"exoscale_secret": exoscaleSecret,
		},
	}

	return &TerraformTestScaffold{
		terraformOptions: terraformOptions,
		t:                t,
		ExoscaleKey:      exoscaleKey,
		ExoscaleSecret:   exoscaleSecret,
	}
}

func (s *TerraformTestScaffold) Apply() {
	terraform.InitAndApply(s.t, s.terraformOptions)
}

func (s *TerraformTestScaffold) Destroy() {
	terraform.Destroy(s.t, s.terraformOptions)
}
