/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ControllerManagerConfiguration defines the configuration for the Gardener controller manager.
type ControllerManagerConfiguration struct {
	// +optional
	Kind string `yaml:"kind"`
	// +optional
	APIVersion string `yaml:"apiVersion"`

	// Controllers defines the configuration of the controllers.
	Controllers ControllerManagerControllerConfiguration `yaml:"controllers"`
	// Webhooks defines the configuration of the admission webhooks.
	Webhooks ControllerManagerWebhookConfiguration `yaml:"webhooks"`
}

// ControllerManagerControllerConfiguration defines the configuration of the controllers.
type ControllerManagerControllerConfiguration struct {
	// Shoot defines the configuration of the Shoot controller.
	Shoot ShootControllerConfiguration `yaml:"shoot"`
}

// ShootControllerConfiguration defines the configuration of the Shoot controller.
type ShootControllerConfiguration struct {
	// MaxConcurrentReconciles is the maximum number of concurrent Reconciles which can be run. Defaults to 50.
	MaxConcurrentReconciles int `yaml:"maxConcurrentReconciles"`

	// MaxConcurrentReconcilesPerNamespace is the maximum number of concurrent Reconciles which can be run per Namespace (independent of the user who created the Shoot resource). Defaults to 3.
	MaxConcurrentReconcilesPerNamespace int `yaml:"maxConcurrentReconcilesPerNamespace"`

	// QuotaExceededRetryDelay is the duration, after which the reconciliation will be retried again in case the configMap quota is exceeded.
	// Note that in case the resource quota for count/configmaps is increased or configMap quota was freed a reconciliation is requested for all shoots in the namespace that do not already have a corresponding <shootname>.kubeconfig configMap.
	// Defaults to 24 hours.
	QuotaExceededRetryDelay time.Duration `yaml:"quotaExceededRetryDelay"`
}

// ControllerManagerWebhookConfiguration defines the configuration of the admission webhooks.
type ControllerManagerWebhookConfiguration struct {
	// ConfigMapValidation defines the configuration of the validating webhook.
	ConfigMapValidation ConfigMapValidatingWebhookConfiguration `yaml:"configMapValidation"`
}

// ConfigMapValidatingWebhookConfiguration defines the configuration of the validating webhook.
type ConfigMapValidatingWebhookConfiguration struct {
	// MaxObjectSize is the maximum size of a configMap resource in bytes. Defaults to 102400.
	MaxObjectSize int `yaml:"maxObjectSize"`
}

// ReadControllerManagerConfiguration returns a valid ControllerManagerConfiguration struct.
// The ControllerManagerConfiguration is initialized by reading the config file from the given file path (if the value is not empty), with defaults applied.
func ReadControllerManagerConfiguration(configFile string) (*ControllerManagerConfiguration, error) {
	// Default configuration
	cfg := ControllerManagerConfiguration{
		Controllers: ControllerManagerControllerConfiguration{
			Shoot: ShootControllerConfiguration{
				MaxConcurrentReconciles:             50,
				MaxConcurrentReconcilesPerNamespace: 3,
				QuotaExceededRetryDelay:             24 * time.Hour,
			},
		},
		Webhooks: ControllerManagerWebhookConfiguration{
			ConfigMapValidation: ConfigMapValidatingWebhookConfiguration{
				MaxObjectSize: 100 * 1024,
			},
		},
	}

	if configFile != "" {
		if err := readFile(configFile, &cfg); err != nil {
			return nil, err
		}
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func readFile(configFile string, cfg *ControllerManagerConfiguration) error {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}

	defer func() {
		utilruntime.HandleError(f.Close())
	}()

	decoder := yaml.NewDecoder(f)

	return decoder.Decode(cfg)
}

func validateConfig(cfg *ControllerManagerConfiguration) error {
	if cfg.Controllers.Shoot.MaxConcurrentReconciles < 1 {
		fldPath := field.NewPath("controllers", "shootState", "maxConcurrentReconciles")
		return field.Invalid(fldPath, cfg.Controllers.Shoot.MaxConcurrentReconciles, "must be 1 or greater")
	}

	if cfg.Controllers.Shoot.MaxConcurrentReconcilesPerNamespace > cfg.Controllers.Shoot.MaxConcurrentReconciles {
		fldPath := field.NewPath("controllers", "shootState", "maxConcurrentReconcilesPerNamespace")
		return field.Invalid(fldPath, cfg.Controllers.Shoot.MaxConcurrentReconcilesPerNamespace, "must not be greater than maxConcurrentReconciles")
	}

	return nil
}
