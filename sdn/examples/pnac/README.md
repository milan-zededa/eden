# SDN Example with 802.1x Port Authentication

TODO

Run the example with:

```shell
make clean && make build-tests
./eden config add default
./eden config set default --key sdn.disable --value false
./eden setup
./eden start --sdn-network-model $(pwd)/sdn/examples/pnac/network-model.json
./eden eve onboard
./eden controller edge-node set-config --file $(pwd)/sdn/examples/pnac/device-config.json
```
