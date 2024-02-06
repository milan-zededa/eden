package configitems

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lf-edge/eden/sdn/vm/api"
	"github.com/lf-edge/eden/sdn/vm/pkg/maclookup"
	"github.com/lf-edge/eve/libs/depgraph"
	"github.com/lf-edge/eve/libs/reconciler"
	log "github.com/sirupsen/logrus"
)

const (
	hostapdBinary        = "/usr/local/sbin/hostapd"
	hostapdStartTimeout  = 3 * time.Second
	hostapdStopTimeout   = 30 * time.Second
	hostapdConfParentDir = "/etc/hostapd"
	hostapdRunParentDir  = "/run/hostapd"
)

// PNAC : represents 802.1x authenticator performing port-based network access control
// for a given physical interface.
type PNAC struct {
	api.PNAC
	// PhysIf : target physical network interface for PNAC.
	PhysIf PhysIf
}

// Name returns MAC address of the physical interface as the unique identifier
// for the PNAC instance.
func (p PNAC) Name() string {
	return p.PhysIf.MAC.String()
}

// Label is used only for the visualization purposes of the config/state depgraphs.
func (p PNAC) Label() string {
	return p.PhysIf.LogicalLabel + " (PNAC)"
}

// Type assigned to PNAC
func (p PNAC) Type() string {
	return PNACTypename
}

// Equal is a comparison method for two equally-named PNAC instances.
func (p PNAC) Equal(other depgraph.Item) bool {
	p2, isPNAC := other.(PNAC)
	if !isPNAC {
		return false
	}
	return p.PNAC == p2.PNAC
}

// External returns false.
func (p PNAC) External() bool {
	return false
}

// String describes PNAC instance.
func (p PNAC) String() string {
	return fmt.Sprintf("PNAC: %#+v", p)
}

// Dependencies lists the physical interface as the only dependency.
func (p PNAC) Dependencies() (deps []depgraph.Dependency) {
	return []depgraph.Dependency{
		{
			RequiredItem: depgraph.ItemRef{
				ItemType: PhysIfTypename,
				ItemName: p.PhysIf.MAC.String(),
			},
			Description: "Underlying physical network interface must exist",
		},
	}
}

// PNACConfigurator implements Configurator interface for PNAC.
type PNACConfigurator struct {
	MacLookup *maclookup.MacLookup
}

// Create starts 802.1x authenticator for the given physical interface.
func (c *PNACConfigurator) Create(ctx context.Context, item depgraph.Item) error {
	pnac, isPNAC := item.(PNAC)
	if !isPNAC {
		return fmt.Errorf("invalid item type %T, expected PNAC", item)
	}
	netIf, found := c.MacLookup.GetInterfaceByMAC(pnac.PhysIf.MAC, false)
	if !found {
		err := fmt.Errorf("failed to get physical interface with MAC %v", pnac.PhysIf.MAC)
		log.Error(err)
		return err
	}
	if err := c.createHostapdConfFile(pnac, netIf.IfName); err != nil {
		return err
	}
	done := reconciler.ContinueInBackground(ctx)
	go func() {
		err := startHostapd(netIf.IfName)
		done(err)

		// TODO: block all non-EAP traffic until authenticated.
	}()
	return nil
}

