package samples

import (
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/test/resources/builders"
)

func TrafficTraceBackend() *mesh.TrafficTraceResource {
	return TrafficTraceBackendBuilder().Build()
}

func TrafficTraceBackendBuilder() *builders.TrafficTraceBuilder {
	return builders.TrafficTrace().
		WithName("tt-backend").
		WithServiceSelector("backend")
}

func TrafficTraceWeb() *mesh.TrafficTraceResource {
	return TrafficTraceWebBuilder().Build()
}

func TrafficTraceWebBuilder() *builders.TrafficTraceBuilder {
	return builders.TrafficTrace().
		WithName("tt-web").
		WithServiceSelector("web")
}
