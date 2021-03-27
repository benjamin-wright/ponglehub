def load_files_encoded(basePath, fileDir):
  lines = []
  
  fileNames = str(local("ls %s" % fileDir)).rstrip('\n')
  for file in fileNames.split('\n'):
    data = str(local('cat %s/%s | base64' % (fileDir, file))).rstrip('\n')

    lines.append('%s."%s"=%s' % (basePath, file.replace('.', '\\.'), data))
  
  return blob(','.join(lines))