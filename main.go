// Command p12extract writes the certificate, the private key and the CA
// certificates of a PKCS#12 (.p12/.pfx) file to separate PEM files.
package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"os"

	"software.sslmate.com/src/go-pkcs12"
)

func main() {
	var p12 onceString
	flag.Var(&p12, "p12", "PKCS#12 file to read")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}
	if p12.value == "" || flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2)
	}

	if err := run(p12.value); err != nil {
		fmt.Fprintln(os.Stderr, "p12extract:", err)
		os.Exit(1)
	}
}

func run(p12 string) error {
	pfxData, err := os.ReadFile(p12)
	if err != nil {
		return err
	}
	password, err := readPassword()
	if err != nil {
		return err
	}

	key, cert, caCerts, err := pkcs12.DecodeChain(pfxData, password)
	if err != nil {
		return err
	}
	leaf, cas, err := findLeaf(key, append([]*x509.Certificate{cert}, caCerts...))
	if err != nil {
		return err
	}

	files, err := pemFiles(stem(p12), key, leaf, cas)
	if err != nil {
		return err
	}
	if err := writeAllOrNone(files); err != nil {
		return err
	}
	for _, f := range files {
		fmt.Println(f.path)
	}
	return nil
}
