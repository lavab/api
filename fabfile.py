import os, random, string
from fabric.api import run, env, cd

def deploy():
	branch = os.getenv('DRONE_BRANCH', 'master')
	commit = os.getenv('DRONE_COMMIT', 'master')
	tmp_dir = '/tmp/' + ''.join(random.choice(string.lowercase) for i in xrange(10))

	with cd(tmp_dir):
		run('git clone git@github.com:lavab/api.git')
		with cd('api'):
			run('docker build -t registry.lavaboom.io/lavaboom/api-' + branch + ' .')

		run('git clone git@github.com:lavab/docker.git')
		with cd('docker/runners'):
			run('docker rm -f api-' + branch)
			run('./api-' + branch + '.sh')

	run('rm -r ' + tmp_dir)