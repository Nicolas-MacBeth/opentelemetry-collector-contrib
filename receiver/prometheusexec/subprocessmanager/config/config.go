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

package config

// SubprocessConfig is the config definition for the subprocess manager
type SubprocessConfig struct {
	// CommandString is the command to be run (binary + flags)
	CommandString string `mapstructure:"exec"`
	// Port is the port assigned to the Receiver, and to the {{port}} template variables
	Port int `mapstructure:"port"`
	// CustomName is a custom user-specified name to keep track of a certain process in logs
	CustomName string `mapstructure:"custom_name"`
}
