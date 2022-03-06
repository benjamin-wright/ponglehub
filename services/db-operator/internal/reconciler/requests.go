package reconciler

// import (
// 	"ponglehub.co.uk/operators/db/internal/crds"
// 	"ponglehub.co.uk/operators/db/internal/database"
// 	"ponglehub.co.uk/operators/db/internal/deployments"
// )

// func getDatabaseRequests(databases map[string]crds.Database) (map[string]deployments.StatefulSet, map[string]deployments.Service) {
// 	sets := map[string]deployments.StatefulSet{}
// 	services := map[string]deployments.Service{}

// 	for _, db := range databases {
// 		set := deployments.StatefulSet{
// 			Name:      db.Name,
// 			Namespace: db.Namespace,
// 			Storage:   db.Storage,
// 		}
// 		setKey := set.Key()
// 		sets[setKey] = set

// 		service := deployments.Service{
// 			Name:      db.Name,
// 			Namespace: db.Namespace,
// 		}
// 		serviceKey := service.Key()
// 		services[serviceKey] = service
// 	}

// 	return sets, services
// }

// func getClientRequests(clients map[string]crds.Client, statefulSets map[string]deployments.StatefulSet) map[string]database.Client {
// 	return nil
// }

// func getSecretRequests(clients map[string]crds.Client, existing map[string]database.Client) map[string]deployments.ClientSecret {
// 	return nil
// }
