#!/bin/bash

set -Eeuo pipefail

template="${1:?Provide template file as arg[1]}"
target="${2:?Provide a target CSV file as arg[2]}"

# shellcheck disable=SC1091,SC1090
source "$(dirname "${BASH_SOURCE[0]}")/../lib/metadata.bash"

registry="registry.svc.ci.openshift.org/openshift"
serving="${registry}/knative-v$(metadata.get dependencies.serving):knative-serving"
eventing="${registry}/knative-v$(metadata.get dependencies.eventing):knative-eventing"
eventing_contrib="${registry}/knative-v$(metadata.get dependencies.eventing_contrib):knative-eventing-sources"

declare -a serving_images
declare -a eventing_images
declare -a kafka_images
declare -A serving_images_addresses
declare -A eventing_images_addresses
declare -A kafka_images_addresses

function serving_image {
  local name address
  name="${1:?Pass a image name as arg[1]}"
  address="${2:?Pass a image address as arg[2]}"
  serving_images+=("${name}")
  serving_images_addresses["${name}"]="${address}"
}

function eventing_image {
  local name address
  name="${1:?Pass a image name as arg[1]}"
  address="${2:?Pass a image address as arg[2]}"
  eventing_images+=("${name}")
  eventing_images_addresses["${name}"]="${address}"
}

function kafka_image {
  local name address
  name="${1:?Pass a image name as arg[1]}"
  address="${2:?Pass a image address as arg[2]}"
  kafka_images+=("${name}")
  kafka_images_addresses["${name}"]="${address}"
}

serving_image "queue-proxy"    "${serving}-queue"
serving_image "activator"      "${serving}-activator"
serving_image "autoscaler"     "${serving}-autoscaler"
serving_image "autoscaler-hpa" "${serving}-autoscaler-hpa"
serving_image "controller"     "${serving}-controller"
serving_image "webhook"        "${serving}-webhook"
serving_image "storage-version-migration-serving-$(metadata.get dependencies.serving)__migrate" "${serving}-storage-version-migration"

serving_image "3scale-kourier-gateway" "docker.io/maistra/proxyv2-ubi8:$(metadata.get dependencies.maistra)"
serving_image "3scale-kourier-control" "${registry}/knative-v$(metadata.get dependencies.kourier):kourier"

serving_image "KN_CLI_ARTIFACTS"     "${registry}/knative-v$(metadata.get dependencies.cli):kn-cli-artifacts"

eventing_image "eventing-controller__eventing-controller"    "${eventing}-controller"
eventing_image "sugar-controller__controller"                "${eventing}-sugar-controller"
eventing_image "eventing-webhook__eventing-webhook"          "${eventing}-webhook"
eventing_image "storage-version-migration-eventing__migrate" "${eventing}-storage-version-migration"
eventing_image "mt-broker-controller__mt-broker-controller"  "${eventing}-mtchannel-broker"
eventing_image "mt-broker-filter__filter"                    "${eventing}-mtbroker-filter"
eventing_image "mt-broker-ingress__ingress"                  "${eventing}-mtbroker-ingress"
eventing_image "imc-controller__controller"                  "${eventing}-channel-controller"
eventing_image "imc-dispatcher__dispatcher"                  "${eventing}-channel-dispatcher"

eventing_image "v0.17.0-pingsource-cleanup__pingsource" "${eventing}-pingsource-cleanup"
eventing_image "PING_IMAGE"           "${eventing}-ping"
eventing_image "MT_PING_IMAGE"        "${eventing}-mtping"
eventing_image "APISERVER_RA_IMAGE"   "${eventing}-apiserver-receive-adapter"
eventing_image "BROKER_INGRESS_IMAGE" "${eventing}-broker-ingress"
eventing_image "BROKER_FILTER_IMAGE"  "${eventing}-broker-filter"
eventing_image "DISPATCHER_IMAGE"     "${eventing}-channel-dispatcher"

kafka_image "kafka-controller-manager__manager"    "${eventing_contrib}-kafka-source-controller"
kafka_image "KAFKA_RA_IMAGE"                       "${eventing_contrib}-kafka-source-adapter"
kafka_image "kafka-ch-controller__controller"      "${eventing_contrib}-kafka-channel-controller"
kafka_image "DISPATCHER_IMAGE"                     "${eventing_contrib}-kafka-channel-dispatcher"
kafka_image "kafka-ch-dispatcher__dispatcher"      "${eventing_contrib}-kafka-channel-dispatcher"
kafka_image "kafka-webhook__kafka-webhook"         "${eventing_contrib}-kafka-channel-webhook"

declare -A values
values[spec.version]="$(metadata.get project.version)"
values[metadata.name]="$(metadata.get project.name).v$(metadata.get project.version)"
values['metadata.annotations[olm.skipRange]']="$(metadata.get olm.skipRange)"
values[spec.minKubeVersion]="$(metadata.get requirements.kube.minVersion)"
values[spec.replaces]="$(metadata.get project.name).v$(metadata.get olm.replaces)"

function add_related_image {
  cat << EOF | yq write --inplace --script - "$1"
- command: update
  path: spec.relatedImages[+]
  value:
    name: "IMAGE_${2}_${3}"
    image: "${4}"
EOF
}

function add_downstream_operator_deployment_image {
  cat << EOF | yq write --inplace --script - "$1"
- command: update
  path: spec.install.spec.deployments(name==knative-openshift).spec.template.spec.containers[0].env[+]
  value:
    name: "IMAGE_${2}_${3}"
    value: "${4}"
EOF
}

# since we also parse the environment variables in the upstream (actually midstream) operator,
# we don't add scope prefixes to image overrides here. We don't have a clash anyway without any scope prefixes!
# there was a naming clash between eventing and kafka, but we won't provide the Kafka overrides to the
# midstream operator.
function add_upstream_operator_deployment_image {
  cat << EOF | yq write --inplace --script - "$1"
- command: update
  path: spec.install.spec.deployments(name==knative-operator).spec.template.spec.containers[0].env[+]
  value:
    name: "IMAGE_${2}"
    value: "${3}"
EOF
}

# Start fresh
cp "$template" "$target"

for name in "${serving_images[@]}"; do
  echo "serving Image: ${name} -> ${serving_images_addresses[$name]}"
  add_related_image "$target" "SERVING" "$name" "${serving_images_addresses[$name]}"
  add_downstream_operator_deployment_image "$target" "SERVING" "$name" "${serving_images_addresses[$name]}"
  add_upstream_operator_deployment_image "$target" "$name" "${serving_images_addresses[$name]}"
done

for name in "${eventing_images[@]}"; do
  echo "eventing Image: ${name} -> ${eventing_images_addresses[$name]}"
  add_related_image "$target" "EVENTING" "$name" "${eventing_images_addresses[$name]}"
  add_downstream_operator_deployment_image "$target" "EVENTING" "$name" "${eventing_images_addresses[$name]}"
  add_upstream_operator_deployment_image "$target" "$name" "${eventing_images_addresses[$name]}"
done

# don't add Kafka image overrides to upstream operator
for name in "${kafka_images[@]}"; do
  echo "kafka Image: ${name} -> ${kafka_images_addresses[$name]}"
  add_related_image "$target" "KAFKA" "$name" "${kafka_images_addresses[$name]}"
  add_downstream_operator_deployment_image "$target" "KAFKA" "$name" "${kafka_images_addresses[$name]}"
done

for name in "${!values[@]}"; do
  echo "Value: ${name} -> ${values[$name]}"
  yq write --inplace "$target" "$name" "${values[$name]}"
done
