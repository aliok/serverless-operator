// Code generated by injection-gen. DO NOT EDIT.

package fake

import (
	context "context"

	ingress "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/injection/informers/config/v1/ingress"
	fake "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/injection/informers/factory/fake"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
)

var Get = ingress.Get

func init() {
	injection.Fake.RegisterInformer(withInformer)
}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := fake.Get(ctx)
	inf := f.Config().V1().Ingresses()
	return context.WithValue(ctx, ingress.Key{}, inf), inf.Informer()
}
