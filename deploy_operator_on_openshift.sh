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

# deoploy skydive
oc_deploy_skydive() {
  make manifests
  oc create -f config/crd/bases
  oc create -f config/skydive_v1_skydivesuite.yaml
  make run
}

# main
main() {
  echo "starting"
  oc_delete_skydive_project_if_exists
  oc_create_skydive_project
  oc_set_credentials
  oc_deploy_skydive
}
main