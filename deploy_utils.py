# CONSTANTS
DEFAULT_TEMPLATE = 'k8s-conf.template.yaml'

ENV_DEVEL = 'devel'
ENV_PROD = 'prod'

ENV_TO_PROJ = {
    ENV_DEVEL: 'blorg-dev',
    ENV_PROD: 'blorg-prod'  # probably? idk!
}


def docker_tag(env, owner):
    return '%s-%s' % (env, owner)


def image_name(env, owner):
    """Generate the canonical name of the docker image for this server+env+user."""
    server = 'blorg-backend'
    gcloud_proj = ENV_TO_PROJ[env]
    tag = docker_tag(env, owner)

    return 'gcr.io/%(gcloud_proj)s/%(server)s:%(tag)s' % {
        'gcloud_proj': gcloud_proj,
        'server': server,
        'tag': tag,
     }


def tab_lines(s):
    lines = s.split('\n')
    lines[0] = '\t' + lines[0]
    return '\n\t'.join(lines)

