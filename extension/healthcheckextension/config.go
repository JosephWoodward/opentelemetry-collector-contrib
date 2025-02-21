// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package healthcheckextension

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
)

// Config has the configuration for the extension enabling the health check
// extension, used to report the health status of the service.
type Config struct {
	config.ExtensionSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	// Port is the port used to publish the health check status.
	// The default value is 13133.
	// Deprecated: use Endpoint instead.
	Port uint16 `mapstructure:"port"`

	// TCPAddr represents a tcp endpoint address that is to publish the health
	// check status.
	// The default endpoint is "0.0.0.0:13133".
	TCPAddr confignet.TCPAddr `mapstructure:",squash"`

	// CheckCollectorPipeline contains the list of settings of collector pipeline health check
	CheckCollectorPipeline checkCollectorPipelineSettings `mapstructure:"check_collector_pipeline"`
}

var _ config.Extension = (*Config)(nil)
var (
	errNoEndpointProvided                      = errors.New("bad config: endpoint must be specified")
	errInvalidExporterFailureThresholdProvided = errors.New("bad config: exporter_failure_threshold expects a positive number")
)

// Validate checks if the extension configuration is valid
func (cfg *Config) Validate() error {
	_, err := time.ParseDuration(cfg.CheckCollectorPipeline.Interval)
	if err != nil {
		return err
	}
	if cfg.TCPAddr.Endpoint == "" {
		return errNoEndpointProvided
	}
	if cfg.CheckCollectorPipeline.ExporterFailureThreshold <= 0 {
		return errInvalidExporterFailureThresholdProvided
	}
	return nil
}

type checkCollectorPipelineSettings struct {
	// Enabled indicates whether to not enable collector pipeline check.
	Enabled bool `mapstructure:"enabled"`
	// Interval the time range to check healthy status of collector pipeline
	Interval string `mapstructure:"interval"`
	// ExporterFailureThreshold is the threshold of exporter failure numbers during the Interval
	ExporterFailureThreshold int `mapstructure:"exporter_failure_threshold"`
}
