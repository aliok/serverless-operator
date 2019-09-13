package test

import (
	"os"
	"os/signal"
	"testing"

	servingversioned "github.com/knative/serving/pkg/client/clientset/versioned"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	servingoperatorversioned "knative.dev/serving-operator/pkg/client/clientset/versioned"
	servingoperatorv1alpha1 "knative.dev/serving-operator/pkg/client/clientset/versioned/typed/serving/v1alpha1"
)

// Context holds objects related to test execution
type Context struct {
	T           *testing.T
	Clients     *Clients
	CleanupList []CleanupFunc
}

// Clients holds instances of interfaces for making requests to various APIs
type Clients struct {
	Kube            *kubernetes.Clientset
	ServingOperator servingoperatorv1alpha1.ServingV1alpha1Interface
	Serving         *servingversioned.Clientset
	OLM             versioned.Interface
	Dynamic         dynamic.Interface
	Config          *rest.Config
}

// CleanupFunc defines a function that is called when the respective resource
// should be deleted. When creating resources the user should also create a CleanupFunc
// and register with the Context
type CleanupFunc func() error

// Setup creates the context object needed in the e2e tests
func Setup(t *testing.T) *Context {
	clients, err := NewClients()
	if err != nil {
		t.Fatalf("Couldn't initialize clients: %v", err)
	}

	ctx := &Context{
		T:       t,
		Clients: clients,
	}
	return ctx
}

// NewClients instantiates and returns several clientsets required for making request to the
// Knative cluster
func NewClients() (*Clients, error) {
	clients := &Clients{}

	cfg, err := clientcmd.BuildConfigFromFlags("", Flags.Kubeconfig)
	if err != nil {
		return nil, err
	}

	// We poll, so set our limits high.
	cfg.QPS = 100
	cfg.Burst = 200

	clients.Kube, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	clients.Dynamic, err = dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	clients.ServingOperator, err = newKnativeServingClients(cfg)
	if err != nil {
		return nil, err
	}

	clients.Serving, err = servingversioned.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	clients.OLM, err = newOLMClient(Flags.Kubeconfig)
	if err != nil {
		return nil, err
	}

	clients.Config = cfg
	return clients, nil
}

func newOLMClient(configPath string) (versioned.Interface, error) {
	olmclient, err := client.NewClient(configPath)
	if err != nil {
		return nil, err
	}
	return olmclient, nil
}

func newKnativeServingClients(cfg *rest.Config) (servingoperatorv1alpha1.ServingV1alpha1Interface, error) {
	cs, err := servingoperatorversioned.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cs.ServingV1alpha1(), nil
}

// Cleanup iterates through the list of registered CleanupFunc functions and calls them
func (ctx *Context) Cleanup() {
	for _, f := range ctx.CleanupList {
		err := f()
		if err != nil {
			ctx.T.Fatalf("Error cleaning up %v", f)
		}
	}
}

// AddToCleanup adds the cleanup function as the first function to the cleanup list,
// we want to delete the last thing first
func (ctx *Context) AddToCleanup(f CleanupFunc) {
	ctx.CleanupList = append([]CleanupFunc{f}, ctx.CleanupList...)
}

// CleanupOnInterrupt will execute the function cleanup if an interrupt signal is caught
func CleanupOnInterrupt(t *testing.T, cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			t.Logf("Test interrupted, cleaning up.")
			cleanup()
			os.Exit(1)
		}
	}()
}
