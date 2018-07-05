#!/usr/bin/env python

import deploy_utils as utils
from populate_config_template import populate_config_template

import argparse
import getpass
import subprocess
import textwrap


def parse_args():
    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description=textwrap.dedent("""
        Deploy the `blorg backend` app to Kubernetes development cluster (i.e. NOT PROD).

        Deploy consists of the following steps:
        1. (re)build Docker image
        2. push docker image to gcr.io
        3. generate k8s config file (by populating template)
        4. create/update k8s from config file

        If this is your first time deploying to the dev cluster, run with flag --first-deploy.
        """))
    parser.add_argument('--config_template', '-c', type=str,
                        help=('path to config template file (default: %s)' %
                              utils.DEFAULT_TEMPLATE),
                        default=utils.DEFAULT_TEMPLATE)
    parser.add_argument('--first_deploy', action='store_true',
                        help=('use this flag if this is your first blorg deploy ('
                              'will set the appropriate cluster, etc.)'))
    return parser.parse_args()


def main():
    # setup
    args = parse_args()
    owner = getpass.getuser()
    imgname = utils.image_name(utils.ENV_DEVEL, owner)

    if args.first_deploy:
        proj = utils.ENV_TO_PROJ[utils.ENV_DEVEL]
        print('+ Configuring kubectl for approriate cluster...')
        out = subprocess.check_output(
            ['gcloud', 'container', 'clusters', 'get-credentials', 'blorg', '--zone',
             'us-central1-b', '--project', proj])
        print('~~~ Configured kubectl for cluster "%s" with output:\n%s' %
              (proj, utils.tab_lines(out.decode("utf-8"))))

        # NOTE(maia): I'm not actually sure if this will work b/c I think this command
        # requires user input... well, we'll see!
        print('+ Authorizing Docker push for gcloud repo...')
        out = subprocess.check_output(
            ['gcloud', 'auth', 'configure-docker'])
        print('~~~ Docker authorized to gcloud repo with output:\n%s' %
              utils.tab_lines(out.decode("utf-8")))

    # 1. (re)build Docker image
    print('+ (Re)building Docker image...')
    out = subprocess.check_output(['docker', 'build', '-t', imgname, '.'])
    print('~~~ Built Docker image "%s" with output:\n%s' %
          (imgname, utils.tab_lines(out.decode("utf-8"))))

    # 2. push docker image to gcr.io
    print('+ Pushing Docker image...')
    out = subprocess.check_output(['docker', 'push', imgname])
    print('~~~ Pushed Docker image with output:\n%s' % utils.tab_lines(out.decode("utf-8")))

    # 3. generate k8s config file (by populating template)
    print('+ Generating k8s file from template "%s"...' % args.config_template)
    config = populate_config_template(args.config_template, utils.ENV_DEVEL, owner)
    print('~~~ Generated config file: "%s"\n' % config)

    # 4. create/update k8s from config file
    print('+ Deleting existing pods for this app+owner+env...')
    # TODO: template these keys/vals in conf file so everything is controlled by Python
    labels = {
        'app': 'blorg',
        'environment': utils.ENV_DEVEL,
        'owner': owner,
        'tier': 'backend'
    }
    selectors = []
    for selector in ['%s=%s' % (k, v) for k, v in labels.iteritems()]:
        selectors.append(selector)
    selectors_string = ",".join(selectors)
    cmd = ['kubectl', 'delete', 'pods', '-l', selectors_string]
    out = subprocess.check_output(cmd)
    print('~~~ Deleted existing pods (if any) with output:\n%s' %
          utils.tab_lines(out.decode("utf-8")))

    print('+ Applying generated k8s config...')
    out = subprocess.check_output(['kubectl', 'apply', '-f', config])
    print('~~~ Successfully applied config with output:\n%s' %
          utils.tab_lines(out.decode("utf-8")))


if __name__ == '__main__':
    main()
