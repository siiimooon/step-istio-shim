package cmd

import (
	"github.com/siiimooon/istio-ca-shim-step/internal/server"
	"github.com/spf13/cobra"
	"os"
)

var (
	caUrl         = ""
	caFingerprint = ""
	name          = ""
)

var rootCmd = &cobra.Command{
	Use: "istio-ca-shim-step",
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := server.New()
		_ = server.Start(caUrl, caFingerprint)
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
	rootCmd.MarkFlagRequired("ca-url")
	rootCmd.MarkFlagRequired("ca-fingerprint")
}
