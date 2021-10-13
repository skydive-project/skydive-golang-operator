module skydive

go 1.13

require (
	github.com/go-logr/logr v0.4.0
	github.com/imdario/mergo v0.3.12
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/openshift/api v0.0.0-20200722170803-0ba2c3658da6
	github.com/openshift/client-go v0.0.0-20200722173614-5a1b0aaeff15
	github.com/pkg/errors v0.9.1
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.9.7
)

replace (
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
)
