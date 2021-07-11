#!/bin/bash
# Run before "./eden start"

qemu_user=${SUDO_USER}
echo "Creating bridge with TAPs for the user ${qemu_user}"
bridge_net_prefix=192.168.99

# cleanup of the previous run
docker kill eve-dhcp-server
docker rm eve-dhcp-server
ip link del eve-tap0
ip link del eve-tap1
ip link del eve-bridge
ip link del eve-dhcp-veth1
ip link del eve-dhcp-veth2
ip netns delete eve-dhcp-server
iptables -D DOCKER-USER -i eve-bridge -j ACCEPT
iptables -D DOCKER-USER -o eve-bridge -j ACCEPT

ip link add name eve-bridge type bridge
brctl setfd eve-bridge 0
ip link set eve-bridge up
ip addr add ${bridge_net_prefix}.1/24 dev eve-bridge

tunctl -t eve-tap0 -u ${qemu_user}
ip link set eve-tap0 up
ip link set dev eve-tap0 master eve-bridge

tunctl -t eve-tap1 -u ${qemu_user}
ip link set eve-tap1 up
ip link set dev eve-tap1 master eve-bridge

ip link add eve-dhcp-veth1 type veth peer name eve-dhcp-veth2
ip link set dev eve-dhcp-veth2 master eve-bridge
ip link set eve-dhcp-veth2 up

CONFIG_DIR=`mktemp -d`
cat <<EOF > ${CONFIG_DIR}/dnsmasq.conf
log-queries
log-dhcp
bind-interfaces
except-interface=lo
dhcp-leasefile=/run/dnsmasq.leases
interface=eve-dhcp-veth1
dhcp-range=${bridge_net_prefix}.3,${bridge_net_prefix}.254,60m
dhcp-option=option:router,${bridge_net_prefix}.1
EOF

docker run --cap-add=NET_ADMIN --rm -d -v ${CONFIG_DIR}/dnsmasq.conf:/etc/dnsmasq.conf --name eve-dhcp-server --entrypoint=""\
	strm/dnsmasq /bin/sh -c "sleep 5 && dnsmasq -d"

pid=$(docker inspect -f '{{.State.Pid}}' eve-dhcp-server)
ip netns attach eve-dhcp-server ${pid}
ip link set eve-dhcp-veth1 netns eve-dhcp-server
ip netns exec eve-dhcp-server ip link set eve-dhcp-veth1 up
ip netns exec eve-dhcp-server ip addr add ${bridge_net_prefix}.2/24 dev eve-dhcp-veth1

iptables -N DOCKER-USER
iptables -I DOCKER-USER -i eve-bridge -j ACCEPT
iptables -I DOCKER-USER -o eve-bridge -j ACCEPT