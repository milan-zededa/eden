package utils

import (
	"bytes"
	"fmt"
	"github.com/lf-edge/eden/pkg/defaults"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

//ConfigVars struct with parameters from config file
type ConfigVars struct {
	AdamIP            string
	AdamPort          string
	AdamDomain        string
	AdamDir           string
	AdamCA            string
	AdamRemote        bool
	AdamCaching       bool
	AdamCachingRedis  bool
	AdamCachingPrefix string
	AdamRemoteRedis   bool
	AdamRedisUrlEden  string
	AdamRedisUrlAdam  string
	EveHV             string
	EveSSID           string
	EveName           string
	EveRemote         bool
	EveRemoteAddr     string
	SshKey            string
	CheckLogs         bool
	EveCert           string
	EveSerial         string
	ZArch             string
	DevModel          string
	EdenBinDir        string
	EdenProg          string
	TestProg          string
	TestScenario      string
	EServerImageDist  string
	EServerPort       string
	EServerIp         string
}

//InitVars loads vars from viper
func InitVars() (*ConfigVars, error) {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}
	loaded, err := LoadConfigFile(configPath)
	if err != nil {
		return nil, err
	}
	if loaded {
		var vars = &ConfigVars{
			AdamIP:            viper.GetString("adam.ip"),
			AdamPort:          viper.GetString("adam.port"),
			AdamDomain:        viper.GetString("adam.domain"),
			AdamDir:           ResolveAbsPath(viper.GetString("adam.dist")),
			AdamCA:            ResolveAbsPath(viper.GetString("adam.ca")),
			AdamRedisUrlEden:  viper.GetString("adam.redis.eden"),
			AdamRedisUrlAdam:  viper.GetString("adam.redis.adam"),
			SshKey:            ResolveAbsPath(viper.GetString("eden.ssh-key")),
			CheckLogs:         viper.GetBool("eden.logs"),
			EveCert:           ResolveAbsPath(viper.GetString("eve.cert")),
			EveSerial:         viper.GetString("eve.serial"),
			ZArch:             viper.GetString("eve.arch"),
			EveSSID:           viper.GetString("eve.ssid"),
			EveHV:             viper.GetString("eve.hv"),
			DevModel:          viper.GetString("eve.devmodel"),
			EveName:           viper.GetString("eve.name"),
			EveRemote:         viper.GetBool("eve.remote"),
			EveRemoteAddr:     viper.GetString("eve.remote-addr"),
			AdamRemote:        viper.GetBool("adam.remote.enabled"),
			AdamRemoteRedis:   viper.GetBool("adam.remote.redis"),
			AdamCaching:       viper.GetBool("adam.caching.enabled"),
			AdamCachingPrefix: viper.GetString("adam.caching.prefix"),
			AdamCachingRedis:  viper.GetBool("adam.caching.redis"),
			EdenBinDir:        viper.GetString("eden.bin-dist"),
			EdenProg:          viper.GetString("eden.eden-bin"),
			TestProg:          viper.GetString("eden.test-bin"),
			TestScenario:      viper.GetString("eden.test-scenario"),
			EServerImageDist:  ResolveAbsPath(viper.GetString("eden.images.dist")),
			EServerPort:       viper.GetString("eden.eserver.port"),
			EServerIp:         viper.GetString("eden.eserver.ip"),
		}
		return vars, nil
	}
	return nil, nil
}

