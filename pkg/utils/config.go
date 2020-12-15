package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/lf-edge/eden/pkg/defaults"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	AdamRedisURLEden  string
	AdamRedisURLAdam  string
	EveHV             string
	EveSSID           string
	EveUUID           string
	EveName           string
	EveRemote         bool
	EveRemoteAddr     string
	EveQemuPorts      map[string]string
	EveQemuConfig     string
	EveDist           string
	SSHKey            string
	EveCert           string
	EveDeviceCert     string
	EveSerial         string
	ZArch             string
	DevModel          string
	EdenBinDir        string
	EdenProg          string
	TestProg          string
	TestScenario      string
	EServerImageDist  string
	EServerPort       string
	EServerIP         string
	LogLevel          string
	AdamLogLevel      string
}

//InitVars loads vars from viper
func InitVars() (*ConfigVars, error) {
	loaded := true
	if viper.ConfigFileUsed() == "" {
		configPath, err := DefaultConfigPath()
		if err != nil {
			return nil, err
		}
		loaded, err = LoadConfigFile(configPath)
		if err != nil {
			return nil, err
		}
	}
	if loaded {
		edenHome, err := DefaultEdenDir()
		if err != nil {
			log.Fatal(err)
		}
		globalCertsDir := filepath.Join(edenHome, defaults.DefaultCertsDist)
		if _, err := os.Stat(globalCertsDir); os.IsNotExist(err) {
			if err = os.MkdirAll(globalCertsDir, 0755); err != nil {
				log.Fatal(err)
			}
		}
		caCertPath := filepath.Join(globalCertsDir, "root-certificate.pem")
		var vars = &ConfigVars{
			AdamIP:            viper.GetString("adam.ip"),
			AdamPort:          viper.GetString("adam.port"),
			AdamDomain:        viper.GetString("adam.domain"),
			AdamDir:           ResolveAbsPath(viper.GetString("adam.dist")),
			AdamCA:            caCertPath,
			AdamRedisURLEden:  viper.GetString("adam.redis.eden"),
			AdamRedisURLAdam:  viper.GetString("adam.redis.adam"),
			SSHKey:            ResolveAbsPath(viper.GetString("eden.ssh-key")),
			EveCert:           ResolveAbsPath(viper.GetString("eve.cert")),
			EveDeviceCert:     ResolveAbsPath(viper.GetString("eve.device-cert")),
			EveSerial:         viper.GetString("eve.serial"),
			EveDist:           viper.GetString("eve.dist"),
			EveQemuConfig:     viper.GetString("eve.qemu-config"),
			ZArch:             viper.GetString("eve.arch"),
			EveSSID:           viper.GetString("eve.ssid"),
			EveHV:             viper.GetString("eve.hv"),
			DevModel:          viper.GetString("eve.devmodel"),
			EveName:           viper.GetString("eve.name"),
			EveUUID:           viper.GetString("eve.uuid"),
			EveRemote:         viper.GetBool("eve.remote"),
			EveRemoteAddr:     viper.GetString("eve.remote-addr"),
			EveQemuPorts:      viper.GetStringMapString("eve.hostfwd"),
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
			EServerIP:         viper.GetString("eden.eserver.ip"),
			LogLevel:          viper.GetString("eve.log-level"),
			AdamLogLevel:      viper.GetString("eve.adam-log-level"),
		}
		return vars, nil
	}
	return nil, nil
}

var defaultEnvConfig = `#config is generated by eden
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

    #forward of ports in qemu [(HOST:EVE)]
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

eden:
    #root directory of eden
    root: '{{parse "eden.root"}}'
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
`

//DefaultEdenDir returns path to default directory
func DefaultEdenDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, defaults.DefaultEdenHomeDir), nil
}

