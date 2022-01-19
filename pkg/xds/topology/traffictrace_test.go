package topology_test

import (
	"context"

	"github.com/kumahq/kuma/pkg/test/resources/samples"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	resources_manager "github.com/kumahq/kuma/pkg/core/resources/manager"
	plugins_memory "github.com/kumahq/kuma/pkg/plugins/resources/memory"
	"github.com/kumahq/kuma/pkg/xds/topology"
)

var _ = Describe("GetTrafficTrace", func() {

	It("should return matched TrafficTrace", func() {
		// given
		store := plugins_memory.NewStore()
		manager := resources_manager.NewResourceManager(store)

		Expect(samples.TrafficTraceBackendBuilder().Create(store)).To(Succeed())
		Expect(samples.TrafficTraceWebBuilder().Create(store)).To(Succeed())

		// when
		picked, err := topology.GetTrafficTrace(context.Background(), samples.DataplaneBackend(), manager)

		// then
		Expect(err).ToNot(HaveOccurred())

		// and
		Expect(picked.Meta.GetName()).To(Equal(samples.TrafficTraceBackend().GetMeta().GetName()))
	})

	It("should return nil when there are no matching traffic traces", func() {
		// given
		store := plugins_memory.NewStore()
		manager := resources_manager.NewResourceManager(store)

		// when
		picked, err := topology.GetTrafficTrace(context.Background(), samples.DataplaneBackend(), manager)

		// then
		Expect(err).ToNot(HaveOccurred())

		// and
		Expect(picked).To(BeNil())
	})
})
