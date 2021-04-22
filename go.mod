module skydive-operator

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/imdario/mergo v0.3.12
	github.com/onsi/ginkgo v1.16.0
	github.com/onsi/gomega v1.11.0
	github.com/openshift/api v3.9.0+incompatible
	github.com/openshift/client-go v0.0.0-20200722173614-5a1b0aaeff15
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.20.5
	k8s.io/apiextensions-apiserver v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.8.3
)
