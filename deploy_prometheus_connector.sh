clear_prometheus_connector_crd() {
  PROMETHEUS_SERVICE=$(oc get service | grep prometheus)
  if [ ! -z "$PROMETHEUS_SERVICE" ]; then
    kubectl delete service prometheus
  fi

  PROMETHEUS_CONNECTOR_SERVICE=$(oc get service | grep skydive-prometheus-connector)
  if [ ! -z "$PROMETHEUS_CONNECTOR_SERVICE" ]; then
    kubectl delete service skydive-prometheus-connector
  fi

  PROMETHEUS_CONNECTOR=$(oc get PrometheusConnector)
  if [ ! -z "$PROMETHEUS_CONNECTOR" ]; then
    kubectl delete -f config/crd/bases/skydive.example.com_prometheusconnectors.yaml
  fi
}

# deploy skydive flow exporter
deploy_prometheus_connector() {
  make manifests
  kubectl create -f config/crd/bases/skydive.example.com_prometheusconnectors.yaml
  kubectl create -f config/skydive_v1_prometheusconnector.yaml
  make run
}

# main
main() {
  clear_prometheus_connector_crd
  deploy_prometheus_connector
}

main