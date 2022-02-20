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

k8s_yaml('services/db-operator/crds/cockroach-client.crd.yaml')
k8s_yaml('services/db-operator/crds/cockroach-db.crd.yaml')
k8s_yaml('services/event-gateway/crds/user.crd.yaml')
k8s_yaml('services/event-broker/crds/event-trigger.crd.yaml')

# Operators

custom_build(
  'db-operator',
  'just ./services/db-operator/image $EXPECTED_REF',
  ['./services/db-operator'],
  ignore=['Tiltfile', './dist']
)

k8s_resource(
  'db-operator',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_yaml(namespace_yaml('operators'))
k8s_yaml(helm(
  'helm',
  name='operator',
  namespace='operators',
  set=[
    'servers.db-operator.image=db-operator',
    'servers.db-operator.rbac.apiGroups={ponglehub.co.uk,apps,}',
    'servers.db-operator.rbac.resources={cockroachdbs,cockroachdbs/status,cockroachclients,cockroachclients/status,statefulsets,secrets,services,persistentvolumeclaims}',
    'servers.db-operator.rbac.verbs={get,list,watch,create,update,delete}',
    'servers.db-operator.rbac.clusterWide=true',
    'servers.db-operator.resources.limits.memory=64Mi',
    'servers.db-operator.resources.requests.memory=64Mi',
  ]
))

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
    'servers.gateway.env.BROKER_URL="http://gateway:80"',
    'servers.gateway.env.REDIS_URL="redis:6379"',
    'servers.gateway.env.KEY_FILE="/secrets/keyfile"',
    'servers.gateway.env.TOKEN_DOMAIN="ponglehub.co.uk"',
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
  ]
))

# Static files

custom_build(
  'landing-page',
  'just ./static/landing-page/image $EXPECTED_REF',
  ['./static/landing-page'],
  ignore=['./dist']
)

custom_build(
  'naughts-and-crosses',
  'just ./static/naughts-and-crosses/image $EXPECTED_REF',
  ['./static/naughts-and-crosses'],
  ignore=['./dist']
)

custom_build(
  'draughts',
  'just ./static/draughts/image $EXPECTED_REF',
  ['./static/draughts'],
  ignore=['./dist']
)

k8s_resource(
  'landing-page',
  trigger_mode=TRIGGER_MODE_MANUAL,
  port_forwards=["3000:80"]
)

k8s_resource(
  'naughts-and-crosses',
  trigger_mode=TRIGGER_MODE_MANUAL,
  port_forwards=["3001:80"]
)

k8s_resource(
  'draughts',
  trigger_mode=TRIGGER_MODE_MANUAL,
  port_forwards=["3002:80"]
)

k8s_yaml(namespace_yaml('static-files'))
k8s_yaml(helm(
  'helm',
  name='ui',
  namespace='static-files',
  set=[
    'servers.landing-page.image=landing-page',
    'servers.landing-page.resources.limits.memory=64Mi',
    'servers.landing-page.resources.requests.memory=64Mi',
    'servers.landing-page.host=games.ponglehub.co.uk',
    'servers.naughts-and-crosses.image=naughts-and-crosses',
    'servers.naughts-and-crosses.resources.limits.memory=64Mi',
    'servers.naughts-and-crosses.resources.requests.memory=64Mi',
    'servers.naughts-and-crosses.host=games.ponglehub.co.uk',
    'servers.naughts-and-crosses.path=/naughts-and-crosses',
    'servers.naughts-and-crosses.stripPrefix=true',
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
  'add_user',
  'ponglehub users add --resource-name test-user --display-name pingu --email test@user.com --password password',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL
)