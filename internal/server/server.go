package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/siiimooon/istio-ca-shim-step/internal/common"
	"github.com/siiimooon/istio-ca-shim-step/internal/monitoring"
	"github.com/smallstep/certificates/api"
	"github.com/smallstep/certificates/ca"
	"go.step.sm/crypto/pemutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	securityapi "istio.io/api/security/v1alpha1"
	"istio.io/istio/pkg/security"
	"net"
)

func New(logger *zap.SugaredLogger) (*Server, error) {
	server := Server{}
	server.logger = logger
	return &server, nil
}

func (s *Server) Start(caUrl, caFingerprint string, getCertificate func(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error)) error {
	addr := ":9696"
	ctx := context.TODO()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %v", ":9696", err)
	}

	var grpcServer *grpc.Server
	var serverOpts []grpc.ServerOption
	if getCertificate != nil {
		creds := credentials.NewTLS(&tls.Config{
			GetCertificate: getCertificate,
		})
		serverOpts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer = grpc.NewServer(serverOpts...)
	s.logger.Infof("server started. listening to %s", addr)
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
	monitoring.IncProcessedRequests()
	token, err := security.ExtractBearerToken(ctx)
	if err != nil {
		err = fmt.Errorf("failed at extracting auth token from request")
		monitoring.IncFailedRequests()
		s.logger.Warnw(err.Error())
		return nil, err
	}
	certificateResponse, err := s.sign(token, request)
	if err != nil {
		monitoring.IncFailedRequests()
		s.logger.Warnw("failed at signing certificate", zap.Error(err))
	}
	return certificateResponse, err
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
