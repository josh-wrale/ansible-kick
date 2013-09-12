# Ansible Kick

An alternative to [ansible-pull](http://www.ansibleworks.com/docs/playbooks2.html#pull-mode-playbooks) designed to make it easy for ec2 instances to request an Ansible push during initial boot; think auto-scaling.

[![Build Status](https://travis-ci.org/kelseyhightower/ansible-kick.png)](https://travis-ci.org/kelseyhightower/ansible-kick)

## How it works

 * ec2 instance request an Ansible push via a remote SSH forced command during initial boot

```
ssh -i /tmp/ansible_kick -o RequestTTY=no -o StrictHostKeyChecking=no ansiblekick@203.0.113.100
```

 * ansible-kick queries the ec2 API filtering by `private_ip_address` obtained from the `SSH_CLIENT` environment variable
 * ansible-kick searches the instance tags for a key named "role"
 * ansible-kick generates a temporary Ansible inventory file:

```TOML
[role]
ec2instance
```

 * ansible-kick locates an Ansible playbook matching the role name (role.yml)
 * ansible-kick calls `ansible-playbook` using the temporary inventory file and matching playbook

```
ansible-playbook -i /path/to/temp-inventory /path/to/role.yml
```

 * ansible-kick blocks until the `ansible-playbook` command exits
 * ec2 instance continues booting

## Install

    chmod +x /usr/local/bin/ansible-kick

## Setup and Configuration

[docs/setup.md](https://github.com/kelseyhightower/ansible-kick/blob/master/docs/setup.md)

## Build

ansible-kick requires the following dependencies:

 - launchpad.net/goamz/aws (bzr)
 - launchpad.net/goamz/ec2 (bzr)
 - github.com/BurntSushi/toml (git)

Clone this repository into `$GOPATH/src/github.com/kelseyhightower/ansible-kick`

    mkdir -p $GOPATH/src/github.com/kelseyhightower
    cd $GOPATH/src/github.com/kelseyhightower
    git clone https://github.com/kelseyhightower/ansible-kick.git
    cd ansible-kick
    
Then run:
    
    go get
    go build 

You should end-up with a working `ansible-kick` executable
