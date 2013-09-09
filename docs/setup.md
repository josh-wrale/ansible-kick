# Setup

## SSH key/pair

Generate a SSH key pair:

    ssh-keygen -b 1024 -C ansible-kick -f ansible_kick

You should now have the following files:

    ansible_kick
    ansible_kick.pub

## Ansible Kick

Create an ansiblekick user on the machine running ansible:

    sudo useradd -m ansiblekick

Add the `ansible_kick.pub` public key to the ansiblekick users `authorized_keys` file with restricted access.

    command="/usr/local/bin/ansible-kick",no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty ssh-rsa AAAAB3NzaC1y......xevPQ== ansible-kick

This is how we force only the `ansible-kick` command to be run and nothing else.

### Configuration

The default configuration file location is `/etc/ansible-kick/ansible-kick.conf` and can be set with the -c flag:

    /usr/local/bin/ansible-kick -c /etc/ansible-kick/ansible-kick.conf

The `ansible-kick.conf` file must contain the following settings:

    [aws]
    access_key_id = "AKIAIOSFODNN7EXAMPLE"
    secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    region = "us-east-1"

    [ansible]
    command = "/usr/bin/ansible-playbook"
    hosts_template = "/etc/ansible-kick/hosts.tmpl"
    playbook_path = "/etc/ansible/playbooks"

### Inventory Template

ansible-kick uses Golang templates when generating temporary inventory files

    [{{.Role}}]
    {{.Host}}

If the ip address of the ec2 instance is `203.0.113.10` and the role is nginx, the template will result in:

     [nginx]
     203.0.113.10

Templates support any of the Ansible [inventory patterns](http://www.ansibleworks.com/docs/patterns.html)

    [{{.Role}}]
    {{.Host}} ansible_connection=ssh ansible_ssh_user=ec2-user

## Client

Ansible kick was designed to be used during the initial ec2 instance boot process via a [userdata](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-waitcondition-article.html) script.

### Userdata

Create the following `userdata.sh` script: 

```bash
#!/bin/bash

KICK_USER="ansiblekick"
KICK_KEY="/tmp/ansible_kick"
ANSIBLE_HOST="203.0.113.100"

cat > ${KICK_KEY} <<-ENDOFKEY
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQC20SxTqn/YNEI3IcI3U77gwlbifevakeEMCt0AqOFQO/Gd2xCT
Example+Key+5rZnnyn6v+Ej6MBd/f9gQZB5OQU8mplxrAnv30qLK5nuIxIVgNJB
zPgA5CuKtwGBMYWFTAaDf2fL24EiA9xhhH8h/ZefqjhcUW6FAFWCOxevPQIDAQAB
AoGAG2vIgEwknONJw0c3AGF5UfEUYyiMBd63qLzAZWbvVL+JACpplCHBoE9MIxDX
GykgfXdMGDALPR8nfUirDZstco2/+Example+Key+4sERj83p2xCabfgb9lWhfuD
6dLSCkZxMlr+Y+SLwdmvp2oPM9QS0dXCMwPDXuyCwE9B7QECQQDrpxsbaJMRDzWf
IJIgwUSQCTO5DkYSUuvbxxe32tFYGP+tJlD61/HOdiqz5oyOcKqybGgasZsgKihn
Example+Key+EAxpouiQZVpitIYSrkGOnpZ0tn28lqk10DAtiaSoncokn/u60guw
GQgXTqvKL9r01qMEFc1CGtw5mRUpa4dxIQJAHpE1he+hrAPSC8sYyWDoeNqIuAdu
9W+GIqMHo5ShtRDBEX+332Hlfsd7MIzGTK+2pKBFPLkvCxQM263uQIWiTQJABKVE
o2XniPyIM+Wp8j8+e3ETG9wIljhNDvKhjg+w9gftdWSRCSrG0r8StH9mOlpX0dF8
XfowKqquGjuZfW9soQJAXI2ObA1vKje8rNH1cjXRp8uqwizZDs3/+Example+Key
dc/ylVFG6mP+l0NL4xYrVKk0Phouf73WgIgYzhfDZA==
-----END RSA PRIVATE KEY-----
ENDOFKEY

chmod 0600 ${KICK_KEY}
ssh -i ${KICK_KEY} -o RequestTTY=no -o StrictHostKeyChecking=no ${KICK_USER}@${ANSIBLE_HOST}
rm ${KICK_KEY}
```

Now base64 encode it:

    base64 userdata.sh
    ...

You should end up with a really long string which you can use as userdata.
