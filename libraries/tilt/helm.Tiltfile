def load_files_encoded(basePath, fileDir):
  lines = []
  
  fileNames = str(local("ls %s" % fileDir)).rstrip('\n')
  for file in fileNames.split('\n'):
    data = str(local('cat %s/%s | base64' % (fileDir, file))).rstrip('\n')

    lines.append('%s."%s"=%s' % (basePath, file.replace('.', '\\.'), data))
  
  return blob(','.join(lines))

def file(name):
  return str(local("cat %s | base64" % name)).rstrip('\n')

def namespace_yaml(name):
  """Returns YAML for a namespace
  Args:
    name: The namespace name. Currently not validated.
  Returns:
    The namespace YAML as a blob
  """

  return blob("""apiVersion: v1
kind: Namespace
metadata:
  name: %s
  labels:
    istio-enabled: true
""" % name)

def envvar(name):
  return str(local("echo $%s" % name)).rstrip('\n')
