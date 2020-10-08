allow_k8s_contexts(['k3d-pongle'])
default_registry('localhost:5000')

load('ext://restart_process', 'docker_build_with_restart')

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
    dockerfile='services/dockerfiles/golang.Dockerfile',
    build_args={
      'EXE_NAME': name
    }
  )

def vue(name):
  docker_build(
    name,
    'services/node/%s/dist' % name,
    dockerfile='services/dockerfiles/static.Dockerfile',
    live_update=[
      sync('services/node/%s/dist' % name, '/usr/share/nginx/html')
    ],
  )

microservice('wait-for-service')
microservice('keycloak-init')
vue('landing-page')

def envvar(name):
  return str(local("echo $%s" % name)).rstrip('\n')

def file(name):
  return str(local("cat %s | base64" % name)).rstrip('\n')

k8s_yaml(helm(
  'deployment',
  name='ponglehub',
  namespace='ponglehub',
  set=[
    'global.keycloak.user='+envvar('KEYCLOAK_USER'),
    'global.keycloak.password='+envvar('KEYCLOAK_PASSWORD'),
    'global.smtp.email='+envvar('KEYCLOAK_EMAIL'),
    'global.smtp.password='+envvar('KEYCLOAK_EMAIL_PASSWORD'),
    'global.smtp.host='+envvar('KEYCLOAK_SMTP_SERVER'),
    'global.smtp.port='+envvar('KEYCLOAK_SMTP_PORT'),
    'global.smtp.from='+envvar('KEYCLOAK_SMTP_FROM'),
    'global.ssl.key='+file('infra/cluster/ssl/ponglehub.co.uk.key'),
    'global.ssl.crt='+file('infra/cluster/ssl/ponglehub.co.uk.crt'),
    'oauthProxy.clientSecret=1790d1f2-18f7-4fbb-bed2-de1287924e72',
    'oauthProxy.cookieSecret=some-other-secret'
  ]
))