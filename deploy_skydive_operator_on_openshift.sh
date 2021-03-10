# delete skydive project if it exists (to get new frech deployment)
oc_delete_skydive_project_if_exists() {
PROJECT_SKYDIVE=$(oc get project | grep skydive)
if [ ! -z "$PROJECT_SKYDIVE" ]; then
  oc delete project skydive
  while : ; do
    PROJECT_SKYDIVE=$(oc get project | grep skydive)
    if [ -z "$PROJECT_SKYDIVE" ]; then break; fi
    sleep 1
  done
fi
}

# create and switch context to skydive project
oc_create_skydive_project() {
  oc adm new-project --node-selector='' skydive
  oc project skydive
}

# set credentials for skydive deployment
oc_set_credentials() {
  # analyzer and agent run as privileged container
  oc adm policy add-scc-to-user privileged -z default
  # analyzer need cluster-reader access get all informations from the cluster
  oc adm policy add-cluster-role-to-user cluster-reader -z default
}

oc_clear_skydive_crd() {
SKYDIVE_SUITE_CRD=$(oc get skydive)
if [ ! -z "$SKYDIVE_SUITE_CRD" ]; then
  oc delete -f config/crd/bases/skydive.example.com_skydives.yaml
fi
}

# deoploy skydive
oc_deploy_skydive() {
  make manifests
  oc create -f config/crd/bases/skydive.example.com_skydives.yaml
  oc create -f config/skydive_v1_skydive.yaml
  make run_skydive
}

# main
main() {
  oc_delete_skydive_project_if_exists
  oc_create_skydive_project
  oc_set_credentials
  oc_clear_skydive_crd
  oc_deploy_skydive
}

main