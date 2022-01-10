package builders

import (
	"context"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	model2 "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/test/resources/model"
)

//
// SampleDataplane().Inbounds().SampleTransparentProxy()

type DataplaneBuilder struct {
	*mesh.DataplaneResource
}

func Dataplane() *DataplaneBuilder {
	return &DataplaneBuilder{
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
							Port:        80,
							ServicePort: 8080,
							Tags: map[string]string{
								"kuma.io/serivce": "backend",
							},
						},
					},
				},
			},
		},
	}
}

func (d *DataplaneBuilder) WithName(name string) *DataplaneBuilder {
	d.Meta.(*model.ResourceMeta).Name = name
	return d
}

func (d *DataplaneBuilder) WithMesh(mesh string) *DataplaneBuilder {
	d.Meta.(*model.ResourceMeta).Mesh = mesh
	return d
}

func (d *DataplaneBuilder) Key() model2.ResourceKey {
	return model2.MetaToResourceKey(d.Meta)
}

func (d *DataplaneBuilder) Create(manager manager.ResourceManager) error {
	return manager.Create(context.Background(), d.DataplaneResource, store.CreateBy(d.Key()))
}

func (d *DataplaneBuilder) WithTags(tags map[string]string) *DataplaneBuilder {
	d.Spec.Networking.Inbound = []*mesh_proto.Dataplane_Networking_Inbound{
		{
			Port: 80,
			ServicePort: 8080,
			Tags: tags,
		},
	}
	return d
}

func (d *DataplaneBuilder) WithServices(services ...string) *DataplaneBuilder {
	d.Spec.Networking.Inbound = []*mesh_proto.Dataplane_Networking_Inbound{}
	for i, service := range services {
		d.Spec.Networking.Inbound = append(d.Spec.Networking.Inbound, &mesh_proto.Dataplane_Networking_Inbound{
			Port:        uint32(80 + i),
			ServicePort: uint32(8080 + i),
			Tags: map[string]string{
				"kuma.io/service": service,
			},
		})
	}
	return d
}

func (d *DataplaneBuilder) WithSampleTransparentProxy() *DataplaneBuilder {
	d.Spec.Networking.TransparentProxying = &mesh_proto.Dataplane_Networking_TransparentProxying{
		RedirectPortInbound:   0,
		RedirectPortOutbound:  0,
		DirectAccessServices:  nil,
		RedirectPortInboundV6: 0,
	}
	return d
}
