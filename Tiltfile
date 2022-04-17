# Setup

allow_k8s_contexts(['pongle'])
load('libraries/tilt/helm.Tiltfile', 'namespace_yaml')

k8s_kind('CockroachDB', api_version='ponglehub.co.uk/v1alpha1')

# Conditional resources

config.define_string_list("to-run", args=True)
cfg = config.parse()
groups = {
  'servers': ['db-operator', 'gateway', 'broker', 'redis'],
}

resources = []
for arg in cfg.get('to-run', []):
  if arg in groups:
    resources += groups[arg]
  else:
    resources.append(arg)

config.set_enabled_resources(resources)

# CRDs

k8s_yaml('services/event-gateway/crds/user.crd.yaml')
k8s_yaml('services/event-broker/crds/event-trigger.crd.yaml')

# Operators

load_dynamic('./services/db-operator/shared.Tiltfile')

# Auth & Comms services

custom_build(
  'event-gateway',
  'just ./services/event-gateway/image $EXPECTED_REF',
  ['./services/event-gateway'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'event-broker',
  'just ./services/event-broker/image $EXPECTED_REF',
  ['./services/event-broker'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'event-responder',
  'just ./services/event-responder/image $EXPECTED_REF',
  ['./services/event-responder'],
  ignore=['Tiltfile', './dist']
)

k8s_resource(
  'gateway',
  trigger_mode=TRIGGER_MODE_MANUAL,
  port_forwards=["4000:80"]
)

k8s_resource(
  'broker',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_resource(
  'responder',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_resource(
  'redis',
  port_forwards=["6379:6379"]
)

k8s_yaml(namespace_yaml('auth-service'))
k8s_yaml(helm(
  'helm',
  name='auth',
  namespace='auth-service',
  set=[
    'redis.redis.storage=256Mi',
    'secrets.gateway-key.keyfile=abcdefg',
    'servers.gateway.image=event-gateway',
    'servers.gateway.env.BROKER_URL="http://broker:80"',
    'servers.gateway.env.REDIS_URL="redis:6379"',
    'servers.gateway.env.KEY_FILE="/secrets/keyfile"',
    'servers.gateway.env.TOKEN_DOMAIN="ponglehub.co.uk"',
    'servers.gateway.env.ALLOWED_ORIGINS="http://ponglehub.co.uk\\,http://games.ponglehub.co.uk"',
    'servers.gateway.volFromSecret.gateway-key.path=/secrets',
    'servers.gateway.rbac.apiGroups={ponglehub.co.uk}',
    'servers.gateway.rbac.resources={authusers,authusers/status}',
    'servers.gateway.rbac.verbs={get,list,watch,patch,update}',
    'servers.gateway.rbac.clusterWide=true',
    'servers.gateway.resources.limits.memory=64Mi',
    'servers.gateway.resources.requests.memory=64Mi',
    'servers.gateway.host=ponglehub.co.uk',
    'servers.broker.image=event-broker',
    'servers.broker.rbac.apiGroups={ponglehub.co.uk}',
    'servers.broker.rbac.resources={eventtriggers}',
    'servers.broker.rbac.verbs={list,watch}',
    'servers.broker.rbac.clusterWide=true',
    'servers.broker.resources.limits.memory=32Mi',
    'servers.broker.resources.requests.memory=32Mi',
    'servers.responder.image=event-responder',
    'servers.responder.env.REDIS_URL="redis:6379"',
    'servers.responder.events={\'**.response\'}',
    'servers.responder.resources.limits.memory=32Mi',
    'servers.responder.resources.requests.memory=32Mi',
  ]
))

# Static files

custom_build(
  'landing-page',
  'just ./static/landing-page/image $EXPECTED_REF',
  ['./static/landing-page/dist'],
  live_update=[
    sync('./static/landing-page/dist', '/usr/share/nginx/html')
  ]
)

custom_build(
  'naughts-and-crosses',
  'just ./static/naughts-and-crosses/image $EXPECTED_REF',
  ['./static/naughts-and-crosses/dist'],
  live_update=[
    sync('./static/naughts-and-crosses/dist', '/usr/share/nginx/html')
  ]
)

custom_build(
  'naughts-and-crosses-migrations',
  'just ./services/naughts-and-crosses/image-migrations $EXPECTED_REF',
  ['./services/naughts-and-crosses'],
  ignore=['./dist']
)

custom_build(
  'naughts-and-crosses-server',
  'just ./services/naughts-and-crosses/image $EXPECTED_REF',
  ['./services/naughts-and-crosses'],
  ignore=['./dist']
)

custom_build(
  'draughts',
  'just ./static/draughts/image $EXPECTED_REF',
  ['./static/draughts/dist'],
  live_update=[
    sync('./static/draughts/dist', '/usr/share/nginx/html')
  ]
)

k8s_resource(
  'naughts-and-crosses-static',
  new_name='nac: static'
)

k8s_resource(
  'naughts-and-crosses-server',
  new_name='nac: server',
  resource_deps=[
    'nac: migrations'
  ]
)

k8s_resource(
  'naughts-and-crosses-migrations',
  new_name='nac: migrations',
  resource_deps=[
    'db'
  ]
)

k8s_resource(
  'db',
  extra_pod_selectors=[{'db-operator.ponglehub.co.uk/deployment': 'db'}],
  port_forwards=["26257:26257"]
)

k8s_yaml(namespace_yaml('apps'))

k8s_yaml(helm(
  'helm',
  name='landing',
  namespace='apps',
  set=[
    'servers.landing-page.image=landing-page',
    'servers.landing-page.resources.limits.memory=64Mi',
    'servers.landing-page.resources.requests.memory=64Mi',
    'servers.landing-page.host=games.ponglehub.co.uk',
  ]
))

k8s_yaml(helm(
  'helm',
  name='naughts-and-crosses',
  namespace='apps',
  set=[
    'servers.naughts-and-crosses-static.image=naughts-and-crosses',
    'servers.naughts-and-crosses-static.resources.limits.memory=64Mi',
    'servers.naughts-and-crosses-static.resources.requests.memory=64Mi',
    'servers.naughts-and-crosses-static.host=games.ponglehub.co.uk',
    'servers.naughts-and-crosses-static.path=/naughts-and-crosses',
    'servers.naughts-and-crosses-static.stripPrefix=true',
    'cockroach.db=256Mi',
    'jobs.naughts-and-crosses-migrations.image=naughts-and-crosses-migrations',
    'jobs.naughts-and-crosses-migrations.db.cluster=db',
    'jobs.naughts-and-crosses-migrations.db.username=nac_mig_user',
    'jobs.naughts-and-crosses-migrations.db.database=naughts_and_crosses',
    'jobs.naughts-and-crosses-migrations.resources.limits.memory=64Mi',
    'jobs.naughts-and-crosses-migrations.resources.requests.memory=64Mi',
    'servers.naughts-and-crosses-server.image=naughts-and-crosses-server',
    'servers.naughts-and-crosses-server.env.BROKER_URL="http://broker.auth-service.svc.cluster.local:80"',
    'servers.naughts-and-crosses-server.db.cluster=db',
    'servers.naughts-and-crosses-server.db.username=nac_user',
    'servers.naughts-and-crosses-server.db.database=naughts_and_crosses',
    'servers.naughts-and-crosses-server.resources.limits.memory=64Mi',
    'servers.naughts-and-crosses-server.resources.requests.memory=64Mi',
    'servers.naughts-and-crosses-server.events={\'naughts-and-crosses.*\'}',
  ]
))
k8s_yaml(helm(
  'helm',
  name='draughts',
  namespace='apps',
  set=[
    'servers.draughts.image=draughts',
    'servers.draughts.resources.limits.memory=64Mi',
    'servers.draughts.resources.requests.memory=64Mi',
    'servers.draughts.host=games.ponglehub.co.uk',
    'servers.draughts.path=/draughts',
    'servers.draughts.stripPrefix=true',
  ]
))

# Utils

local_resource(
  'add_users',
  '''
    ponglehub users add --resource-name test-user --display-name pingu --email test@user.com --password password
    ponglehub users add --resource-name other-user --display-name pongo --email other@user.com --password password
  ''',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL
)