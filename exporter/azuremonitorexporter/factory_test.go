// Copyright OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azuremonitorexporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.uber.org/zap"
)

// An inappropriate config
type badConfig struct {
	configmodels.ExporterSettings `mapstructure:",squash"`
}

func TestExporterTypeKey(t *testing.T) {
	factory := Factory{}

	assert.Equal(t, configmodels.Type(typeStr), factory.Type())
}

func TestCreateMetricsExporter(t *testing.T) {
	factory := Factory{}

	ctx := context.Background()
	params := component.ExporterCreateParams{Logger: zap.NewNop()}
	exporter, err := factory.CreateMetricsExporter(ctx, params, &Config{})

	// unsupported
	assert.Nil(t, exporter)
	assert.Equal(t, configerror.ErrDataTypeIsNotSupported, err)
}

func TestCreateTraceExporterUsingSpecificTransportChannel(t *testing.T) {
	// mock transport channel creation
	factory := Factory{TransportChannel: &mockTransportChannel{}}
	ctx := context.Background()
	params := component.ExporterCreateParams{Logger: zap.NewNop()}
	exporter, err := factory.CreateTraceExporter(ctx, params, factory.CreateDefaultConfig())
	assert.NotNil(t, exporter)
	assert.Nil(t, err)
}

func TestCreateTraceExporterUsingDefaultTransportChannel(t *testing.T) {
	// We get the default transport channel creation, if we don't specify one during factory creation
	factory := Factory{}
	assert.Nil(t, factory.TransportChannel)
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	params := component.ExporterCreateParams{Logger: logger}
	exporter, err := factory.CreateTraceExporter(ctx, params, factory.CreateDefaultConfig())
	assert.NotNil(t, exporter)
	assert.Nil(t, err)
	assert.NotNil(t, factory.TransportChannel)
}

func TestCreateTraceExporterUsingBadConfig(t *testing.T) {
	// We get the default transport channel creation, if we don't specify one during factory creation
	factory := Factory{}
	assert.Nil(t, factory.TransportChannel)
	ctx := context.Background()
	params := component.ExporterCreateParams{Logger: zap.NewNop()}

	badConfig := &badConfig{}

	exporter, err := factory.CreateTraceExporter(ctx, params, badConfig)
	assert.Nil(t, exporter)
	assert.NotNil(t, err)
}
