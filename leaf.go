package main

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
)

// Not DecodeChain's leaf: it takes the first certificate bag, which need not
// be the one that belongs to the key.
func findLeaf(key crypto.PrivateKey, certs []*x509.Certificate) (leaf *x509.Certificate, others []*x509.Certificate, err error) {
	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported private key type %T", key)
	}
	public, ok := signer.Public().(interface{ Equal(crypto.PublicKey) bool })
	if !ok {
		return nil, nil, fmt.Errorf("unsupported public key type %T", signer.Public())
	}

	for i, cert := range certs {
		if public.Equal(cert.PublicKey) {
			others = append(append(others, certs[:i]...), certs[i+1:]...)
			return cert, others, nil
		}
	}
	return nil, nil, errors.New("no certificate matches the private key")
}
