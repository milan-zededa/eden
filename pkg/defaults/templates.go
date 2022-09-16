package defaults

//DefaultEdenTemplate is configuration template for Eden
const DefaultEdenTemplate = `#config is generated by eden
adam:
    #tag on adam container to pull
    tag: '{{parse "adam.tag"}}'

    #location of adam
    dist: '{{parse "adam.dist"}}'

    #port of adam
    port: {{parse "adam.port"}}

    #domain of adam
    domain: '{{parse "adam.domain"}}'

    #ip of adam for EVE access
    eve-ip: '{{parse "adam.eve-ip"}}'

    #ip of adam for EDEN access
    ip: '{{parse "adam.ip"}}'

    redis:
      #host of adam's redis for EDEN access
      eden: '{{parse "adam.redis.eden"}}'
      #host of adam's redis for ADAM access
      adam: '{{parse "adam.redis.adam"}}'

    #force adam rebuild
    force: {{parse "adam.force"}}

    #certificate for communication with adam
    ca: '{{parse "adam.ca"}}'

    #use remote adam
    remote:
        enabled: {{parse "adam.remote.enabled"}}

        #load logs and info from redis instead of http stream
        redis: {{parse "adam.remote.redis"}}

    #use v1 api
    v1: {{parse "adam.v1"}}

    caching:
        enabled: {{parse "adam.caching.enabled"}}

        #caching logs and info to redis instead of local
        redis: {{parse "adam.caching.redis"}}

        #prefix for directory/redis stream
        prefix: '{{parse "adam.caching.prefix"}}'

eve:
    #name
    name: '{{parse "eve.name"}}'

    #devmodel
    devmodel: '{{parse "eve.devmodel"}}'

    #devmodel file overwrite
    devmodelfile: '{{parse "eve.devmodelfile"}}'

    #Path to a file with JSON-formatted device config (EdgeDevConfig, but mostly just networking),
    #used to bootstrap device (i.e. establish connectivity with the controller and onboard).
    #The config will be reformatted to binary proto, signed and embedded into the EVE image.
    #Note: for legacy override.json use "eve-config-dir" arg of "eden setup".
    bootstrap-file: '{{parse "eve.bootstrap-file"}}'

    #Path to a file with JSON-formatted network config (modeled by DevicePortConfig
    #struct from EVE repo), applied in runtime using a specially formatted USB stick.
    #Also known as "usb.json".
    #Typically used for device bootstrapping as a second step after EVE installation.
    #This is a legacy method soon to be replaced by EdgeDevConfig-based bootstrapping.
    usbnetconf-file: '{{parse "eve.usbnetconf-file"}}'

    #EVE arch (amd64/arm64)
    arch: '{{parse "eve.arch"}}'

    #EVE os (linux/darwin)
    os: '{{parse "eve.os"}}'

    #EVE acceleration (set to false if you have problems with qemu)
    accel: {{parse "eve.accel"}}

    #variant of hypervisor of EVE (kvm/xen)
    hv: '{{parse "eve.hv"}}'

    #serial number in SMBIOS
    serial: '{{parse "eve.serial"}}'

    #onboarding certificate of EVE to put into adam
    cert: '{{parse "eve.cert"}}'

    #device certificate of EVE to put into adam
    device-cert: '{{parse "eve.device-cert"}}'

    #EVE pid file
    pid: '{{parse "eve.pid"}}'

    #EVE log file
    log: '{{parse "eve.log"}}'

    #EVE firmware
    firmware: {{parse "eve.firmware"}}

    #eve repo used in clone mode (eden.download = false)
    repo: '{{parse "eve.repo"}}'

    #eve registry to use
    registry: '{{parse "eve.registry"}}'

    #eve tag
    tag: '{{parse "eve.tag"}}'

    #port forwarding for EVE VM [(HOST:EVE)] when running without Eden-SDN
    hostfwd: '{{parse "eve.hostfwd"}}'

    #location of eve directory
    dist: '{{parse "eve.dist"}}'

    #file to save qemu config
    qemu-config: '{{parse "eve.qemu-config"}}'

    #uuid of EVE to use in cert
    uuid: '{{parse "eve.uuid"}}'

    #live image of EVE
    image-file: '{{parse "eve.image-file"}}'

    #dtb directory of EVE
    dtb-part: '{{parse "eve.dtb-part"}}'

    #config part of EVE
    config-part: '{{parse "eve.config-part"}}'

    #is EVE remote or local
    remote: {{parse "eve.remote"}}

    #EVE address for access from Eden
    remote-addr: '{{parse "eve.remote-addr"}}'

    #min level of logs saved in files on device
    log-level: '{{parse "eve.log-level"}}'

    #min level of logs sent to controller
    adam-log-level: '{{parse "eve.adam-log-level"}}'

    #port for telnet (console access)
    telnet-port: {{parse "eve.telnet-port"}}

    #ssid for wifi
    ssid: '{{parse "eve.ssid"}}'

    #cpu count
    cpu: {{parse "eve.cpu"}}

    #memory (MB)
    ram: {{parse "eve.ram"}}

    #disk (MB)
    disk: {{parse "eve.disk"}}

    #tpm
    tpm: {{parse "eve.tpm"}}

    #additional disks count
    disks: {{parse "eve.disks"}}

    #configuration specific to QEMU-emulated device
    qemu:
        #port for QEMU Monitor
        monitor-port: {{parse "eve.qemu.monitor-port"}}

        #base port for socket-based ethernet interfaces used in QEMU
        netdev-socket-port: {{parse "eve.qemu.netdev-socket-port"}}

eden:
    #root directory of eden
    root: '{{parse "eden.root"}}'
    #directory with tests
    tests: '{{parse "eden.tests"}}'
    images:
        #directory to save images
        dist: '{{parse "eden.images.dist"}}'

    #download eve instead of build
    download: {{parse "eden.download"}}

    #eserver is tool for serve images
    eserver:
        #ip (domain name) of eserver for EVE access
        eve-ip: '{{parse "eden.eserver.eve-ip"}}'

        #ip of eserver for EDEN access
        ip: '{{parse "eden.eserver.ip"}}'

        #port for eserver
        port: {{parse "eden.eserver.port"}}

        #tag of eserver container
        tag: '{{parse "eden.eserver.tag"}}'

        #force eserver rebuild
        force: {{parse "eden.eserver.force"}}

    #eclient is tool we use in tests
    eclient:
        #tag of eclient container
        tag: '{{parse "eden.eclient.tag"}}'
        #image of eclient container
        image: '{{parse "eden.eclient.image"}}'

    #directory to save certs
    certs-dist: '{{parse "eden.certs-dist"}}'

    #directory to save binaries
    bin-dist: '{{parse "eden.bin-dist"}}'

    #ssh-key to put into EVE
    ssh-key: '{{parse "eden.ssh-key"}}'

    #eden binary
    eden-bin: '{{parse "eden.eden-bin"}}'

    #test binary
    test-bin: '{{parse "eden.test-bin"}}'

    #test scenario
    test-scenario: '{{parse "eden.test-scenario"}}'

gcp:
    #path to the key to interact with gcp
    key: '{{parse "gcp.key"}}'

packet:
    #path to the key to interact with packet
    key: '{{parse "packet.key"}}'

redis:
    #port for access redis
    port: {{parse "redis.port"}}

    #tag for redis image
    tag: '{{parse "redis.tag"}}'

    #directory to use for redis persistence
    dist: '{{parse "redis.dist"}}'

registry:
    #port for registry access
    port: {{parse "registry.port"}}

    #tag for registry image
    tag: '{{parse "registry.tag"}}'

    #ip of registry for EDEN access
    ip: '{{parse "registry.ip"}}'

    # dist path to store registry data
    dist: '{{parse "registry.dist"}}'

sdn:
    #disable SDN
    disable: '{{parse "sdn.disable"}}'

    #directory with SDN source code
    source-dir: '{{parse "sdn.source-dir"}}'

    #directory where to put generated SDN-related config files
    config-dir: '{{parse "sdn.config-dir"}}'

    #live image of SDN
    image-file: '{{parse "sdn.image-file"}}'

    #path to linuxkit binary used to build SDN VM
    linuxkit-bin: '{{parse "sdn.linuxkit-bin"}}'

    #CPU count for SDN VM
    cpu: {{parse "sdn.cpu"}}

    #memory (MB) for SDN VM
    ram: {{parse "sdn.ram"}}

    #SDN pid file
    pid: '{{parse "sdn.pid"}}'

    #SDN file where console output is logged
    #Not as useful as logs from the SDN mgmt agent (get with: eden sdn logs)
    console-log: '{{parse "sdn.console-log"}}'

    #port for telnet (console access) to SDN VM
    telnet-port: {{parse "sdn.telnet-port"}}

    #port for SSH access to SDN VM
    ssh-port: {{parse "sdn.ssh-port"}}

    #port for access to the management agent running inside SDN VM
    mgmt-port: {{parse "sdn.mgmt-port"}}

    #path to JSON file with network model to apply into SDN
    #leave empty for default network model
    network-model: '{{parse "sdn.network-model"}}'
`

