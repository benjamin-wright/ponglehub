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
  'auth-operator',
  'mudly ./services/auth-operator+image && docker tag localhost:5000/auth-operator $EXPECTED_REF',
  ['./services/auth-operator'],
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
  'auth-operator',
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

k8s_resource(
  'cockroach',
  extra_pod_selectors={
    'db-operator.ponglehub.co.uk/deployment': 'cockroach'
  }
)

k8s_yaml(namespace_yaml('auth-service'))
k8s_yaml(helm(
  'helm',
  name='auth',
  namespace='auth-service',
  set=[
    'cockroach.cockroach=1Gi',
    'servers.auth-operator.env.BROKER_URL="http://recorder:80"',
    'servers.auth-operator.image=auth-operator',
    'servers.auth-operator.rbac.apiGroups={ponglehub.co.uk}',
    'servers.auth-operator.rbac.resources={authusers,authusers/status}',
    'servers.auth-operator.rbac.verbs={get,list,watch,patch,update}',
    'servers.auth-operator.rbac.clusterWide=true',
    'servers.auth-operator.resources.limits.memory=64Mi',
    'servers.auth-operator.resources.requests.memory=64Mi',
    'servers.broker.image=event-broker',
    'servers.broker.rbac.apiGroups={ponglehub.co.uk}',
    'servers.broker.rbac.resources={eventtriggers}',
    'servers.broker.rbac.verbs={list,watch}',
    'servers.broker.rbac.clusterWide=true',
    'servers.broker.resources.limits.memory=32Mi',
    'servers.broker.resources.requests.memory=32Mi',
  ]
))