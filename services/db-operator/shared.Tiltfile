load('../../libraries/tilt/helm.Tiltfile', 'namespace_yaml')

k8s_yaml('crds/cockroach-client.crd.yaml')
k8s_yaml('crds/cockroach-db.crd.yaml')
k8s_kind('CockroachDB')

custom_build(
  'db-operator',
  'just image $EXPECTED_REF',
  ['./'],
  ignore=['Tiltfile', './dist']
)

k8s_resource(
  'db-operator',
  trigger_mode=TRIGGER_MODE_MANUAL
)

k8s_yaml(namespace_yaml('operators'))
k8s_yaml(helm(
  '../../helm',
  name='db-operator',
  namespace='operators',
  set=[
    'servers.db-operator.image=db-operator',
    'servers.db-operator.rbac.apiGroups={ponglehub.co.uk,apps,}',
    'servers.db-operator.rbac.resources={cockroachdbs,cockroachdbs/status,cockroachclients,cockroachclients/status,statefulsets,services,persistentvolumeclaims}',
    'servers.db-operator.rbac.verbs={get,list,watch,create,update,delete}',
    'servers.db-operator.rbac.clusterWide=true',
    'servers.db-operator.resources.limits.memory=64Mi',
    'servers.db-operator.resources.requests.memory=64Mi',
  ]
))