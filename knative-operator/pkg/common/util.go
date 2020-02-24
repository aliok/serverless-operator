package common

import (
	servingv1alpha1 "knative.dev/serving-operator/pkg/apis/serving/v1alpha1"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
	"time"
)

const MutationTimestampKey = "knative-serving-openshift/mutation"

var Log = logf.Log.WithName("knative").WithName("openshift")

// config helper to set value for key if not already set
func Configure(ks *servingv1alpha1.KnativeServing, cm, key, value string) bool {
	if ks.Spec.Config == nil {
		ks.Spec.Config = map[string]map[string]string{}
	}
	if _, found := ks.Spec.Config[cm][key]; !found {
		if ks.Spec.Config[cm] == nil {
			ks.Spec.Config[cm] = map[string]string{}
		}
		ks.Spec.Config[cm][key] = value
		Log.Info("Configured", "map", cm, key, value)
		return true
	}
	return false
}

// buildImageOverrideMapFromEnviron creates a map to overrides registry images
func buildImageOverrideMapFromEnviron() map[string]string {
	overrideMap := map[string]string{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "IMAGE_") {
			name := strings.SplitN(pair[0], "_", 2)[1]
			if pair[1] != "" {
				overrideMap[name] = pair[1]
			}
		}
	}
	return overrideMap
}

// Mark the time when instance configured for OpenShift
func annotateTimestamp(annotations map[string]string) {
	annotations[MutationTimestampKey] = time.Now().Format(time.RFC3339)
}
