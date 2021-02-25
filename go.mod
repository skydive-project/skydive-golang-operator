module skydive

go 1.13

require (
	github.com/go-logr/logr v0.3.0
	github.com/imdario/mergo v0.3.10
	github.com/kr/pretty v0.2.1 // indirect
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/openshift/api v0.0.0-20200722170803-0ba2c3658da6
	github.com/openshift/client-go v0.0.0-20200722173614-5a1b0aaeff15
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/tools v0.0.0-20201008025239-9df69603baec // indirect
	k8s.io/api v0.20.2
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.8.2
)

replace (
	k8s.io/api => k8s.io/api v0.19.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.4
	k8s.io/client-go => k8s.io/client-go v0.19.4
)