func (c *PNACConfigurator) createHostapdConfFile(pnac PNAC, ifName string) error {
	if err := ensureDir(hostapdConfigDir(ifName)); err != nil {
		return err
	}
	// Prepare certificates for the EAP server integrated into hostapd.
	caCertPath := hostapdCACertPath(ifName)
	err := os.WriteFile(caCertPath, []byte(pnac.CACertPEM), 0600)
	if err != nil {
		err = fmt.Errorf("failed to write CA certificate %s: %w", caCertPath, err)
		log.Error(err)
		return err
	}
	srvCertPath := hostapdServerCertPath(ifName)
	err = os.WriteFile(srvCertPath, []byte(pnac.ServerCertPEM), 0600)
	if err != nil {
		err = fmt.Errorf("failed to write server certificate %s: %w", srvCertPath, err)
		log.Error(err)
		return err
	}
	srvKeyPath := hostapdServerKeyPath(ifName)
	err = os.WriteFile(srvKeyPath, []byte(pnac.ServerKeyPEM), 0600)
	if err != nil {
		err = fmt.Errorf("failed to write server key %s: %w", srvKeyPath, err)
		log.Error(err)
		return err
	}
	// Prepare user database.
	dbPath := hostapdUserDBPath(ifName)
	dbFile, err := os.Create(dbPath)
	if err != nil {
		err = fmt.Errorf("failed to create user DB file %s: %w", dbPath, err)
		log.Error(err)
		return err
	}
	defer dbFile.Close()
	switch pnac.EAPMethod {
	case api.EAPMethodTLS:
		dbFile.WriteString("\"client\" TLS")
	case api.EAPMethodTTLS_PAP:
		dbFile.WriteString("\"client\" TTLS")
		dbFile.WriteString(fmt.Sprintf("\"client\" TTLS-PAP \"%s\" [2]", pnac.Password))
	case api.EAPMethodTTLS_CHAP:
		dbFile.WriteString("\"client\" TTLS")
		dbFile.WriteString(fmt.Sprintf("\"client\" TTLS-CHAP \"%s\" [2]", pnac.Password))
	case api.EAPMethodTTLS_MSCHAPV2:
		dbFile.WriteString("\"client\" TTLS")
		dbFile.WriteString(fmt.Sprintf("\"client\" TTLS-MSCHAPV2 \"%s\" [2]", pnac.Password))
	}
	if err = dbFile.Sync(); err != nil {
		err = fmt.Errorf("failed to sync DB file %s: %w", dbPath, err)
		log.Error(err)
		return err
	}
	// Prepare hostapd config file.
	cfgPath := hostapdConfigPath(ifName)
	cfgFile, err := os.Create(cfgPath)
	if err != nil {
		err = fmt.Errorf("failed to create config cfgFile %s: %w", cfgPath, err)
		log.Error(err)
		return err
	}
	defer cfgFile.Close()
	cfgFile.WriteString(fmt.Sprintf("interface=%s\n", ifName))
	cfgFile.WriteString("driver=wired")
	cfgFile.WriteString("logger_stdout=-1")
	cfgFile.WriteString("logger_stdout_level=1")
	cfgFile.WriteString("ctrl_interface=/run/hostapd")
	cfgFile.WriteString("ieee8021x=1")
	// Default EAPoL version.
	eapolVer := uint8(2)
	if pnac.MACsec {
		eapolVer = uint8(3)
	}
	if pnac.EAPoLVersion != 0 {
		eapolVer = pnac.EAPoLVersion
	}
	cfgFile.WriteString(fmt.Sprintf("eapol_version=%d\n", eapolVer))
	cfgFile.WriteString(fmt.Sprintf("eap_reauth_period=%d\n", pnac.EAPReauthPeriod))
	cfgFile.WriteString("eap_server=0")
	cfgFile.WriteString(fmt.Sprintf("eap_user_file=%s\n", hostapdUserDBPath(ifName)))
	cfgFile.WriteString(fmt.Sprintf("ca_cert=%s\n", hostapdCACertPath(ifName)))
	cfgFile.WriteString(fmt.Sprintf("server_cert=%s\n", hostapdServerCertPath(ifName)))
	cfgFile.WriteString(fmt.Sprintf("private_key=%s\n", hostapdServerKeyPath(ifName)))
	macsec := "0"
	if pnac.MACsec {
		macsec = "1"
	}
	cfgFile.WriteString(fmt.Sprintf("macsec_policy=%s\n", macsec))
	// TODO: password?
	if err = cfgFile.Sync(); err != nil {
		err = fmt.Errorf("failed to sync config file %s: %w", cfgPath, err)
		log.Error(err)
		return err
	}
	return nil
}

