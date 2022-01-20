package builders

import (
	"context"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	model2 "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/test/resources/model"
)

type MeshBuilder struct {
	*mesh.MeshResource
}

func Mesh() *MeshBuilder {
	return &MeshBuilder{
		MeshResource: &mesh.MeshResource{
			Meta: &model.ResourceMeta{
				Mesh: "",
				Name: "default",
			},
			Spec: &mesh_proto.Mesh{
				Mtls: &mesh_proto.Mesh_Mtls{},
			},
		},
	}
}

func (m *MeshBuilder) Build() *mesh.MeshResource {
	if err := m.MeshResource.Validate(); err != nil {
		panic(err)
	}
	return m.MeshResource
}

func (m *MeshBuilder) Create(s store.ResourceStore) error {
	return s.Create(context.Background(), m.Build(), store.CreateBy(m.Key()))
}

func (m *MeshBuilder) Key() model2.ResourceKey {
	return model2.MetaToResourceKey(m.GetMeta())
}

func (m *MeshBuilder) WithBuiltinMTLSBackend(name string) *MeshBuilder {
	return m.AddBuiltinMTLSBackend(name)
}

func (m *MeshBuilder) AddBuiltinMTLSBackend(name string) *MeshBuilder {
	m.Spec.Mtls.Backends = append(m.Spec.Mtls.Backends, &mesh_proto.CertificateAuthorityBackend{
		Name: "builtin-1",
		Type: "builtin",
	})
	return m
}

func (m *MeshBuilder) WithEnabledBackend(name string) *MeshBuilder {
	m.Spec.Mtls.EnabledBackend = name
	return m
}
