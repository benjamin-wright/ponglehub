allow_k8s_contexts(['k3d-pongle'])
default_registry('localhost:5000')

def namespace(name):
  return blob("""apiVersion: v1
kind: Namespace
metadata:
  name: %s
  annotations:
    linkerd.io/inject: enabled
""" % name)

k8s_yaml(namespace('ponglehub'))

def microservice(name):
  docker_build(
    name,
    'services/golang/%s/build/ponglehub.co.uk' % name,
    dockerfile='services/golang/Dockerfile',
    build_args={
      'EXE_NAME': name
    }
  )

microservice('wait-for-service')
microservice('keycloak-init')

k8s_yaml(helm(
  'deployment',
  name='ponglehub',
  namespace='ponglehub'
))