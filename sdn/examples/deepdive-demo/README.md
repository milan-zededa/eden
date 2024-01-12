```
cp -r go/src/github.com/lf-edge/eden/sdn/examples/deepdive-demo ~/dev/
```

Start with:
```shell
(cd ~/go/src/github.com/lf-edge/eden && make clean && make build-tests)
eden config add default
eden config set default --key sdn.disable --value false
eden config set default --key eve.tag --value demo
eden setup --zedcontrol=zedcloud.canary.zededa.net --soft-serial=milan0123456789
eden start --zedcontrol=zedcloud.canary.zededa.net --sdn-network-model $(pwd)/network-model.json
```

Create device, network instances and app instance.
```
eth0: mgmt, default IPv4 network
eth1: mgmt (by mistake), default IPv4 network

NI1: 10.10.1.0/24; 10.10.1.1; 10.10.1.2 - 10.10.1.254; NTP: 185.242.56.3
NI2: 10.10.2.0/24; 10.10.2.1; 10.10.2.2 - 10.10.2.254

INFO[0002] Please use 5d0767ee-0547-4569-b530-387e526f8cb9 as Onboarding Key 
INFO[0002] use milan0123456789 as Serial Number
```

Enable SSH access:
```
docker run -it --rm  -v "${HOME}/.ssh":/ssh --name  zcli zededa/zcli
zcli> zcli edge-node update milan-zedvirtual-demo --config="debug.enable.ssh:$(cat /ssh/id_rsa.pub)"
```

Backup:
```shell
make clean && make build-tests
./eden config add default
./eden config set default --key sdn.disable --value false
./eden config set default --key eve.tag --value demo
./eden setup 
./eden start --sdn-network-model $(pwd)/network-model.json
./eden eve onboard

###
./eden controller edge-node set-config --file $(pwd)/device-config.json
./eden pod deploy -n app --networks=netinst1 --networks=netinst2 --acl=netinst2:172.30.6.10 -p 2224:22 --memory=1GB ${IMG} \
  "http://10.10.10.102/ubuntu/ubuntu-22.04-server-cloudimg-amd64.img" \
  --metadata="#cloud-config\npassword: 123456\nchpasswd: { expire: False }\nssh_pwauth: True\nnetwork:\n  version: 2\n  ethernets:\n    enp3s0:\n      dhcp4: true"
```

/////////////////////////////////////////////////////////////////

- 4 tabs:
  - console access
  - 2 tabs for ssh access + fix DHCP / override
  - HTTP server
- have app definition, NIs and network configs open in zedcloud tabs


1. Console already opened, diag logs visible, explain them; eve verbose off; cat /run/diag.out;
2. Enter debug container; tools are there
3. collect-info.sh;  // maybe:  tar -xvf ...; show network folder
4. show netdump and the request attempts
5. ifconfig - IPs are OK
6. ip rule; ip route;
7. cat /etc/resolv.conf
8. dhcpcd -U -4 eth0; dhcpcd -U -4 eth1   - point out IPs, DNS servers, NTP server
9. ./fix-dhcp-server-config.sh
10. ip link set down eth0; dhcpcd -U -4 eth0; ip link set up eth0; dhcpcd -U -4 eth0
11. cat /etc/resolv.conf
12. Show diag.out, show device online in zedcloud
13. ssh into the device
14. cat /run/nim/DeviceNetworkStatus/global.json | jq
15. show eth1 error in zedui; fix eth1 usage
16. cat /persist/status/nim/DevicePortConfigList/global.json | jq
17. Wait for app to be running
18. cat /run/zedrouter/AppNetworkStatus/647470b5-3d4e-4a36-8819-a74bf354bb6b.json | jq
19. Point out Vif names; AllocatedIPv4Addr (empty for the second interface)
20. cat /run/zedrouter/dnsmasq.bn1.conf

21. eve attach to VM console
22. Try: curl http-server.demo/test
23. ip addr; point out missing IP; show wrong cloud-init config
24. Try: dhclient -v enp4s0
25. Check DNS: resolvectl status; this is OK
26. Try: nslookup http-server.demo 10.10.2.1; OK
27. Missing route: ip route add 172.16.0.0/12 dev enp4s0 via 10.10.2.1
28. Try: ping 172.30.5.10; fails
29. Go to another tab; ssh and trace packets: tcpdump -n -i any icmp 
30. check iptables: iptables -t mangle -L PREROUTING-nbu2x1-OUT -v --line-numbers -n
31. clear counters; retry: iptables -t mangle -Z PREROUTING-nbu2x1-OUT
32. Fix the iptables rule: iptables -t mangle -R PREROUTING-nbu2x1-OUT 6 -d 172.30.5.10/32 -m comment --comment "Manually edited rule" -j bn2-nbu2x1-2
33. Go back to the app; ping 172.30.5.10
34. Try: curl http-server.demo/test
35. Check NTP: systemctl status systemd-timesyncd


////////////////////////////////////////////////////////////


Get one-line `/run/global/DevicePortConfig/override.json`
```
cat override.json | jq -c 
echo '...' > /run/global/DevicePortConfig/override.json
```

Can be removed to unpublish.

Try httpserver:
```
curl http-server.demo/test
```

Try VM:
```
sudo dhclient -v enp4s0
resolvectl status
nslookup http-server.demo 10.10.2.1
ping http-server.demo

ip route add 172.16.0.0/12 dev enp4s0 via 10.10.2.1

iptables -t mangle -L PREROUTING-nbu2x1-OUT -v --line-numbers -n
iptables -t mangle -S PREROUTING-nbu2x1-OUT -v
iptables -t mangle -R PREROUTING-nbu2x1-OUT 6 -d 172.30.5.10/32 -m comment --comment "Manually edited rule" -j bn2-nbu2x1-2

curl http-server.demo/test
```