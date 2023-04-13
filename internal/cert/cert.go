package cert

import (
	"crypto/tls"
	"log"
)

func Loader(tlsCrt, tlsKey string) func(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if tlsCrt != "" && tlsKey != "" {
		return func(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cer, err := tls.LoadX509KeyPair(tlsCrt, tlsKey)
			if err != nil {
				log.Fatal("failed at loading server certificate. terminating.")
			}
			return &cer, nil
		}
	}
	return nil
}
