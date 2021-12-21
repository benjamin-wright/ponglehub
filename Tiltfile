allow_k8s_contexts(['pongle'])

load('libraries/tilt/helm.Tiltfile', 'namespace_yaml')

k8s_kind('CockroachDB', api_version='ponglehub.co.uk/v1alpha1')

default_registry('localhost:5000', host_from_cluster='pongle-registry.localhost:5000')

custom_build(
  'auth-operator',
  'mudly ./services/auth-operator+image && docker tag auth-operator $EXPECTED_REF',
  ['./services/auth-operator'],
  ignore=['Tiltfile', './dist']
)

k8s_resource(
  'auth-operator',
  trigger_mode=TRIGGER_MODE_MANUAL
)

custom_build(
  'auth-server',
  'mudly ./services/auth-server+server && docker tag auth-server $EXPECTED_REF',
  ['./services/auth-server'],
  ignore=['Tiltfile', './dist']
)


k8s_resource(
  'auth-server',
  extra_pod_selectors=[{'serving.knative.dev/configuration': 'auth-server'}],
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_yaml('services/auth-operator/crds/user.crd.yaml')
k8s_yaml(namespace_yaml('auth-service'))

k8s_yaml(helm(
  'helm',
  name='auth',
  namespace='auth-service',
  set=[
    'db.enabled=true',
    'servers.auth-operator.enabled=true',
    'servers.auth-operator.env.BROKER_URL="http://auth-service:80"',
    'servers.auth-operator.image=auth-operator',
    'servers.auth-operator.rbac.apiGroups={ponglehub.co.uk}',
    'servers.auth-operator.rbac.resources={authusers,authusers/status}',
    'servers.auth-operator.rbac.verbs={get,list,watch,patch,update}',
    'servers.auth-operator.rbac.clusterWide=true',
    'servers.auth-operator.resources.limits.memory=64Mi',
    'servers.auth-operator.resources.requests.memory=64Mi',
    'servers.auth-server.enabled=true',
    'servers.auth-server.db.user=auth-user',
    'servers.auth-server.db.database=auth-db',
    'servers.auth-server.image=auth-server'
  ]
))