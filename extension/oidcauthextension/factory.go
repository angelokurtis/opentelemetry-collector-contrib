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

package oidcauthextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	// The value of extension "type" in configuration.
	typeStr = "oidc"

	defaultAttribute = "authorization"
)

// NewFactory creates a factory for the OIDC Authenticator extension.
func NewFactory() component.ExtensionFactory {
	return component.NewExtensionFactory(
		typeStr,
		createDefaultConfig,
		createExtension,
		component.StabilityLevelBeta,
	)
}

func createDefaultConfig() component.ExtensionConfig {
	return &Config{
		ExtensionSettings: config.NewExtensionSettings(component.NewID(typeStr)),
		Attribute:         defaultAttribute,
	}
}

func createExtension(_ context.Context, set component.ExtensionCreateSettings, cfg component.ExtensionConfig) (component.Extension, error) {
	return newExtension(cfg.(*Config), set.Logger)
}
