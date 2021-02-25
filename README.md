<!-- ABOUT THE PROJECT -->

# Skydive Operator (WIP)

This is a simple operator to deploy skydive analyzer and agents.

<!-- GETTING STARTED -->

## Getting Started

To set up this operator follow the instructions below:

### Prerequisites

* Make sure you have golang installed on your machine, The go-lang version that this operator was built on is go1.15.7
  darwin/amd64

* An openshift cluster
    * If you don't have an openshift cluster you may locally install one
      using [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)

### Installation

1. If you're running a kind cluster the following commands won't work for you and you will have to modify the project
   namespace in the next part

```zsh
oc adm new-project --node-selector='' skydive
oc project skydive
```

2. Run this to grant privileges to your user (#TODO this is maybe not needed)

```zsh
oc adm policy add-scc-to-user privileged -z default
oc adm policy add-cluster-role-to-user cluster-reader -z default
```

3. Clone this repo into your $GOPATH/src folder
4. From skydive-operator folder run:

```bash
make manifests
kubectl create -f config/crd/bases
kubectl create -f config/samples/skydive_v1beta1_skydiveanalyzer.yaml 
kubectl create -f config/samples/skydive_v1beta1_skydiveagents.yaml 
make run
```

You might need to install some golang packages 
(#TODO make sure to see what my IDE auto installed for me)