allow_k8s_contexts(['pongle'])
load('libraries/tilt/helm.Tiltfile', 'namespace_yaml')

k8s_kind('CockroachDB', api_version='ponglehub.co.uk/v1alpha1')
k8s_kind('Service', api_version='serving.knative.dev/v1', image_json_path='{.spec.template.spec.containers[*].image}')

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
  'auth-server',
  'mudly ./services/auth-server+server && docker tag localhost:5000/auth-server $EXPECTED_REF',
  ['./services/auth-server'],
  ignore=['Tiltfile', './dist']
)

custom_build(
  'auth-server-events',
  'mudly ./services/auth-server+events && docker tag localhost:5000/auth-server-events $EXPECTED_REF',
  ['./services/auth-server'],
  ignore=['Tiltfile', './dist']
)

k8s_resource(
  'auth-operator',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_resource(
  'auth-server',
  extra_pod_selectors=[{'serving.knative.dev/configuration': 'auth-server'}],
  trigger_mode=TRIGGER_MODE_MANUAL,
  resource_deps=['cockroach']
)

k8s_yaml('services/db-operator/crds/cockroach-client.crd.yaml')
k8s_yaml('services/db-operator/crds/cockroach-db.crd.yaml')
k8s_yaml('services/auth-operator/crds/user.crd.yaml')

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
    'eventBrokers.events=true',
    'servers.auth-operator.env.BROKER_URL="http://broker-ingress.knative-eventing.svc.cluster.local/auth-service/events"',
    'servers.auth-operator.image=auth-operator',
    'servers.auth-operator.rbac.apiGroups={ponglehub.co.uk}',
    'servers.auth-operator.rbac.resources={authusers,authusers/status}',
    'servers.auth-operator.rbac.verbs={get,list,watch,patch,update}',
    'servers.auth-operator.rbac.clusterWide=true',
    'servers.auth-operator.resources.limits.memory=64Mi',
    'servers.auth-operator.resources.requests.memory=64Mi',
    'servers.auth-server.db.cluster=cockroach',
    'servers.auth-server.db.username=auth_user',
    'servers.auth-server.db.database=auth_db',
    'servers.auth-server.image=auth-server',
    'functions.auth-server-events.image=auth-server-events',
    'functions.auth-server-events.env.BROKER_URL="http://broker-ingress.knative-eventing.svc.cluster.local/auth-service/events"',
    'functions.auth-server-events.db.cluster=cockroach',
    'functions.auth-server-events.db.username=auth_user',
    'functions.auth-server-events.db.database=auth_db',
    'functions.auth-server-events.eventBroker=events'
  ]
))