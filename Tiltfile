load('ext://namespace', 'namespace_create')
allow_k8s_contexts(['k3d-pongle'])

namespace_create('ponglehub')

docker_build(
    'keycloak-init',
    'services/golang/keycloak/build/ponglehub.co.uk',
    dockerfile='services/golang/Dockerfile',
    build_args={
        'EXE_NAME': 'keycloak-init'
    },
    ignore=[
        'keycloak-init-go-tmp-umask'
    ]
)

helm(
  'deployment',
  name='ponglehub',
  namespace='ponglehub'
)