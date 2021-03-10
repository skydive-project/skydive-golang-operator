kubectl_clear_flow_exporter_env() {
FLOW_EXPORTER_ENV_DEPLOY=$(kubectl get deployment -A | grep skydive-flow-exporter-env)
if [ ! -z "$FLOW_EXPORTER_ENV_DEPLOY" ]; then
  kubectl delete -f ../assets/enviourment-deployment.yaml
fi
}

kubectl_create_skydive_project() {
  kubectl apply -f ../assets/enviourment-deployment.yaml
}

# main
main() {
  kubectl_clear_flow_exporter_env
  kubectl_create_skydive_project
}

main