//GetConfig return path to config file
func GetConfig(name string) string {
	edenDir, err := DefaultEdenDir()
	if err != nil {
		log.Fatalf("GetCurrentConfig DefaultEdenDir error: %s", err)
	}
	return filepath.Join(edenDir, defaults.DefaultContextDirectory, fmt.Sprintf("%s.yml", name))
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

func loadConfigFile(config string, local bool) (loaded bool, err error) {
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
	if local {
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
	}
	return true, nil
}

//LoadConfigFile load config from file with viper
func LoadConfigFile(config string) (loaded bool, err error) {
	return loadConfigFile(config, true)
}

//LoadConfigFileContext load config from context file with viper
func LoadConfigFileContext(config string) (loaded bool, err error) {
	return loadConfigFile(config, false)
}

//GenerateConfigFile is a function to generate default yml
func GenerateConfigFile(filePath string) error {
	context, err := ContextInit()
	if err != nil {
		return err
	}
	context.Save()
	return generateConfigFileFromTemplate(filePath, defaultEnvConfig, context)
}

func generateConfigFileFromTemplate(filePath string, templateString string, context *Context) error {
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

	ip, err := GetIPForDockerAccess()
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	imageDist := fmt.Sprintf("%s-%s", context.Current, defaults.DefaultImageDist)

	certsDist := fmt.Sprintf("%s-%s", context.Current, defaults.DefaultCertsDist)

	parse := func(inp string) interface{} {
		switch inp {
		case "adam.tag":
			return defaults.DefaultAdamTag
		case "adam.dist":
			return defaults.DefaultAdamDist
		case "adam.port":
			return defaults.DefaultAdamPort
		case "adam.domain":
			return defaults.DefaultDomain
		case "adam.eve-ip":
			return ip
		case "adam.ip":
			return ip
		case "adam.redis.eden":
			return fmt.Sprintf("redis://%s:%d", ip, defaults.DefaultRedisPort)
		case "adam.redis.adam":
			return fmt.Sprintf("redis://%s:%d", defaults.DefaultRedisContainerName, defaults.DefaultRedisPort)
		case "adam.force":
			return true
		case "adam.ca":
			return filepath.Join(certsDist, "root-certificate.pem")
		case "adam.remote.enabled":
			return true
		case "adam.remote.redis":
			return true
		case "adam.v1":
			return true
		case "adam.caching.enabled":
			return false
		case "adam.caching.redis":
			return false
		case "adam.caching.prefix":
			return "cache"

		case "eve.name":
			return strings.ToLower(context.Current)
		case "eve.devmodel":
			return defaults.DefaultQemuModel
		case "eve.arch":
			return runtime.GOARCH
		case "eve.os":
			return runtime.GOOS
		case "eve.accel":
			return true
		case "eve.hv":
			return defaults.DefaultEVEHV
		case "eve.serial":
			return defaults.DefaultEVESerial
		case "eve.cert":
			return filepath.Join(certsDist, "onboard.cert.pem")
		case "eve.device-cert":
			return filepath.Join(certsDist, "device.cert.pem")
		case "eve.pid":
			return fmt.Sprintf("%s-eve.pid", strings.ToLower(context.Current))
		case "eve.log":
			return fmt.Sprintf("%s-eve.log", strings.ToLower(context.Current))
		case "eve.firmware":
			return fmt.Sprintf("[%s %s]",
				filepath.Join(imageDist, "eve", "OVMF_CODE.fd"),
				filepath.Join(imageDist, "eve", "OVMF_VARS.fd"))
		case "eve.repo":
			return defaults.DefaultEveRepo
		case "eve.registry":
			return defaults.DefaultEveRegistry
		case "eve.tag":
			return defaults.DefaultEVETag
		case "eve.hostfwd":
			return fmt.Sprintf("{\"%d\":\"22\",\"5912\":\"5902\",\"5911\":\"5901\",\"8027\":\"8027\",\"8028\":\"8028\"}", defaults.DefaultSSHPort)
		case "eve.dist":
			return fmt.Sprintf("%s-%s", context.Current, defaults.DefaultEVEDist)
		case "eve.qemu-config":
			return filepath.Join(edenDir, fmt.Sprintf("%s-%s", context.Current, defaults.DefaultQemuFileToSave))
		case "eve.uuid":
			return id.String()
		case "eve.image-file":
			return filepath.Join(imageDist, "eve", "live.img")
		case "eve.dtb-part":
			return ""
		case "eve.config-part":
			return certsDist
		case "eve.remote":
			return defaults.DefaultEVERemote
		case "eve.remote-addr":
			return defaults.DefaultEVEHost
		case "eve.log-level":
			return defaults.DefaultEveLogLevel
		case "eve.adam-log-level":
			return defaults.DefaultAdamLogLevel
		case "eve.telnet-port":
			return defaults.DefaultTelnetPort
		case "eve.ssid":
			return ""

		case "eden.root":
			return filepath.Join(currentPath, defaults.DefaultDist)
		case "eden.images.dist":
			return defaults.DefaultEserverDist
		case "eden.download":
			return true
		case "eden.eserver.eve-ip":
			return defaults.DefaultDomain
		case "eden.eserver.ip":
			return ip
		case "eden.eserver.port":
			return defaults.DefaultEserverPort
		case "eden.eserver.tag":
			return defaults.DefaultEServerTag
		case "eden.eserver.force":
			return true
		case "eden.certs-dist":
			return certsDist
		case "eden.bin-dist":
			return defaults.DefaultBinDist
		case "eden.ssh-key":
			return fmt.Sprintf("%s-%s", context.Current, defaults.DefaultSSHKey)
		case "eden.eden-bin":
			return "eden"
		case "eden.test-bin":
			return defaults.DefaultTestProg
		case "eden.test-scenario":
			return defaults.DefaultTestScenario

		case "gcp.key":
			return ""

		case "redis.port":
			return defaults.DefaultRedisPort
		case "redis.tag":
			return defaults.DefaultRedisTag
		case "redis.dist":
			return defaults.DefaultRedisDist

		case "registry.port":
			return defaults.DefaultRegistryPort
		case "registry.tag":
			return defaults.DefaultRegistryTag
		case "registry.ip":
			return ip
		case "registry.dist":
			return defaults.DefaultRegistryDist
		default:
			log.Fatalf("Not found argument %s in config", inp)
		}
		return ""
	}
	var fm = template.FuncMap{
		"parse": parse,
	}
	t := template.New("t").Funcs(fm)
	_, err = t.Parse(templateString)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, nil)
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func generateConfigFileFromViperTemplate(filePath string, templateString string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	parse := func(inp string) interface{} {
		result := viper.Get(inp)
		if result != nil {
			return result
		}
		log.Fatalf("Not found argument %s in config", inp)
		return ""
	}
	var fm = template.FuncMap{
		"parse": parse,
	}
	t := template.New("t").Funcs(fm)
	_, err = t.Parse(templateString)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, nil)
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//GenerateConfigFileFromViper is a function to generate yml from viper config
func GenerateConfigFileFromViper() error {
	configFile, err := DefaultConfigPath()
	if err != nil {
		log.Fatalf("fail in DefaultConfigPath: %s", err)
	}
	return generateConfigFileFromViperTemplate(configFile, defaultEnvConfig)
}

//GenerateConfigFileDiff is a function to generate diff yml for new context
func GenerateConfigFileDiff(filePath string, context *Context) error {
	return generateConfigFileFromTemplate(filePath, defaultEnvConfig, context)
}
