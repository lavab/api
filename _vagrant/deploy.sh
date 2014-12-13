#!/usr/bin/env bash

export DEBIAN_FRONTEND=noninteractive

apt-get update -qq
apt-get install -y python2.7 python-pip

pip install ansible

mkdir /etc/ansible
echo "localhost" > /etc/ansible/hosts

cd /vagrant
ansible-playbook -vvvv playbook.yml
