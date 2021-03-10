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
    * If you wish to run this locally by using kind please check [KindInstallationREADME.md](hack/KindInstallationREADME.md)



### Modyfing the CRD
Modify the file config/skydive_v1_skydivesuite.yaml and pick what you wish to deploy (insert true or false in the relevant field), the options are as follows:
* Skydive agents
* Skydive analyzer
    * Service route
    * Flow exporter
    
Choose your logging level (defaults to INFO)


### Installation - Open-Shift

run the zsh script :
    ```sh
    deploy_operator_on_openshift.sh
    ```

#### Analyzer UI
1a. Check that all the pods, services and routes are running and afterwards run the following command:

```sh
oc get routes
  ```

2a. post the url into your web-browser (make sure you have got an access to the cluster and are not blocked by it's firewall)

1b. If routes option doesn't work, run the following command 
  ```sh
oc port-forward service/skydive-analyzer 8082:8082 --namespace=skydive
  ```
 
2b. Now open web-browser on localhost:8082


