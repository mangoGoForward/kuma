package builders

import (
	"context"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	model2 "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/test/resources/model"
)

type TrafficTraceBuilder struct {
	*mesh.TrafficTraceResource
}

func TrafficTrace() *TrafficTraceBuilder {
	return &TrafficTraceBuilder{
		TrafficTraceResource: &mesh.TrafficTraceResource{
			Meta: &model.ResourceMeta{
				Mesh: "default",
				Name: "tt-1",
			},
			Spec: &mesh_proto.TrafficTrace{},
		},
	}
}

func (d *TrafficTraceBuilder) Create(s store.ResourceStore) error {
	return s.Create(context.Background(), d.Build(), store.CreateBy(d.Key()))
}

func (d *TrafficTraceBuilder) Key() model2.ResourceKey {
	return model2.MetaToResourceKey(d.GetMeta())
}


func (d *TrafficTraceBuilder) WithName(name string) *TrafficTraceBuilder {
	d.Meta.(*model.ResourceMeta).Name = name
	return d
}

func (d *TrafficTraceBuilder) WithMesh(mesh string) *TrafficTraceBuilder {
	d.Meta.(*model.ResourceMeta).Mesh = mesh
	return d
}

func (d *TrafficTraceBuilder) Build() *mesh.TrafficTraceResource {
	if err := d.TrafficTraceResource.Validate(); err != nil {
		panic(err)
	}
	return d.TrafficTraceResource
}

func (t *TrafficTraceBuilder) WithoutSelectors() *TrafficTraceBuilder {
	t.Spec.Selectors = nil
	return t
}

func (t *TrafficTraceBuilder) WithServiceSelector(service string) *TrafficTraceBuilder {
	return t.WithoutSelectors().AddSelector(mesh_proto.ServiceTag, service)
}

func (t *TrafficTraceBuilder) AddSelector(tags ...string) *TrafficTraceBuilder {
	t.Spec.Selectors = append(t.Spec.Selectors, &mesh_proto.Selector{
		Match: tagsKVToMap(tags),
	})
	return t
}
