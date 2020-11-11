package netx

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
)

var errPEM = errors.New("unable to load PEM certs")

// LoadCACert returns a new CertPool with certificates loaded from given path.
func LoadCACert(path string) (*x509.CertPool, error) {
	ca := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(path) //nolint:gosec // False positive.
	if err == nil && !ca.AppendCertsFromPEM(caCert) {
		err = fmt.Errorf("%w: %q", errPEM, path)
	}
	if err != nil {
		return nil, err
	}
	return ca, nil
}
