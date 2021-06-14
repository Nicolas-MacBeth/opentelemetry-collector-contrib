module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/stackdriverexporter

go 1.14

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.7
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.2.2-0.20200728233621-2752da7eaab7
	github.com/golang/protobuf v1.4.3
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.6
	go.opentelemetry.io/collector v0.8.1-0.20200815205113-8e5c6065eb0e
	go.opentelemetry.io/otel v0.9.0
	go.uber.org/zap v1.15.0
	google.golang.org/api v0.37.0
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/examples v0.0.0-20200728194956-1c32b02682df // indirect
)