// Modify is not implemented.
func (c *PNACConfigurator) Modify(_ context.Context, _, _ depgraph.Item) (err error) {
	return errors.New("not implemented")
}

// Delete stops 802.1x authenticator for the given physical interface.
func (c *PNACConfigurator) Delete(ctx context.Context, item depgraph.Item) error {
	pnac, isPNAC := item.(PNAC)
	if !isPNAC {
		return fmt.Errorf("invalid item type %T, expected PNAC", item)
	}
	netIf, found := c.MacLookup.GetInterfaceByMAC(pnac.PhysIf.MAC, false)
	if !found {
		err := fmt.Errorf("failed to get physical interface with MAC %v", pnac.PhysIf.MAC)
		log.Error(err)
		return err
	}
	done := reconciler.ContinueInBackground(ctx)
	go func() {
		err := stopHostapd(netIf.IfName)
		if err == nil {
			// ignore errors from here
			_ = removeHostapdConfDir(netIf.IfName)
			_ = removeHostapdRunDir(netIf.IfName)
		}
		done(err)
	}()
	return nil
}

// NeedsRecreate returns true, Modify is not implemented.
func (c *PNACConfigurator) NeedsRecreate(_, _ depgraph.Item) (recreate bool) {
	return true
}

func hostapdConfigDir(ifName string) string {
	return filepath.Join(hostapdConfParentDir, ifName)
}

func hostapdConfigPath(ifName string) string {
	return filepath.Join(hostapdConfigDir(ifName), "config")
}

func hostapdUserDBPath(ifName string) string {
	return filepath.Join(hostapdConfigDir(ifName), "userdb")
}

func hostapdCACertPath(ifName string) string {
	return filepath.Join(hostapdConfigDir(ifName), "ca.pem")
}

func hostapdServerCertPath(ifName string) string {
	return filepath.Join(hostapdConfigDir(ifName), "server.pem")
}

func hostapdServerKeyPath(ifName string) string {
	return filepath.Join(hostapdConfigDir(ifName), "server.key")
}

func hostapdRunDir(ifName string) string {
	return filepath.Join(hostapdRunParentDir, ifName)
}

func hostapdPidPath(ifName string) string {
	return filepath.Join(hostapdRunDir(ifName), "pid")
}

func hostapdLogPath(ifName string) string {
	return filepath.Join(hostapdRunDir(ifName), "log")
}

func removeHostapdConfDir(ifName string) error {
	dirPath := hostapdConfigDir(ifName)
	if err := os.RemoveAll(dirPath); err != nil {
		err = fmt.Errorf("failed to remove hostapd config dir %s: %w",
			dirPath, err)
		log.Error(err)
		return err
	}
	return nil
}

func removeHostapdRunDir(ifName string) error {
	dirPath := hostapdRunDir(ifName)
	if err := os.RemoveAll(dirPath); err != nil {
		err = fmt.Errorf("failed to remove hostapd run dir %s: %w",
			dirPath, err)
		log.Error(err)
		return err
	}
	return nil
}

func startHostapd(ifName string) error {
	if err := ensureDir(hostapdRunDir(ifName)); err != nil {
		return err
	}
	cfgPath := hostapdConfigPath(ifName)
	pidFile := hostapdPidPath(ifName)
	cmd := hostapdBinary
	logPath := hostapdLogPath(ifName)
	logFile, err := os.Create(logPath)
	if err != nil {
		err = fmt.Errorf("failed to create log file %s: %w", cfgPath, err)
		log.Error(err)
		return err
	}
	args := []string{
		"-B", "-d", "-t", "-P", pidFile,
		cfgPath,
	}
	return startProcess(MainNsName, cmd, args, logFile, pidFile,
		hostapdStartTimeout, true)
}

func stopHostapd(ifName string) error {
	pidFile := hostapdPidPath(ifName)
	return stopProcess(pidFile, hostapdStopTimeout)
}