//DefaultQemuTemplate is configuration template for qemu
const DefaultQemuTemplate = `#qemu config file generated by eden
{{- if .Firmware }}
{{ $firmwareLength := len .Firmware }}{{ if eq $firmwareLength 1 }}
[machine]
  firmware = "{{ index .Firmware 0 }}"
{{- else if eq $firmwareLength 2 }}
[drive]
  if = "pflash"
  format = "raw"
  unit = "0"
  readonly = "on"
  file = "{{ index .Firmware 0 }}"

[drive]
  if = "pflash"
  format = "raw"
  unit = "1"
  file = "{{ index .Firmware 1 }}"
{{end}}
{{end}}
{{if .DTBDrive }}
[drive]
  file = "fat:rw:{{ .DTBDrive }}"
  format = "vvfat"
  label = "QEMU_DTB""
{{end}}
[rtc]
  base = "utc"
  clock = "rt"

[global]
  driver = "ICH9-LPC"
  property = "noreboot"
  value = "false"

[memory]
  size = "{{ .MemoryMB }}"

[smp-opts]
  cpus = "{{ .CPUs }}"

[device "usb"]
  driver = "qemu-xhci"

{{- if .USBTablets -}}
{{ range $i := .USBTablets }}
[device]
  driver = "usb-tablet"
{{ end }}
{{- end -}}

{{- if .USBSerials -}}
{{ range $i := .USBSerials }}
[chardev "charserial{{ $i }}"]
  backend = "pty"

[device "serial{{ $i }}"]
  driver = "usb-serial"
  chardev = "charserial{{ $i }}"
{{ end }}
{{- end -}}

{{ range .Disks }}
[drive]
  format = "qcow2"
  file = "{{.}}"
{{ end }}
`

