module ponglehub.co.uk/auth/auth-operator

go 1.16

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.3
	ponglehub.co.uk/lib/user-events v1.0.0
)

replace ponglehub.co.uk/lib/user-events => ./../../libraries/golang/user-events

replace ponglehub.co.uk/lib/events => ./../../libraries/golang/events
