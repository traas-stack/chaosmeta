package v1alpha1

import (
	"chaosmeta-platform/pkg/gateway/apis/chaosmetacluster/v1alpha1"
	"chaosmeta-platform/pkg/service/kubernetes/clients/v1alpha1/scheme"
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ChaosmetaClustersGetter has a method to return a ChaosmetaClusterInterface.
// A group's client should implement this interface.
type ChaosmetaClustersGetter interface {
	ChaosmetaClusters() ChaosmetaClusterInterface
}

// ChaosmetaClusterInterface has methods to work with ChaosmetaCluster resources.
type ChaosmetaClusterInterface interface {
	Create(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.CreateOptions) (*v1alpha1.ChaosmetaCluster, error)
	Update(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.UpdateOptions) (*v1alpha1.ChaosmetaCluster, error)
	UpdateStatus(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.UpdateOptions) (*v1alpha1.ChaosmetaCluster, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ChaosmetaCluster, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ChaosmetaClusterList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ChaosmetaCluster, err error)
	ChaosmetaClusterExpansion
}

// ChaosmetaClusters implements ChaosmetaClusterInterface
type ChaosmetaClusters struct {
	client rest.Interface
}

// newChaosmetaClusters returns a ChaosmetaClusters
func newChaosmetaClusters(c *ChaosmetaclusterV1alpha1Client) *ChaosmetaClusters {
	return &ChaosmetaClusters{
		client: c.RESTClient(),
	}
}

// Get takes name of the ChaosmetaCluster, and returns the corresponding ChaosmetaCluster object, and an error if there is any.
func (c *ChaosmetaClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	result = &v1alpha1.ChaosmetaCluster{}
	err = c.client.Get().
		Resource("Chaosmetaclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ChaosmetaClusters that match those selectors.
func (c *ChaosmetaClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ChaosmetaClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ChaosmetaClusterList{}
	err = c.client.Get().
		Resource("Chaosmetaclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested ChaosmetaClusters.
func (c *ChaosmetaClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("Chaosmetaclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a ChaosmetaCluster and creates it.  Returns the server's representation of the ChaosmetaCluster, and an error, if there is any.
func (c *ChaosmetaClusters) Create(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.CreateOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	result = &v1alpha1.ChaosmetaCluster{}
	err = c.client.Post().
		Resource("Chaosmetaclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ChaosmetaCluster).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a ChaosmetaCluster and updates it. Returns the server's representation of the ChaosmetaCluster, and an error, if there is any.
func (c *ChaosmetaClusters) Update(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.UpdateOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	result = &v1alpha1.ChaosmetaCluster{}
	err = c.client.Put().
		Resource("Chaosmetaclusters").
		Name(ChaosmetaCluster.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ChaosmetaCluster).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *ChaosmetaClusters) UpdateStatus(ctx context.Context, ChaosmetaCluster *v1alpha1.ChaosmetaCluster, opts v1.UpdateOptions) (result *v1alpha1.ChaosmetaCluster, err error) {
	result = &v1alpha1.ChaosmetaCluster{}
	err = c.client.Put().
		Resource("Chaosmetaclusters").
		Name(ChaosmetaCluster.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ChaosmetaCluster).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the ChaosmetaCluster and deletes it. Returns an error if one occurs.
func (c *ChaosmetaClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("Chaosmetaclusters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *ChaosmetaClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("Chaosmetaclusters").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched ChaosmetaCluster.
func (c *ChaosmetaClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ChaosmetaCluster, err error) {
	result = &v1alpha1.ChaosmetaCluster{}
	err = c.client.Patch(pt).
		Resource("Chaosmetaclusters").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