// ParallelsDiskTemplate is template for disk annotation of parallels
const ParallelsDiskTemplate = `<?xml version='1.0' encoding='UTF-8'?>
<Parallels_disk_image Version="1.0">
    <Disk_Parameters>
        <Disk_size>{{ .DiskSize }}</Disk_size>
        <Cylinders>{{ .Cylinders }}</Cylinders>
        <PhysicalSectorSize>512</PhysicalSectorSize>
        <Heads>16</Heads>
        <Sectors>32</Sectors>
        <Padding>0</Padding>
        <Encryption>
            <Engine>{00000000-0000-0000-0000-000000000000}</Engine>
            <Data></Data>
        </Encryption>
        <UID>{{ .UID }}</UID>
        <Name>eve</Name>
        <Miscellaneous>
            <CompatLevel>level2</CompatLevel>
            <Bootable>1</Bootable>
            <SuspendState>0</SuspendState>
        </Miscellaneous>
    </Disk_Parameters>
    <StorageData>
        <Storage>
            <Start>0</Start>
            <End>{{ .DiskSize }}</End>
            <Blocksize>2048</Blocksize>
            <Image>
                <GUID>{{ .SnapshotUID }}</GUID>
                <Type>Compressed</Type>
                <File>live.0.{{ .SnapshotUID }}.hds</File>
            </Image>
        </Storage>
    </StorageData>
    <Snapshots>
        <Shot>
            <GUID>{{ .SnapshotUID }}</GUID>
            <ParentGUID>{00000000-0000-0000-0000-000000000000}</ParentGUID>
        </Shot>
    </Snapshots>
</Parallels_disk_image>`
