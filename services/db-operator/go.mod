module ponglehub.co.uk/operators/db

go 1.16

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	ponglehub.co.uk/lib/postgres v1.0.0
)

replace ponglehub.co.uk/lib/postgres => ./../../libraries/golang/postgres
