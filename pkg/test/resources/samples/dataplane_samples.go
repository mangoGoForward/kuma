package samples

import (
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/test/resources/builders"
)

func DataplaneBackend() *mesh.DataplaneResource {
	return DataplaneBackendBuilder().Build()
}

func DataplaneBackendBuilder() *builders.DataplaneBuilder {
	return builders.Dataplane().
		WithAddress("192.168.0.1").
		WithServices("backend")
}
