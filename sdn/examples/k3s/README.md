Under EVE repo:
```shell
# checkout https://github.com/zedi-pramodh/eve/tree/clustered-eve-poc

make pkg/* && make eve HV=kubevirt
...
Tagging docker.io/lfedge/eve:v7.4.0-master-private-1124-gcacab6707-kubevirt-amd64 as docker.io/lfedge/eve:v7.4.0-master-private-1124-gcacab6707-kubevirt
Build complete, not pushing, all done.
rm images/rootfs-kubevirt.yml.in

docker tag lfedge/eve:v7.4.0-master-private-1124-gcacab6707-kubevirt-amd64 lfedge/eve:test-kubevirt-amd64
```

Under Eden repo:
```shell
# Single network port:
make clean && make build-tests
./eden config add default
./eden config set default --key eve.hv --value kubevirt
./eden config set default --key eve.tag --value test
./eden config set default --key eve.log-level --value debug
./eden config set default --key sdn.disable --value false
./eden config set default --key=eve.disks --value=2
./eden config set default --key=eve.disk --value=8192
./eden setup -v debug --grub-options='set_global dom0_extra_args "$dom0_extra_args eve_install_clustered_storage_sizeGB=2 "'
./eden start 
./eden eve onboard
```

```shell
# Two network ports:
make clean && make build-tests
./eden config add default
./eden config set default --key eve.hv --value kubevirt
./eden config set default --key eve.tag --value test
./eden config set default --key eve.log-level --value debug
./eden config set default --key sdn.disable --value false
./eden config set default --key=eve.disks --value=2
./eden config set default --key=eve.disk --value=8192
./eden setup -v debug --eve-bootstrap-file $(pwd)/sdn/examples/k3s/device-config.json --grub-options='set_global dom0_extra_args "$dom0_extra_args eve_install_clustered_storage_sizeGB=2 "'
./eden start --sdn-network-model $(pwd)/sdn/examples/k3s/network-model.json 
./eden eve onboard
./eden controller edge-node set-config --file $(pwd)/sdn/examples/k3s/device-config.json 
```