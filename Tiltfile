allow_k8s_contexts(['k3d-pongle'])
default_registry('localhost:5000')

load('ext://restart_process', 'docker_build_with_restart')

def namespace(name):
  return blob("""apiVersion: v1
kind: Namespace
metadata:
  name: %s
  # annotations:
  #   linkerd.io/inject: enabled
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

def rust(name):
  docker_build_with_restart(
    name,
    'services/rust/%s/build' % name,
    dockerfile='services/dockerfiles/rust.Dockerfile',
    build_args={
      'EXE_NAME': name
    },
    entrypoint='/rust_binary',
    live_update=[
      sync('services/rust/%s/build/%s' % (name, name), '/rust_binary')
    ]
  )

def migration(name):
  docker_build(
    '%s-migrations' % name,
    'migrations/%s' % name,
    dockerfile='migrations/flyway.Dockerfile'
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

# microservice('wait-for-service')
# microservice('keycloak-init')
migration('auth')
rust('db-init')
rust('auth-server')
rust('auth-controller')
rust('gatekeeper')
vue('landing-page')
microservice('game-state')

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
    'global.ssl.key='+file('infra/terraform/infra/.scratch/ingress.key'),
    'global.ssl.crt='+file('infra/terraform/infra/.scratch/ingress.crt'),
    'seed.name='+envvar('SEED_USER_NAME'),
    'seed.email='+envvar('SEED_USER_EMAIL'),
  ]
))

k8s_resource(
  'migrations',
  trigger_mode=TRIGGER_MODE_MANUAL,
  auto_init=True
)