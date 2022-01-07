allow_k8s_contexts(['pongle'])
load('libraries/tilt/helm.Tiltfile', 'namespace_yaml')

k8s_kind('CockroachDB', api_version='ponglehub.co.uk/v1alpha1')

custom_build(
  'db-operator',
  'mudly ./services/db-operator+image && docker tag localhost:5000/db-operator $EXPECTED_REF',
  ['./services/db-operator'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'event-gateway',
  'mudly ./services/event-gateway+image && docker tag localhost:5000/event-gateway $EXPECTED_REF',
  ['./services/event-gateway'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'event-broker',
  'mudly ./services/event-broker+image && docker tag localhost:5000/event-broker $EXPECTED_REF',
  ['./services/event-broker'],
  ignore=['Tiltfile', './dist']
)

k8s_resource(
  'db-operator',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_resource(
  'gateway',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_resource(
  'broker',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_yaml('services/db-operator/crds/cockroach-client.crd.yaml')
k8s_yaml('services/db-operator/crds/cockroach-db.crd.yaml')
k8s_yaml('services/auth-operator/crds/user.crd.yaml')
k8s_yaml('services/event-broker/crds/event-trigger.crd.yaml')

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
    'servers.gateway.env.TOKEN_DOMAIN="localhost"',
    'servers.gateway.volFromSecret.gateway-key.path=/secrets',
    'servers.gateway.rbac.apiGroups={ponglehub.co.uk}',
    'servers.gateway.rbac.resources={authusers,authusers/status}',
    'servers.gateway.rbac.verbs={get,list,watch,patch,update}',
    'servers.gateway.rbac.clusterWide=true',
    'servers.gateway.resources.limits.memory=64Mi',
    'servers.gateway.resources.requests.memory=64Mi',
    'servers.broker.image=event-broker',
    'servers.broker.rbac.apiGroups={ponglehub.co.uk}',
    'servers.broker.rbac.resources={eventtriggers}',
    'servers.broker.rbac.verbs={list,watch}',
    'servers.broker.rbac.clusterWide=true',
    'servers.broker.resources.limits.memory=32Mi',
    'servers.broker.resources.requests.memory=32Mi',
  ]
))