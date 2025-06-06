# Edge-Sync test

https://github.com/zededa/edge-sync/blob/main/TESTING.md

```shell
make clean && make build-tests
./eden config add default
./eden config set default --key sdn.disable --value false
./eden config set default --key eve.tag --value 0.0.0-airgap-mode-20add8b8
./eden setup --zedcontrol=zedcloud.alpha.zededa.net --soft-serial=milan-edge-sync-test
./eden start --sdn-network-model $(pwd)/sdn/examples/edge-sync/network-model.json \
             --zedcontrol=zedcloud.alpha.zededa.net
```

```shell
http://edge-sync.sdn:1180
```

```shell
zcli edge-node update milan-edge-sync-test --config="debug.enable.console:true"
zcli edge-node update milan-edge-sync-test --config="debug.enable.vga:true"
zcli edge-node update milan-edge-sync-test --config="debug.enable.usb:true"
zcli edge-node update milan-edge-sync-test --config="airgap.mode:enabled"
```

```shell
curl -F file=@milan-edge-sync-test_config.json http://localhost:1180/api/v1/user/config
```

```shell
git apply ./sdn/examples/edge-sync/deny-cloud-access.diff
./eden sdn net-model apply ./sdn/examples/edge-sync/network-model.json

./eden sdn ssh
killall goproxy
```

```shell
2025-06-20T16:54:15.457692066Z;pillar.out;{"file":"/pillar/conntester/cloud.go:215","func":"github.com/lf-edge/eve/pkg/pillar/conntester.(*CloudConnectivityTester).testLOCConnectivity","level":"info","msg":"HEY! testLOCConnectivity, rv: {CloudReachable:false RemoteTempFailure:false IntfStatusMap:{StatusMap:map[eth0:{LastFailed:2025-06-20 16:54:15.456330732 +0000 UTC m=+448.144845819 LastSucceeded:0001-01-01 00:00:00 +0000 UTC LastError:SendOnIntf to http://edge-sync.sdn:1180/api/v2/edgedevice/ping reqlen 0 statuscode 404 Not Found body:\n00000000  34 30 34 20 70 61 67 65  20 6e 6f 74 20 66 6f 75  |404 page not fou|\n00000010  6e 64                                             |nd| LastWarning:}]} TracedReqs:[]}, err: All attempts to connect to http://edge-sync.sdn:1180/api/v2/edgedevice/ping failed: send via eth0: SendOnIntf to http://edge-sync.sdn:1180/api/v2/edgedevice/ping reqlen 0 statuscode 404 Not Found body:\n00000000  34 30 34 20 70 61 67 65  20 6e 6f 74 20 66 6f 75  |404 page not fou|\n00000010  6e 64                                             |nd|","pid":2430,"source":"nim","time":"2025-06-20T16:54:15.456577229Z"}
```