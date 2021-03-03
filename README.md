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


### Modyfing the CRD
Modify the file config/skydive_v1_skydivesuite.yaml and pick what you wish to deploy (insert true or false in the relevant field), the options are as follows:
* Skydive agents
* Skydive analyzer
    * Service route - TODO: doesn't work on IBM cluster.   
    
Choose your logging level (defaults to INFO)


### Installation - Open-Shift

run the zsh script :
    ```deploy_operator_on_openshift.sh```

### Installation - Kubernetes
*** to be continued
