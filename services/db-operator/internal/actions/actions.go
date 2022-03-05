package actions

import "ponglehub.co.uk/operators/db/internal/deployments"

type AddStatefulSet struct {
	StatefulSet deployments.StatefulSet
}

type DeleteStatefulSet struct {
	StatefulSet deployments.StatefulSet
}

type AddService struct {
	Service deployments.Service
}

type DeleteService struct {
	Service deployments.Service
}
