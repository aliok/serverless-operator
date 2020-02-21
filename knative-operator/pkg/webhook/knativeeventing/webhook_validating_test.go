package knativeeventing_test

import (
	"context"
	"os"
	"testing"

	. "github.com/openshift-knative/serverless-operator/knative-operator/pkg/webhook/knativeeventing"
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	eventingv1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	configv1.AddToScheme(scheme.Scheme)
}

func TestInvalidNamespace(t *testing.T) {
	os.Clearenv()
	os.Setenv("REQUIRED_EVENTING_NAMESPACE", "knative-eventing")
	validator := KnativeEventingValidator{}
	// The mock will return a KE in the 'default' namespace
	validator.InjectDecoder(&mockDecoder{})
	result := validator.Handle(context.TODO(), types.Request{})
	if result.Response.Allowed {
		t.Error("The required namespace is wrong, but the request is allowed")
	}
}

func TestInvalidVersion(t *testing.T) {
	os.Clearenv()
	os.Setenv("MIN_OPENSHIFT_VERSION", "4.1.13")
	validator := KnativeEventingValidator{}
	validator.InjectDecoder(&mockDecoder{})
	validator.InjectClient(fake.NewFakeClient(mockClusterVersion("3.2.0")))
	result := validator.Handle(context.TODO(), types.Request{})
	if result.Response.Allowed {
		t.Error("The version is too low, but the request is allowed")
	}
}

func mockClusterVersion(version string) *configv1.ClusterVersion {
	return &configv1.ClusterVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "version",
		},
		Status: configv1.ClusterVersionStatus{
			Desired: configv1.Update{
				Version: version,
			},
		},
	}
}

type mockDecoder struct {
	ke *eventingv1alpha1.KnativeEventing
}

var _ types.Decoder = (*mockDecoder)(nil)

func (mock *mockDecoder) Decode(_ types.Request, obj runtime.Object) error {
	if p, ok := obj.(*eventingv1alpha1.KnativeEventing); ok {
		if mock.ke != nil {
			*p = *mock.ke
		} else {
			*p = eventingv1alpha1.KnativeEventing{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-eventing",
					Namespace: "default",
				},
			}
		}
	}
	return nil
}