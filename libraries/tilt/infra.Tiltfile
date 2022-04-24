# Setup

load('./helm.Tiltfile', 'namespace_yaml')

k8s_kind('CockroachDB', api_version='ponglehub.co.uk/v1alpha1')

# CRDs

k8s_yaml('../../services/event-gateway/crds/user.crd.yaml')
k8s_yaml('../../services/event-broker/crds/event-trigger.crd.yaml')

# Operators

load_dynamic('../../services/db-operator/shared.Tiltfile')

# Auth & Comms services

custom_build(
  'event-gateway',
  'just ../../services/event-gateway/image $EXPECTED_REF',
  ['../../services/event-gateway'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'event-broker',
  'just ../../services/event-broker/image $EXPECTED_REF',
  ['../../services/event-broker'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'event-responder',
  'just ../../services/event-responder/image $EXPECTED_REF',
  ['../../services/event-responder'],
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
  '../../helm',
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