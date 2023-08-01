package fake

import (
	"chaosmeta-platform/pkg/service/kubernetes/clients/v1alpha1/typed/chaosmetacluster/v1alpha1"

	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeChaosmetaclusterV1alpha1 struct {
	*testing.Fake
}

func (c *FakeChaosmetaclusterV1alpha1) ChaosmetaClusters() v1alpha1.ChaosmetaClusterInterface {
	return &FakeChaosmetaClusters{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeChaosmetaclusterV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
