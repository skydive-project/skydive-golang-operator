clear_skydive_flow_exporter_crd() {
SKYDIVE_FLOW_EXPORTER_CRD=$(kubectl get SkydiveFlowExporter -A)
if [ ! -z "$SKYDIVE_FLOW_EXPORTER_CRD" ]; then
  kubectl delete -f config/crd/bases/skydive.example.com_skydiveflowexporters.yaml
fi
}

# deoploy skydive flow exporter
deploy_skydive_flow_exporter_dev() {
  make manifests
  kubectl create -f config/crd/bases/skydive.example.com_skydiveflowexporters.yaml
  kubectl create -f hack/config/skydive_v1_skydiveflowexporter_dev.yaml
  make run_flow_exporter
}

# main
main() {
  cd ..
  clear_skydive_flow_exporter_crd
  deploy_skydive_flow_exporter_dev
}

main