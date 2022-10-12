```
make clean && make build-tests
./eden config add default
./eden config set default --key eve.tag --value test
./eden config set default --key sdn.disable --value false
./eden setup
./eden start --sdn-network-model $(pwd)/sdn/examples/custom-dns-config/network-model.json
./eden eve onboard
```