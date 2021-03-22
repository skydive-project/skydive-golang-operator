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
  

### Installation

#### Kind
Install kind on your machine, use [this](https://kind.sigs.k8s.io/docs/user/quick-start/) quick start guide for help, 
make sure that your have a running kind cluster and that kubectl commands refer to relevant context.

After installing and running the Kind cluster, make sure the env var KUBECONFIG is set to the correct kubeConfig file,
it usually defaults to $HOME/.kube/config, run this command:

```bash
export KUBECONFIG=~/.kube/config
```

run the script :

```./deploy_operator_on_kind.sh```

  *if you get an error 

### Customize skydive deployment / CRD

Modify the current config/skydive_v1_skydive.yaml and pick what you wish to deploy (insert true or false in the relevant
field), the options are as follows:

* Skydive agents
* Skydive analyzer

You can provide the skydive operator with environments variables in order to customize your skydive deployment.
Checkout [this example](config/skydive_v1_skydive_env_example.yaml) of crd to get started with providing environment
variables to the skydive operator, full list of acceptable enviorment variables are listed [here](https://github.com/skydive-project/skydive/blob/master/etc/skydive.yml.default)


#### Analyzer UI
1. Check that all the pods and services are running and afterwards run the following command:

  ```sh
kubectl port-forward service/skydive-analyzer 8082:8082 --namespace=skydive
  ```

2. Now open browser on localhost:8082
