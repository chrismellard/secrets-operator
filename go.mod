module github.com/chrismellard/secret-operator

go 1.13

require (
	github.com/go-logr/logr v0.3.0
	github.com/jenkins-x-plugins/secretfacade v0.0.9
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.2.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/oauth2 v0.0.0-20210113205817-d3ed898aa8a3
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	sigs.k8s.io/controller-runtime v0.8.1
)
