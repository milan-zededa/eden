# SDN Examples with multiple application interfaces and enforced order

Run the example with:

```shell
make clean && make build-tests
./eden stop
./eden config add default
./eden config set default --key eve.tag --value 0.0.0-app-interface-order-c1c26a94
./eden config set default --key sdn.disable --value false
./eden config set default --key eve.hv --value kvm
./eden setup
./eden start --sdn-network-model $(pwd)/sdn/examples/app-interface-order/network-model.json
./eden eve onboard

./eden controller edge-node set-config --file $(pwd)/sdn/examples/app-interface-order/device-config.json


IMG="http://10.10.10.102/images/ubuntu-22.04-server-cloudimg-amd64.img"
./eden pod deploy -n ubuntu --networks=ni-eth0:3 --networks=ni-eth1:2 --networks=ni-eth2:0 --adapters eth3:1 --adapters eth4:5 --adapters eth5:4 -p 2224:22 --memory=1GB ${IMG} --metadata="#cloud-config\nssh_pwauth: Yes\nhostname: myubuntu\nusers:\n  - name: pocuser\n    shell: /bin/bash\n    sudo: ALL=(ALL) NOPASSWD:ALL\nchpasswd:\n  list: |\n    pocuser:pocuser\n  expire: false\nruncmd:\n  - hostnamectl set-hostname myubuntu"

./eden pod deploy -n eclient docker://lfedge/eden-eclient:b96434e --networks=ni-eth0:3 --networks=ni-eth1:2 --networks=ni-eth2:0 --adapters eth3:1 --adapters eth4:5 --adapters eth5:4 -p 2224:22 --memory=1024MB
```

Once deployed, login to the application:

```shell
./eden eve ssh
CONSOLE="$(eve list-app-consoles | grep 89d15df7-0a37-49b3-8b90-ca7cf0c59cc4 | grep CONTAINER | awk '{print $4}')"
eve attach-app-console "$CONSOLE"
```