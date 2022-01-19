package samples

import (
	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/test/resources/builders"
	"github.com/kumahq/kuma/pkg/test/resources/model"
)

func DataplaneBackend() *mesh.DataplaneResource {
	return DataplaneBackendBuilder().Build()
}

func DataplaneBackendBuilder() *builders.DataplaneBuilder {
	return &builders.DataplaneBuilder{
		DataplaneResource: &mesh.DataplaneResource{
			Meta: &model.ResourceMeta{
				Mesh: "default",
				Name: "dp-1",
			},
			Spec: &mesh_proto.Dataplane{
				Networking: &mesh_proto.Dataplane_Networking{
					Address: "192.168.0.1",
					Inbound: []*mesh_proto.Dataplane_Networking_Inbound{
						{
							Port:        builders.FirstInboundPort,
							ServicePort: builders.FirstInboundServicePort,
							Tags: map[string]string{
								mesh_proto.ServiceTag: "backend",
							},
						},
					},
				},
			},
		},
	}
}
