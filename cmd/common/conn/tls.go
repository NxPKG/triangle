// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package conn

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/khulnasoft/triangle/cmd/common/config"
	"github.com/khulnasoft/triangle/pkg/defaults"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func grpcOptionTLS(vp *viper.Viper) (grpc.DialOption, error) {
	target := vp.GetString(config.KeyServer)
	if !(vp.GetBool(config.KeyTLS) || strings.HasPrefix(target, defaults.TargetTLSPrefix)) {
		return grpc.WithInsecure(), nil
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: vp.GetBool(config.KeyTLSAllowInsecure),
		ServerName:         vp.GetString(config.KeyTLSServerName),
	}

	// optional custom CAs
	caFiles := vp.GetStringSlice(config.KeyTLSCACertFiles)
	if len(caFiles) > 0 {
		ca := x509.NewCertPool()
		for _, path := range caFiles {
			certPEM, err := os.ReadFile(filepath.Clean(path))
			if err != nil {
				return nil, fmt.Errorf("cannot load cert '%s': %s", path, err)
			}
			if ok := ca.AppendCertsFromPEM(certPEM); !ok {
				return nil, fmt.Errorf("cannot process cert '%s': must be a PEM encoded certificate", path)
			}
		}
		tlsConfig.RootCAs = ca
	}

	// optional mTLS
	clientCertFile := vp.GetString(config.KeyTLSClientCertFile)
	clientKeyFile := vp.GetString(config.KeyTLSClientKeyFile)
	if clientCertFile != "" && clientKeyFile != "" {
		c, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &c, nil
		}
	}

	creds := credentials.NewTLS(&tlsConfig)
	return grpc.WithTransportCredentials(creds), nil
}
