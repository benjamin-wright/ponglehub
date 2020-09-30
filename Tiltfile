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

def envvar(name):
  return str(local("echo $%s" % name)).rstrip('\n')

k8s_yaml(helm(
  'deployment',
  name='ponglehub',
  namespace='ponglehub',
  set=[
    'global.keycloakUser='+envvar('KEYCLOAK_USER'),
    'global.keycloakPassword='+envvar('KEYCLOAK_PASSWORD'),
    'keycloak.postgresql.postgresqlPassword='+envvar('KEYCLOAK_DB_PASSWORD')
  ]
))