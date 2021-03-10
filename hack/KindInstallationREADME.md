<!-- ABOUT THE PROJECT -->

# Running Skydive Operator on Kind cluster

This is a simple deployment guide for the operator on a kind cluster

<!-- GETTING STARTED -->

## Getting Started

To set up this operator on a kind follow the instructions below:

### Prerequisites

* Make sure you have golang installed on your machine, The go-lang version that this operator was built on is go1.15.7
  darwin/amd64

* Kind installed

### Modyfing the CRD
Modify the file config/skydive_v1_skydive.yaml and pick what you wish to deploy (insert true or false in the relevant field), the options are as follows:
* Skydive agents
* Skydive analyzer
    * Service route

Change Env var : (TODO: more info to be added here)

Choose your logging level (defaults to DEBUG)

### Installation

#### Kind
Install kind on your machine, use [this](https://kind.sigs.k8s.io/docs/user/quick-start/) quick start guide for help, 
make sure that your have a running kind cluster and that kubectl commands refer to relevant context.

run the script :

```bash deploy_operator_on_kind.sh```

#### Analyzer UI
1. Check that all the pods and services are running and afterwards run the following command:

  ```sh
kubectl port-forward service/skydive-analyzer 8082:8082 --namespace=skydive
  ```

2. Now open browser on localhost:8082
