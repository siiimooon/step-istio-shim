package common

import (
	"encoding/pem"
	"github.com/smallstep/certificates/api"
)

func CertChainToPemChain(certificates []api.Certificate) (pems []string) {
	for _, cert := range certificates {
		certPem := pem.Block{
			Type:    "CERTIFICATE",
			Headers: nil,
			Bytes:   cert.Raw,
		}
		pems = append(pems, string(pem.EncodeToMemory(&certPem)))
	}
	return
}
