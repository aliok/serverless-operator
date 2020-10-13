// Code generated by injection-gen. DO NOT EDIT.

package fake

import (
	context "context"

	clusteroperator "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/injection/informers/config/v1/clusteroperator"
	fake "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/client/injection/informers/factory/fake"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
)

var Get = clusteroperator.Get

func init() {
	injection.Fake.RegisterInformer(withInformer)
}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := fake.Get(ctx)
	inf := f.Config().V1().ClusterOperators()
	return context.WithValue(ctx, clusteroperator.Key{}, inf), inf.Informer()
}
