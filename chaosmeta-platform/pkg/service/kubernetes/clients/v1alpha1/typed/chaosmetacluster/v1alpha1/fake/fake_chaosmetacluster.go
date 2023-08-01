package fake

import (
	"chaosmeta-platform/pkg/gateway/apis/chaosmetacluster/v1alpha1"
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeChaosmetaClusters implements ChaosmetaClusterInterface
type FakeChaosmetaClusters struct {
	Fake *FakeChaosmetaclusterV1alpha1
}

var ChaosmetaclustersResource = schema.GroupVersionResource{Group: "Chaosmetacluster", Version: "v1alpha1", Resource: "Chaosmetaclusters"}

var ChaosmetaclustersKind = schema.GroupVersionKind{Group: "Chaosmetacluster", Version: "v1alpha1", Kind: "ChaosmetaCluster"}

// Get takes name of the ChaosmetaCluster, and returns the corresponding ChaosmetaCluster object, and an error if there is any.
func (c *FakeChaosmetaClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(ChaosmetaclustersResource, name), &v1alpha1.ChaosmetaCluster{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ChaosmetaCluster), err
}

// List takes label and field selectors, and returns the list of ChaosmetaClusters that match those selectors.
func (c *FakeChaosmetaClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ChaosmetaClusterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(ChaosmetaclustersResource, ChaosmetaclustersKind, opts), &v1alpha1.ChaosmetaClusterList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ChaosmetaClusterList{ListMeta: obj.(*v1alpha1.ChaosmetaClusterList).ListMeta}
	for _, item := range obj.(*v1alpha1.ChaosmetaClusterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ChaosmetaClusters.
func (c *FakeChaosmetaClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(ChaosmetaclustersResource, opts))
}

// Create takes the representation of a ChaosmetaCluster and creates it.  Returns the server's representation of the ChaosmetaCluster, and an error, if there is any.
func (c *FakeChaosmetaClusters) Create(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.CreateOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(ChaosmetaclustersResource, ChaosmetaCluster), &v1alpha1.ChaosmetaCluster{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ChaosmetaCluster), err
}

// Update takes the representation of a ChaosmetaCluster and updates it. Returns the server's representation of the ChaosmetaCluster, and an error, if there is any.
func (c *FakeChaosmetaClusters) Update(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.UpdateOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(ChaosmetaclustersResource, ChaosmetaCluster), &v1alpha1.ChaosmetaCluster{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ChaosmetaCluster), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeChaosmetaClusters) UpdateStatus(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.UpdateOptions) (*v1alpha1.ChaosmetaCluster, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(ChaosmetaclustersResource, "status", ChaosmetaCluster), &v1alpha1.ChaosmetaCluster{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ChaosmetaCluster), err
}

// Delete takes name of the ChaosmetaCluster and deletes it. Returns an error if one occurs.
func (c *FakeChaosmetaClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(ChaosmetaclustersResource, name), &v1alpha1.ChaosmetaCluster{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeChaosmetaClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(ChaosmetaclustersResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ChaosmetaClusterList{})
	return err
}

// Patch applies the patch and returns the patched ChaosmetaCluster.
func (c *FakeChaosmetaClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ChaosmetaCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(ChaosmetaclustersResource, name, pt, data, subresources...), &v1alpha1.ChaosmetaCluster{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ChaosmetaCluster), err
}
