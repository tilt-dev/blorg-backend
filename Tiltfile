def blorg_backend_local():
  return blorg_backend('devel')

def blorg_backend_prod():
  return blorg_backend('prod')

def blorg_backend(env):
  entrypoint = '/app/server'
  image = build_docker_image('Dockerfile.base', 'gcr.io/blorg-dev/blorg-backend:' + env + '-' + local('whoami').rstrip('\n'), entrypoint)
  src_dir = '/go/src/github.com/windmilleng/blorg-backend'
  image.add_mount(src_dir, git_repo('.'))
  image.add_cmd('cd ' + src_dir + '; go get ./...')
  image.add_cmd('mkdir -p /app')
  image.add_cmd('cd ' + src_dir + '; go build -o server; cp server /app/')
  # print(image)
  yaml = local('python populate_config_template.py ' + env + ' 1>&2 && cat k8s-conf.generated.yaml')

  # this api might be cleaner than stderr stuff above
  # run('python populate_config_template.py ' + env')
  # yaml = read('k8s-conf.generated.yaml')

  # print(yaml)
  return k8s_service(yaml, image)
