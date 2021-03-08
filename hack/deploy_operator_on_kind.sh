# delete skydive project if it exists (to get new frech deployment)
kubectl_delete_skydive_project_if_exists() {
NAMESPACE_SKYDIVE=$(kubectl get namespace | grep skydive)
if [ ! -z "$NAMESPACE_SKYDIVE" ]; then
  kubectl delete namespace skydive
  while : ; do
    NAMESPACE_SKYDIVE=$(kubectl get project | grep skydive)
    if [ -z "$PROJECT_SKYDIVE" ]; then break; fi
    sleep 1
  done
fi
}

# create and switch context to skydive project
kubectl_create_skydive_project() {
  kubectl create namespace skydive
}

kubectl_clear_skydivesuite_crd() {
SKYDIVE_SUITE_CRD=$(kubectl get skydivesuite)
if [ ! -z "$SKYDIVE_SUITE_CRD" ]; then
  kubectl delete -f ../config/crd/bases
fi
}

# deoploy skydive
kubectl_deploy_skydive() {
  cd ..
  make manifests
  kubectl create -f config/crd/bases
  kubectl create -f config/skydive_v1_skydivesuite.yaml
  make run
}

# main
main() {
  kubectl_delete_skydive_project_if_exists
  kubectl_create_skydive_project
  kubectl_clear_skydivesuite_crd
  kubectl_deploy_skydive
}

main