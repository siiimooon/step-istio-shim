package server

import (
	"context"
	"fmt"
	"github.com/siiimooon/istio-ca-shim-step/internal/common"
	"github.com/smallstep/certificates/api"
	"github.com/smallstep/certificates/ca"
	"go.step.sm/crypto/pemutil"
	"google.golang.org/grpc"
	securityapi "istio.io/api/security/v1alpha1"
	"istio.io/istio/pkg/security"
	"net"
)

func New() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) Start(caUrl, caFingerprint string) error {
	addr := ":9696"
	ctx := context.TODO()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %v", ":9696", err)
	}
	grpcServer := grpc.NewServer()
	fmt.Printf("server started. listening to %s", addr)
	securityapi.RegisterIstioCertificateServiceServer(grpcServer, s)

	go func() {
		<-ctx.Done()

		s.lock.Lock()
		s.ready = false
		s.lock.Unlock()

		grpcServer.GracefulStop()
	}()

	s.lock.Lock()
	s.ready = true
	s.lock.Unlock()

	s.ca = caUrl
	s.caFingerprint = caFingerprint
	return grpcServer.Serve(listener)
}

func (s *Server) CreateCertificate(ctx context.Context, request *securityapi.IstioCertificateRequest) (*securityapi.IstioCertificateResponse, error) {
	token, err := security.ExtractBearerToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed at extracting auth token from request")
	}
	return s.sign(token, request)
}

func (s *Server) sign(token string, request *securityapi.IstioCertificateRequest) (*securityapi.IstioCertificateResponse, error) {
	client, err := ca.NewClient(s.ca, ca.WithRootSHA256(s.caFingerprint))
	if err != nil {
		return nil, fmt.Errorf("failed at establishing connection to ca")
	}
	csrPEM, err := pemutil.ParseCertificateRequest([]byte(request.GetCsr()))
	if err != nil {
		return nil, fmt.Errorf("failed at parsing csr provided by request")
	}

	req := &api.SignRequest{
		CsrPEM: api.CertificateRequest{CertificateRequest: csrPEM},
		OTT:    token,
	}
	signResponse, err := client.Sign(req)
	if err != nil {
		return nil, fmt.Errorf("failed at signing certificate from provided csr")
	}

	crtResponse := securityapi.IstioCertificateResponse{
		CertChain: common.CertChainToPemChain(signResponse.CertChainPEM),
	}

	return &crtResponse, nil
}
