package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
)

const (
	publicFilePerm  = 0o644
	privateFilePerm = 0o600
)

type pemFile struct {
	path string
	perm os.FileMode
	data []byte
}

func stem(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func pemFiles(stem string, key crypto.PrivateKey, leaf *x509.Certificate, cas []*x509.Certificate) ([]pemFile, error) {
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	files := []pemFile{
		{stem + ".crt", publicFilePerm, encodeCertificates([]*x509.Certificate{leaf})},
		{stem + ".key", privateFilePerm, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})},
	}
	// Not an empty CA file: a tool given an empty bundle fails far from here.
	if len(cas) > 0 {
		files = append(files, pemFile{stem + "-ca.crt", publicFilePerm, encodeCertificates(cas)})
	}
	return files, nil
}

func encodeCertificates(certs []*x509.Certificate) []byte {
	var data []byte
	for _, cert := range certs {
		data = append(data, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})...)
	}
	return data
}
