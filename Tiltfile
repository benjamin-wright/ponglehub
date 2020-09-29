load('ext://namespace', 'namespace_create')
allow_k8s_contexts(['k3d-pongle'])
default_registry('localhost:5000')

namespace_create('ponglehub')

docker_build(
  'keycloak-init',
  'services/golang/keycloak-init/build/ponglehub.co.uk',
  dockerfile='services/golang/Dockerfile',
  build_args={
    'EXE_NAME': 'keycloak-init'
  }
)

k8s_yaml(helm(
  'deployment',
  name='ponglehub',
  namespace='ponglehub'
))