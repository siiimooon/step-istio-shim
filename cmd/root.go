package cmd

import (
	"github.com/siiimooon/istio-ca-shim-step/internal/cert"
	"github.com/siiimooon/istio-ca-shim-step/internal/monitoring"
	"github.com/siiimooon/istio-ca-shim-step/internal/server"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	caUrl         = ""
	caFingerprint = ""
	loglevel      = ""
	tlsCrt        = ""
	tlsKey        = ""
)

var rootCmd = &cobra.Command{
	Use: "istio-ca-shim-step",
	Run: func(cmd *cobra.Command, args []string) {
		monitoring.ServeMetrics()
		server, _ := server.New(monitoring.NewLogger(loglevel))

		getCertificate := cert.Loader(tlsCrt, tlsKey)
		_, _ = getCertificate(nil)
		_ = server.Start(caUrl, caFingerprint, getCertificate)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&caUrl, "ca-url", "", "url of step ca")
	rootCmd.Flags().StringVar(&caFingerprint, "ca-fingerprint", "", "fingerprint of root certificate of the ca")
	rootCmd.Flags().StringVar(&loglevel, "loglevel", "info", "loglevel of server")
	rootCmd.Flags().StringVar(&tlsCrt, "tls-crt", "", "path of the tls crt")
	rootCmd.Flags().StringVar(&tlsKey, "tls-key", "", "path of the tls key")

	err := rootCmd.MarkFlagRequired("ca-url")
	if err != nil {
		log.Panicf("failed at configuring flag for ca-url: %v", err)
	}
	err = rootCmd.MarkFlagRequired("ca-fingerprint")
	if err != nil {
		log.Panicf("failed at configuring flag for ca-fingerprint: %v", err)
	}
}
