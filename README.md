<!-- ABOUT THE PROJECT -->

# Skydive operator
This is an operator to deploy skydive analyzer, agents, flow-exporter and promethues connector.

## Skydive 
Skydive is an open source real-time network topology and protocols analyzer. It aims to provide a comprehensive way of understanding what is happening in the network infrastructure.

Skydive agents collect topology information and flows and forward them to a central agent for further analysis. All the information is stored in an Elasticsearch database.

Skydive is SDN-agnostic but provides SDN drivers in order to enhance the topology and flows information.

## Skydive flow-exporter
The Skydive Flow Exporter provides a framework for building pipelines which extract flows from the Skydive Analyzer (via it WebSocket API), process them and send the results upstream. read more [here](https://github.com/skydive-project/skydive-flow-exporter)

##Prometheus-connector
Prometheus is an open-source systems monitoring and alerting tool. Many different types of data are collected by different tools and are forwarded to Prometheus via various types of exporters. For example, the Prometheus Node Exporter exposes a wide variety of hardware- and kernel-related metrics. We developed an exporter for Skydive that reports metrics of individual captured network flows. The connector translates data from Skydive captured flows into a format that can be consumed by Prometheus. The first implementation of the Skydive-Prometheus connector periodically provides the byte transfer counts for each network connection under observation. The code can be easily tailored to provide additional flow information. read more [here](https://github.com/skydive-project/skydive.network/blob/3fcd5f66c6926d96460af7e7420339e69480cad1/blog/prometheus-connector.md)

<!-- GETTING STARTED -->

## Getting Started

To set up this operator follow the instructions below:

### Prerequisites

* Make sure you have golang installed on your machine, The go-lang version that this operator was built on is go1.15.7
  darwin/amd64

* An openshift cluster
    * If you wish to run this locally by using kind please
      check [KindInstallationREADME.md](hack/KindInstallationREADME.md)
      
#### Analyzer UI

1a. If you've deployed the operator into an openshift cluster and enabled route in the CRD instance, check that all the pods, services and routes are running and afterwards run the following command:

```sh
oc get routes
  ```

2a. paste the url into your web-browser (make sure you have got access to the cluster and are not blocked by its firewall)

1b. If routes option doesn't work, run the following command

  ```sh
kubectl port-forward service/skydive-analyzer 8082:8082 --namespace=skydive
  ```

2b. Now open web-browser (from the machine you used the kubectl command) on localhost:8082