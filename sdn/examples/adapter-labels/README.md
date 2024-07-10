# SDN Example with Static IP Configuration

```shell
make clean && make build-tests
./eden config add default
./eden config set default --key sdn.disable --value false
./eden setup --eve-bootstrap-file $(pwd)/sdn/examples/adapter-labels/device-config.json
./eden start --sdn-network-model $(pwd)/sdn/examples/adapter-labels/network-model.json 
./eden eve onboard
./eden controller edge-node set-config --file $(pwd)/sdn/examples/adapter-labels/device-config.json 
```

```shell
make clean && make build-tests
./eden config add default
./eden config set default --key sdn.disable --value false
./eden setup --eve-config-dir $(pwd)/sdn/examples/adapter-labels/config-overrides
./eden start --sdn-network-model $(pwd)/sdn/examples/adapter-labels/network-model.json 
./eden eve onboard
./eden controller edge-node set-config --file $(pwd)/sdn/examples/adapter-labels/device-config.json 
```

```shell
make clean && make build-tests
./eden config add default
./eden config set default --key eve.log-level --value debug
./eden config set default --key eve.accel --value true
./eden config set default --key=eve.tpm --value=false
./eden setup
EDEN_TEST_SETUP=y ./eden test ./tests/workflow -s networking.tests.txt -v debug
```

```shell
./eden pod deploy -v debug -n eclient docker://lfedge/eden-eclient:b96434e -p 2223:22 --networks=ni0 --memory=512MB

./eden sdn fwd eth0 2223 -- ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -o PasswordAuthentication=no -i /home/mlenco/go/src/github.com/lf-edge/eden/dist/tests/eclient/image/cert/id_rsa root@FWD_IP -p FWD_PORT
```