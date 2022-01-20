package samples

import (
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/test/resources/builders"
)

func MeshDefault() *mesh.MeshResource {
	return MeshDefaultBuilder().Build()
}

func MeshDefaultBuilder() *builders.MeshBuilder {
	return builders.Mesh()
}

func MeshMTLS() *mesh.MeshResource {
	return builders.Mesh().
		AddBuiltinMTLSBackend("builtin-1").
		WithEnabledBackend("builtin-1").
		Build()
}

func MeshMTLSBuilder() *builders.MeshBuilder {
	return builders.Mesh()
}