var defaultEnvConfig = `#config is generated by eden
adam:
    #tag on adam container to pull
    tag: {{ .DefaultAdamTag }}

    #location of adam
    dist: {{ .DefaultAdamDist }}

    #port of adam
    port: {{ .DefaultAdamPort }}

    #domain of adam
    domain: {{ .DefaultDomain }}

    #ip of adam for EVE access
    eve-ip: {{ .IP }}

    #ip of adam for EDEN access
    ip: {{ .IP }}

    redis:
      #host of adam's redis for EDEN access
      eden: redis://{{ .IP }}:{{ .DefaultRedisPort }}
      #host of adam's redis for ADAM access
      adam: redis://{{ .DefaultRedisContainerName }}:{{ .DefaultRedisPort }}

    #force adam rebuild
    force: true

    #certificate for communication with adam
    ca: {{ .DefaultAdamDist }}/run/config/root-certificate.pem

    #use remote adam
    remote: 
        enabled: true

        #load logs and info from redis instead of http stream
        redis: true

    #use v1 api
    v1: true

    caching:
        enabled: false

        #caching logs and info to redis instead of local
        redis: false
        
        #prefix for directory/redis stream
        prefix: cache

eve:
    #name
    name: {{ .DefaultEVEName }}

    #devmodel
    devmodel: {{ .DefaultEVEModel }}

    #EVE arch (amd64/arm64)
    arch: {{ .Arch }}

    #EVE os (linux/darwin)
    os: {{ .OS }}

    #EVE acceleration (set to false if you have problems with qemu)
    accel: true

    #variant of hypervisor of EVE (kvm/xen)
    hv: {{ .DefaultEVEHV }}

    #serial number in SMBIOS
    serial: "{{ .DefaultEVESerial }}"

    #onboarding certificate of EVE to put into adam
    cert: certs/onboard.cert.pem

    #EVE pid file
    pid: eve.pid

    #EVE log file
    log: eve.log

    #EVE firmware
    firmware: [{{ .DefaultImageDist }}/eve/OVMF_CODE.fd,{{ .DefaultImageDist }}/eve/OVMF_VARS.fd]

    #eve repo used in clone mode (eden.download = false)
    repo: {{ .DefaultEveRepo }}

    #eve tag
    tag: {{ .DefaultEVETag }}

    #forward of ports in qemu [(HOST:EVE)]
    hostfwd:
        {{ .DefaultSSHPort }}: 22
        5912: 5901
        5911: 5900
        8027: 8027
        8028: 8028

    #location of eve directory
    dist: {{ .DefaultEVEDist }}

    #file to save qemu config
    qemu-config: {{ .EdenDir }}/qemu.conf

    #uuid of EVE to use in cert
    uuid: {{ .UUID }}

    #live image of EVE
    image-file: {{ .DefaultImageDist }}/eve/live.img

    #dtb directory of EVE
    dtb-part: ""

    #config part of EVE
    config-part: {{ .DefaultAdamDist }}/run/config

    #is EVE remote or local
    remote: {{ .DefaultEVERemote }}

    #EVE address for access from Eden
    remote-addr: {{ .DefaultEVERemoteAddr }}

eden:
    #root directory of eden
    root: {{ .Root }}
    images:
        #directory to save images
        dist: {{ .DefaultImageDist }}

    #download eve instead of build
    download: true

    #eserver is tool for serve images
    eserver:
        #ip (domain name) of eserver for EVE access
        eve-ip: {{ .DefaultDomain }}

        #ip of eserver for EDEN access
        ip: {{ .IP }}

        #port for eserver
        port: {{ .DefaultEserverPort }}

        #tag of eserver container
        tag: {{ .DefaultEServerTag }}

        #force eserver rebuild
        force: true

    #directory to save certs
    certs-dist: {{ .DefaultCertsDist }}

    #directory to save binaries
    bin-dist: {{ .DefaultBinDist }}

    #ssh-key to put into EVE
    ssh-key: {{ .DefaultSSHKey }}

    #observe logs in tests
    logs: false

    #eden binary
    eden-bin: eden

    #test binary
    test-bin: "{{ .DefaultTestProg }}"

    #test scenario
    test-scenario: "{{ .DefaultTestScenario }}"

redis:
    #port for access redis
    port: {{ .DefaultRedisPort }}

    #tag for redis image
    tag: {{ .DefaultRedisTag }}

    #directory to use for redis persistence
    dist: {{ .DefaultRedisDist }}
`

//DefaultEdenDir returns path to default directory
func DefaultEdenDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, defaults.DefaultEdenHomeDir), nil
}

//DefaultConfigPath returns path to default config
func DefaultConfigPath() (string, error) {
	context, err := ContextLoad()
	if err != nil {
		return "", fmt.Errorf("context load error: %s", err)
	}
	return context.GetCurrentConfig(), nil
}

//CurrentDirConfigPath returns path to eden-config.yml in current folder
func CurrentDirConfigPath() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(currentPath, defaults.DefaultCurrentDirConfig), nil
}

//LoadConfigFile load config from file with viper
func LoadConfigFile(config string) (loaded bool, err error) {
	if config == "" {
		config, err = DefaultConfigPath()
		if err != nil {
			return false, fmt.Errorf("fail in DefaultConfigPath: %s", err.Error())
		}
	} else {
		context, err := ContextInit()
		if err != nil {
			return false, fmt.Errorf("context Load DefaultEdenDir error: %s", err)
		}
		contextFile := context.GetCurrentConfig()
		if config != contextFile {
			loaded, err := LoadConfigFile(contextFile)
			if err != nil {
				return loaded, err
			}
		}
	}
	log.Debugf("Will use config from %s", config)
	if _, err = os.Stat(config); os.IsNotExist(err) {
		log.Fatal("no config, please run 'eden config add default'")
	}
	abs, err := filepath.Abs(config)
	if err != nil {
		return false, fmt.Errorf("fail in reading filepath: %s", err.Error())
	}
	viper.SetConfigFile(abs)
	if err := viper.MergeInConfig(); err != nil {
		return false, fmt.Errorf("failed to read config file: %s", err.Error())
	}
	currentFolderDir, err := CurrentDirConfigPath()
	if err != nil {
		log.Errorf("CurrentDirConfigPath: %s", err)
	} else {
		log.Debugf("Try to add config from %s", currentFolderDir)
		if _, err = os.Stat(currentFolderDir); !os.IsNotExist(err) {
			abs, err = filepath.Abs(currentFolderDir)
			if err != nil {
				log.Errorf("CurrentDirConfigPath absolute: %s", err)
			} else {
				viper.SetConfigFile(abs)
				if err := viper.MergeInConfig(); err != nil {
					log.Errorf("failed in merge config file: %s", err.Error())
				} else {
					log.Debugf("Merged config with %s", abs)
				}
			}
		}
	}
	return true, nil
}

