// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllogs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/common"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/logs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/metrics"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/traces"
)

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	TraceStatements  []common.ContextStatements `mapstructure:"trace_statements"`
	MetricStatements []common.ContextStatements `mapstructure:"metric_statements"`
	LogStatements    []common.ContextStatements `mapstructure:"log_statements"`

	// Deprecated.  Use TraceStatements, MetricStatements, and LogStatements instead
	OTTLConfig `mapstructure:",squash"`
}

type OTTLConfig struct {
	Traces  SignalConfig `mapstructure:"traces"`
	Metrics SignalConfig `mapstructure:"metrics"`
	Logs    SignalConfig `mapstructure:"logs"`
}

type SignalConfig struct {
	Statements []string `mapstructure:"statements"`
}

var _ component.ProcessorConfig = (*Config)(nil)

func (c *Config) Validate() error {
	if (len(c.Traces.Statements) > 0 || len(c.Metrics.Statements) > 0 || len(c.Logs.Statements) > 0) &&
		(len(c.TraceStatements) > 0 || len(c.MetricStatements) > 0 || len(c.LogStatements) > 0) {
		return fmt.Errorf("cannot use Traces, Metrics and/or Logs with TraceStatements, MetricStatements and/or LogStatements")
	}

	if len(c.Traces.Statements) > 0 {
		ottlspanp := ottlspan.NewParser(traces.Functions(), component.TelemetrySettings{Logger: zap.NewNop()})
		_, err := ottlspanp.ParseStatements(c.Traces.Statements)
		if err != nil {
			return err
		}
	}

	if len(c.TraceStatements) > 0 {
		pc, err := common.NewTraceParserCollection(traces.Functions(), component.TelemetrySettings{Logger: zap.NewNop()})
		if err != nil {
			return err
		}
		for _, cs := range c.TraceStatements {
			_, err = pc.ParseContextStatements(cs)
			if err != nil {
				return err
			}
		}
	}

	if len(c.Metrics.Statements) > 0 {
		ottldatapointp := ottldatapoint.NewParser(metrics.Functions(), component.TelemetrySettings{Logger: zap.NewNop()})
		_, err := ottldatapointp.ParseStatements(c.Metrics.Statements)
		if err != nil {
			return err
		}
	}

	if len(c.MetricStatements) > 0 {
		pc, err := common.NewMetricParserCollection(metrics.Functions(), component.TelemetrySettings{Logger: zap.NewNop()})
		if err != nil {
			return err
		}
		for _, cs := range c.MetricStatements {
			_, err = pc.ParseContextStatements(cs)
			if err != nil {
				return err
			}
		}
	}

	if len(c.Logs.Statements) > 0 {
		ottllogsp := ottllogs.NewParser(logs.Functions(), component.TelemetrySettings{Logger: zap.NewNop()})
		_, err := ottllogsp.ParseStatements(c.Logs.Statements)
		if err != nil {
			return err
		}
	}

	if len(c.LogStatements) > 0 {
		pc, err := common.NewLogParserCollection(logs.Functions(), component.TelemetrySettings{Logger: zap.NewNop()})
		if err != nil {
			return err
		}
		for _, cs := range c.LogStatements {
			_, err = pc.ParseContextStatements(cs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
