module ponglehub.co.uk/keycloak-init

go 1.15

require (
	github.com/sirupsen/logrus v1.6.0
	ponglehub.co.uk/envreader v0.0.0-00010101000000-000000000000
)

replace ponglehub.co.uk/envreader => ../../../libraries/golang/envreader
