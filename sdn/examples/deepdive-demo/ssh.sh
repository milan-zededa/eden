#!/bin/sh

/home/mlenco/go/src/github.com/lf-edge/eden/eden sdn fwd eth0 22 -- ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -o PasswordAuthentication=no -i ~/.ssh/id_rsa root@FWD_IP -p FWD_PORT
