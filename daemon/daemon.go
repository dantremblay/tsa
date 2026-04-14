package daemon

import (
	"fmt"
	"os"
	"strconv"

	"github.com/juliengk/go-utils/password"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/api/server"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/tls"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/spf13/cobra"
	"log/slog"
)

func serverInitConfig(appDir string) error {
	s, err := storage.NewDriver("sqlite", appDir)
	if err != nil {
		return err
	}
	defer s.End()

	if s.CountConfigKey("jwk") > 0 {
		slog.Info("Server initialization already done")

		return nil
	}

	s.AddConfig("jwk", token.GenerateJWK("", 24))
	s.AddConfig("auth_type", "none")
	s.AddConfig("admin_password", password.GeneratePassword("admin"))

	return nil
}

func runDaemon(cmd *cobra.Command, args []string) {
	var bindPort int
	var fqdn string

	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := serverInitConfig(cfg.App.Dir.Root); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	bindPort = serverBindPort
	if serverTLS && bindPort == 80 {
		bindPort = 443
	}

	if serverTLSGen && len(serverFQDN) == 0 {
		slog.Error("you must specified --fqdn if --tlsgen is enabled")
		os.Exit(1)
	}

	if len(serverFQDN) > 0 {
		fqdn = serverFQDN
	}

	// Input validations
	// IV - API Bind address
	if err := validation.IsValidIP(serverBindAddress); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// IV - API Port
	if err := validation.IsValidPort(bindPort); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Create API certificates
	conf := tls.New(serverTLSKey, serverTLSCert)

	if serverTLSGen {
		if !conf.CertificateExist() || (conf.CertificateExist() && conf.IsCertificateExpire()) {
			if err := conf.CreateSelfSignedCertificate(fqdn, serverTLSDuration); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		}
	}

	if serverTLS {
		if !conf.CertificateExist() {
			slog.Error("No certificate found")
			os.Exit(1)
		}
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer s.End()

	s.RemoveConfig("api_bind", "ALL")
	s.RemoveConfig("api_port", "ALL")
	s.RemoveConfig("api_fqdn", "ALL")
	s.AddConfig("api_bind", serverBindAddress)
	s.AddConfig("api_port", strconv.Itoa(bindPort))

	if len(serverFQDN) > 0 {
		s.AddConfig("api_fqdn", fqdn)
	}

	addr := fmt.Sprintf("%s:%d", serverBindAddress, bindPort)

	server.API(addr, serverTLS)
}
