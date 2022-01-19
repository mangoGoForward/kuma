package builders

import (
	"context"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	model2 "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/test/resources/model"
	"github.com/kumahq/kuma/pkg/util/proto"
)

var FirstInboundPort = uint32(80)
var FirstInboundServicePort = uint32(8080)
var FirstOutboundPort = uint32(10001)

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
							Port:        FirstInboundPort,
							ServicePort: FirstInboundServicePort,
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

func (d *DataplaneBuilder) Build() *mesh.DataplaneResource {
	if err := d.DataplaneResource.Validate(); err != nil {
		panic(err)
	}
	return d.DataplaneResource
}

func (d *DataplaneBuilder) Create(s store.ResourceStore) error {
	return s.Create(context.Background(), d.Build(), store.CreateBy(d.Key()))
}

func (d *DataplaneBuilder) Key() model2.ResourceKey {
	return model2.MetaToResourceKey(d.GetMeta())
}

func (d *DataplaneBuilder) WithName(name string) *DataplaneBuilder {
	d.Meta.(*model.ResourceMeta).Name = name
	return d
}

func (d *DataplaneBuilder) WithMesh(mesh string) *DataplaneBuilder {
	d.Meta.(*model.ResourceMeta).Mesh = mesh
	return d
}

func (d *DataplaneBuilder) WithoutInbounds() *DataplaneBuilder {
	d.Spec.Networking.Inbound = nil
	return d
}

func (d *DataplaneBuilder) WithAddress(address string) *DataplaneBuilder {
	d.DataplaneResource.Spec.Networking.Address = address
	return d
}

func (d *DataplaneBuilder) WithTags(tagsKV ...string) *DataplaneBuilder {
	return d.WithTagsMap(tagsKVToMap(tagsKV))
}

func (d *DataplaneBuilder) With(fn func(*mesh.DataplaneResource)) *DataplaneBuilder {
	fn(d.DataplaneResource)
	return d
}

func (d *DataplaneBuilder) WithTagsMap(tags map[string]string) *DataplaneBuilder {
	return d.WithoutInbounds().AddTagsMap(tags)
}

func (d *DataplaneBuilder) AddTags(tags ...string) *DataplaneBuilder {
	return d.AddTagsMap(tagsKVToMap(tags))
}

func (d *DataplaneBuilder) AddTagsMap(tags map[string]string) *DataplaneBuilder {
	d.Spec.Networking.Inbound = append(d.Spec.Networking.Inbound, &mesh_proto.Dataplane_Networking_Inbound{
		Port:        FirstInboundPort + uint32(len(d.Spec.Networking.Inbound)),
		ServicePort: FirstInboundServicePort + uint32(len(d.Spec.Networking.Inbound)),
		Tags:        tags,
	})
	return d
}

func (d *DataplaneBuilder) AddOutboundToService(service string) *DataplaneBuilder {
	d.Spec.Networking.Outbound = append(d.Spec.Networking.Outbound, &mesh_proto.Dataplane_Networking_Outbound{
		Port:    FirstOutboundPort + uint32(len(d.Spec.Networking.Outbound)),
		Tags: map[string]string{
			mesh_proto.ServiceTag: service,
		},
	})
	return d
}

func (d *DataplaneBuilder) WithServices(services ...string) *DataplaneBuilder {
	d.WithoutInbounds()
	for _, service := range services {
		d.AddTags(mesh_proto.ServiceTag, service)
	}
	return d
}

func tagsKVToMap(tagsKV []string) map[string]string {
	if len(tagsKV) % 2 == 1 {
		panic("tagsKV has to have even number of arguments")
	}
	tags := map[string]string{}
	for i := 0; i < len(tagsKV); i += 2 {
		tags[tagsKV[i]] = tagsKV[i+1]
	}
	return tags
}

func (d *DataplaneBuilder) WithPrometheusMetrics(config *mesh_proto.PrometheusMetricsBackendConfig) *DataplaneBuilder {
	d.Spec.Metrics = &mesh_proto.MetricsBackend{
		Name: "prometheus-1",
		Type: mesh_proto.MetricsPrometheusType,
		Conf: proto.MustToStruct(config),
	}
	return d
}
