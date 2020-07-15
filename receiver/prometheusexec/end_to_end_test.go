// Copyright 2020, OpenTelemetry Authors
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

package prometheusexec

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config"
	"go.uber.org/zap"
)

func TestEndToEnd(t *testing.T) {
	factories, err := config.ExampleComponents()
	assert.NoError(t, err)

	factory := &Factory{}
	factories.Receivers[factory.Type()] = factory

	config, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Receiver without exec, error returned
	errorReceiverConfig := config.Receivers["prometheus_exec"]
	wrapper := new(zap.NewNop(), errorReceiverConfig.(*Config), nil)

	err = wrapper.Start(context.Background(), nil)
	if err == nil {
		t.Errorf("end_to_end_test.go didn't get error, was expecting one")
	}

	// IN PROGRESS

	// Normal test, make sure the process is restarted by reading from a file while the test code writes to it intermittently
	// receiverConfig := config.Receivers["prometheus_exec/test2/secondary"]
	// wrapper = new(zap.NewNop(), receiverConfig.(*Config), nil)

	// err = wrapper.Start(context.Background(), nil)
	// if err != nil {
	// 	t.Errorf("end_to_end_test.go got error = %v", err)
	// }

	// timestamp := time.Now().UnixNano()
}