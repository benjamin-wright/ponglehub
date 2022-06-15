# Setup

allow_k8s_contexts(['pongle'])
load_dynamic('./libraries/tilt/infra.Tiltfile')
load('libraries/tilt/helm.Tiltfile', 'namespace_yaml')

k8s_yaml(namespace_yaml('apps'))

k8s_resource(
  'db',
  extra_pod_selectors=[{'db-operator.ponglehub.co.uk/deployment': 'db'}],
  port_forwards=["26257:26257"]
)

k8s_yaml(helm(
  'helm',
  name='database',
  namespace='apps',
  set=[
    'cockroach.db=256Mi',
  ]
))

# Landing Page

custom_build(
  'landing-page',
  'just ./static/landing-page/image $EXPECTED_REF',
  ['./static/landing-page/dist'],
  live_update=[
    sync('./static/landing-page/dist', '/usr/share/nginx/html')
  ]
)

k8s_yaml(helm(
  'helm',
  name='landing',
  namespace='apps',
  set=[
    'servers.landing-page.image=landing-page',
    'servers.landing-page.resources.limits.memory=64Mi',
    'servers.landing-page.resources.requests.memory=64Mi',
    'servers.landing-page.host=games.ponglehub.co.uk'
  ]
))

# Naughts and Crosses

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

k8s_yaml(helm(
  'helm',
  name='naughts-and-crosses',
  namespace='apps',
  set=[
    'servers.naughts-and-crosses-static.image=naughts-and-crosses',
    'servers.naughts-and-crosses-static.resources.limits.memory=64Mi',
    'servers.naughts-and-crosses-static.resources.requests.memory=64Mi',
    'servers.naughts-and-crosses-static.host=nac.ponglehub.co.uk',
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

# Draughts

custom_build(
  'draughts',
  'just ./static/draughts/image $EXPECTED_REF',
  ['./static/draughts/dist'],
  live_update=[
    sync('./static/draughts/dist', '/usr/share/nginx/html')
  ]
)

custom_build(
  'draughts-migrations',
  'just ./services/draughts/image-migrations $EXPECTED_REF',
  ['./services/draughts'],
  ignore=['./dist']
)

custom_build(
  'draughts-server',
  'just ./services/draughts/image $EXPECTED_REF',
  ['./services/draughts'],
  ignore=['./dist']
)

k8s_resource(
  'draughts-static',
  new_name='draughts: static'
)

k8s_resource(
  'draughts-server',
  new_name='draughts: server',
  resource_deps=[
    'draughts: migrations'
  ]
)

k8s_resource(
  'draughts-migrations',
  new_name='draughts: migrations',
  resource_deps=[
    'db'
  ]
)

k8s_yaml(helm(
  'helm',
  name='draughts',
  namespace='apps',
  set=[
    'servers.draughts-static.image=draughts',
    'servers.draughts-static.resources.limits.memory=64Mi',
    'servers.draughts-static.resources.requests.memory=64Mi',
    'servers.draughts-static.host=draughts.ponglehub.co.uk',
    'jobs.draughts-migrations.image=draughts-migrations',
    'jobs.draughts-migrations.db.cluster=db',
    'jobs.draughts-migrations.db.username=draughts_mig_user',
    'jobs.draughts-migrations.db.database=draughts',
    'jobs.draughts-migrations.resources.limits.memory=64Mi',
    'jobs.draughts-migrations.resources.requests.memory=64Mi',
    'servers.draughts-server.image=draughts-server',
    'servers.draughts-server.env.BROKER_URL="http://broker.auth-service.svc.cluster.local:80"',
    'servers.draughts-server.db.cluster=db',
    'servers.draughts-server.db.username=draughts_user',
    'servers.draughts-server.db.database=draughts',
    'servers.draughts-server.resources.limits.memory=64Mi',
    'servers.draughts-server.resources.requests.memory=64Mi',
    'servers.draughts-server.events={\'draughts.*\'}',
  ]
))
