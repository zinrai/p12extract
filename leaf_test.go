package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func newCertificate(t *testing.T, name string, key *ecdsa.PrivateKey) *x509.Certificate {
	t.Helper()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: name},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

func TestFindLeafReturnsCertificateOfKeyEvenWhenNotFirst(t *testing.T) {
	key := newKey(t)
	ca := newCertificate(t, "ca", newKey(t))
	leaf := newCertificate(t, "leaf", key)

	got, others, err := findLeaf(key, []*x509.Certificate{ca, leaf})
	if err != nil {
		t.Fatal(err)
	}
	if got != leaf {
		t.Errorf("leaf = %s, want leaf", got.Subject.CommonName)
	}
	if len(others) != 1 || others[0] != ca {
		t.Errorf("others = %v, want [ca]", others)
	}
}

func TestFindLeafFailsWhenNoCertificateMatchesKey(t *testing.T) {
	ca := newCertificate(t, "ca", newKey(t))

	if _, _, err := findLeaf(newKey(t), []*x509.Certificate{ca}); err == nil {
		t.Error("findLeaf succeeded with no certificate of the key")
	}
}