//GenerateConfigFile is a function to generate default yml
func GenerateConfigFile(filePath string) error {
	return generateConfigFileFromTemplate(filePath, defaultEnvConfig)
}

func generateConfigFileFromTemplate(filePath string, templateString string) error {
	context, err := ContextInit()
	if err != nil {
		return err
	}
	context.Save()
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	edenDir, err := DefaultEdenDir()
	if err != nil {
		log.Fatal(err)
	}

	t := template.New("t")
	_, err = t.Parse(templateString)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	ip, err := GetIPForDockerAccess()
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	nets, err := GetSubnetsNotUsed(1)
	if err != nil {
		return err
	}
	address := strings.Split(nets[0].FirstAddress.String(), ".")
	eveIP := strings.Join(append(strings.Split(nets[0].FirstAddress.String(), ".")[:len(address)-1], "2"), ".")
	err = t.Execute(buf,
		struct {
			DefaultAdamDist      string
			DefaultAdamTag       string
			DefaultAdamPort      int
			DefaultImageDist     string
			ImageDir             string
			Root                 string
			IP                   string
			EVEIP                string
			UUID                 string
			Arch                 string
			OS                   string
			EdenDir              string
			DefaultBaseOSVersion string
			DefaultBaseOSTag     string
			DefaultEVETag        string
			DefaultDomain        string
			DefaultRedisPort     int
			DefaultRedisTag      string
			DefaultEVEDist       string
			DefaultEserverPort   int
			DefaultEVESerial     string
			DefaultRedisDist     string
			DefaultCertsDist     string
			DefaultBinDist       string
			DefaultEVEHV         string
			DefaultSSHPort       int
			DefaultTestScenario  string
			DefaultTestProg      string
			DefaultSSHKey        string
			DefaultEveRepo       string

			DefaultEVEModel string
			DefaultEVEName  string

			DefaultEVERemote     bool
			DefaultEVERemoteAddr string

			DefaultRedisContainerName string

			DefaultEServerTag string
		}{
			DefaultAdamDist:     defaults.DefaultAdamDist,
			DefaultAdamPort:     defaults.DefaultAdamPort,
			DefaultAdamTag:      defaults.DefaultAdamTag,
			DefaultImageDist:    defaults.DefaultImageDist,
			ImageDir:            filepath.Join(currentPath, defaults.DefaultImageDist),
			Root:                filepath.Join(currentPath, defaults.DefaultDist),
			IP:                  ip,
			EVEIP:               eveIP,
			UUID:                id.String(),
			Arch:                runtime.GOARCH,
			OS:                  runtime.GOOS,
			EdenDir:             edenDir,
			DefaultEVETag:       defaults.DefaultEVETag,
			DefaultDomain:       defaults.DefaultDomain,
			DefaultRedisPort:    defaults.DefaultRedisPort,
			DefaultRedisTag:     defaults.DefaultRedisTag,
			DefaultEVEDist:      defaults.DefaultEVEDist,
			DefaultEserverPort:  defaults.DefaultEserverPort,
			DefaultEVESerial:    defaults.DefaultEVESerial,
			DefaultRedisDist:    defaults.DefaultRedisDist,
			DefaultCertsDist:    defaults.DefaultCertsDist,
			DefaultBinDist:      defaults.DefaultBinDist,
			DefaultEVEHV:        defaults.DefaultEVEHV,
			DefaultSSHPort:      defaults.DefaultSSHPort,
			DefaultTestScenario: defaults.DefaultTestScenario,
			DefaultTestProg:     defaults.DefaultTestProg,
			DefaultSSHKey:       defaults.DefaultSSHKey,
			DefaultEveRepo:      defaults.DefaultEveRepo,

			DefaultEVEModel: defaults.DefaultEVEModel,
			DefaultEVEName:  defaults.DefaultEVEName,

			DefaultEVERemote:     defaults.DefaultEVERemote,
			DefaultEVERemoteAddr: defaults.DefaultEVEHost,

			DefaultRedisContainerName: defaults.DefaultRedisContainerName,

			DefaultEServerTag: defaults.DefaultEServerTag,
		})
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
