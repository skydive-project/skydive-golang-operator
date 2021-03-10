# delete skydive project if it exists (to get new frech deployment)
kubectl_delete_skydive_project_if_exists() {
NAMESPACE_SKYDIVE=$(kubectl get namespace | grep skydive)
if [ ! -z "$NAMESPACE_SKYDIVE" ]; then
  kubectl delete namespace skydive
fi
}

# create and switch context to skydive project
kubectl_create_skydive_project() {
  kubectl create namespace skydive
}

kubectl_clear_skydive_crd() {
SKYDIVE_SUITE_CRD=$(kubectl get skydive)
if [ ! -z "$SKYDIVE_SUITE_CRD" ]; then
  kubectl delete -f config/crd/bases/skydive.example.com_skydives.yaml
fi
}

kubectl_deploy_skydive() {
  make manifests
  kubectl create -f config/crd/bases/skydive.example.com_skydives.yaml
  kubectl create -f hack/config/skydive_v1_skydive_kind.yaml
  make run_skydive
}

# main
main() {
  cd ..
  kubectl_delete_skydive_project_if_exists
  kubectl_create_skydive_project
  kubectl_clear_skydive_crd
  kubectl_deploy_skydive
}

main