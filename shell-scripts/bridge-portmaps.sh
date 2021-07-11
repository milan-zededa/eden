#!/bin/bash
# Run after "./eden start"

bridge_net_prefix=192.168.99
eden_net=$(docker network inspect eden_network -f '{{(index .IPAM.Config 0).Subnet}}')

# cleanup from the previous run
iptables -t nat -D PREROUTING -j EDEN-PORTMAPS
iptables -t nat -D OUTPUT -j EDEN-PORTMAPS
iptables -t nat -F EDEN-PORTMAPS
iptables -t nat -X EDEN-PORTMAPS
iptables -t nat -D POSTROUTING -j EDEN-MASQUERADE
iptables -t nat -F EDEN-MASQUERADE
iptables -t nat -X EDEN-MASQUERADE
rm -rf /etc/xinetd.d/eve-portmap*

# Get IP addresses assigned to EVE by the DHCP server
function get_ip() {
  ifname=${1}
  mac=${2}
  until docker exec eve-dhcp-server cat /run/dnsmasq.leases | grep -q ${mac};
  do
    echo "Waiting for IP address to be assigned to ${ifname}..." 1>&2
    sleep 5
  done
  ip=$(docker exec eve-dhcp-server cat /run/dnsmasq.leases | grep ${mac} | awk '{print $3}')
  echo ${ip}
}
eth0_ip=$(get_ip eth0 02:00:00:12:34:56)
eth1_ip=$(get_ip eth1 02:00:00:12:34:57)

# portmaps
declare -a eth0_portmaps=(2222:22 2223:2223 2224:2224 5912:5902 5911:5901 8027:8027 8028:8028 8029:8029 8030:8030 8031:8031)
declare -a eth1_portmaps=(2233:2233 2234:2234)
iptables -t nat -N EDEN-PORTMAPS
iptables -t nat -I PREROUTING -j EDEN-PORTMAPS
iptables -t nat -I OUTPUT -j EDEN-PORTMAPS

function add_local_redirect() {
  fport=${1}
  lport=${2}
  eve_ip=${3}
  name="eve-portmap-${fport}"
  cat <<EOF > /etc/xinetd.d/${name}
service ${name}
{
 disable = no
 type = UNLISTED
 socket_type = stream
 protocol = tcp
 wait = no
 redirect = ${eve_ip} ${lport}
 bind = 127.0.0.1
 port = ${fport}
 user = nobody
}
EOF
}

for portmap in "${eth0_portmaps[@]}"
do
  fport=$(echo $portmap | cut -d ":" -f 1)
  lport=$(echo $portmap | cut -d ":" -f 2)
  iptables -t nat -A EDEN-PORTMAPS -p tcp -d ${eth0_ip} --dport ${fport} -j DNAT --to-destination ${eth0_ip}:${lport}
  iptables -t nat -A EDEN-PORTMAPS -p tcp -d ${eth1_ip} --dport ${fport} -j DNAT --to-destination ${eth0_ip}:${lport}
  add_local_redirect ${fport} ${lport} ${eth0_ip}
done

for portmap in "${eth1_portmaps[@]}"
do
  fport=$(echo $portmap | cut -d ":" -f 1)
  lport=$(echo $portmap | cut -d ":" -f 2)
  iptables -t nat -A EDEN-PORTMAPS -p tcp -d ${eth0_ip} --dport ${fport} -j DNAT --to-destination ${eth1_ip}:${lport}
  iptables -t nat -A EDEN-PORTMAPS -p tcp -d ${eth1_ip} --dport ${fport} -j DNAT --to-destination ${eth1_ip}:${lport}
  add_local_redirect ${fport} ${lport} ${eth1_ip}
done

sudo /etc/init.d/xinetd restart

# MASQUERADE bridge network
# But do not masquerade within the eden network, especially between Adam and EVE.
iptables -t nat -N EDEN-MASQUERADE
iptables -t nat -I POSTROUTING -j EDEN-MASQUERADE
iptables -t nat -A EDEN-MASQUERADE -s ${bridge_net_prefix}.0/24 ! -o "eve-bridge" ! -d "${eden_net}" -j MASQUERADE
