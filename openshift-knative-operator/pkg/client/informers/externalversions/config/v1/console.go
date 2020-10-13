// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	versioned "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/listers/config/v1"
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ConsoleInformer provides access to a shared informer and lister for
// Consoles.
type ConsoleInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ConsoleLister
}

type consoleInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewConsoleInformer constructs a new informer for Console type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewConsoleInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredConsoleInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredConsoleInformer constructs a new informer for Console type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredConsoleInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().Consoles().List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().Consoles().Watch(options)
			},
		},
		&configv1.Console{},
		resyncPeriod,
		indexers,
	)
}

func (f *consoleInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredConsoleInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *consoleInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv1.Console{}, f.defaultInformer)
}

func (f *consoleInformer) Lister() v1.ConsoleLister {
	return v1.NewConsoleLister(f.Informer().GetIndexer())
}
