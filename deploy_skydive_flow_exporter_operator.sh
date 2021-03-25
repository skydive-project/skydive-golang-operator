clear_skydive_flow_exporter_crd() {
SKYDIVE_FLOW_EXPORTER_CRD=$(oc get SkydiveFlowExporter)
if [ ! -z "$SKYDIVE_FLOW_EXPORTER_CRD" ]; then
  kubectl delete -f config/crd/bases/skydive.example.com_skydiveflowexporters.yaml
fi
}

# deploy skydive flow exporter
deploy_skydive_flow_exporter() {
  make manifests
  kubectl create -f config/crd/bases/skydive.example.com_skydiveflowexporters.yaml
  kubectl create -f config/skydive_v1_skydiveflowexporter.yaml
  make run
}

# main
main() {
  clear_skydive_flow_exporter_crd
  deploy_skydive_flow_exporter
}

main