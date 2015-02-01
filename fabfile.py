import os, random, string
from fabric.api import run, env, cd, settings

env.key_filename = os.getenv('HOME', '/root') + '/.ssh/id_rsa'

def deploy():
	branch = os.getenv('DRONE_BRANCH', 'master')
	commit = os.getenv('DRONE_COMMIT', 'master')
	tmp_dir = '/tmp/' + ''.join(random.choice(string.lowercase) for i in xrange(10))

	run('mkdir ' + tmp_dir)
	with cd(tmp_dir):
		run('git clone git@github.com:lavab/api.git')
		with cd('api'):
			run('git checkout ' + commit)
			run('docker build -t registry.lavaboom.io/lavaboom/api-' + branch + ' .')

		run('git clone git@github.com:lavab/docker.git')
		with settings(warn_only=True):
			run('docker rm -f api-' + branch)
		with cd('docker/runners'):
			run('./api-' + branch + '.sh')

	run('rm -r ' + tmp_dir)