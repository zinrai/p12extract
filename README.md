# p12extract

p12extract writes the certificate, the private key and the CA certificates
of a PKCS#12 (.p12/.pfx) file to separate PEM files.

## Usage

```
p12extract -p12 client.p12
```

The password is read from the terminal. The files are written to the current
directory, named after the input file:

| File | Content |
|---|---|
| `client.crt` | The certificate that belongs to the private key |
| `client.key` | The private key, as unencrypted PKCS#8 |
| `client-ca.crt` | The other certificates; not written when there are none |

p12extract refuses to overwrite an existing file, and writes either all of
the files or, when one cannot be written, none of them.

The private key is written without encryption. On Linux the file is created
with mode `0600`; on Windows it inherits the ACL of the directory it is
written to.

## Supported files

p12extract reads PKCS#12 files encrypted with PBES2 (AES) as well as the legacy
RC2 and 3DES schemes that OpenSSL 3 reads only with `-legacy`. Files must be
DER-encoded; BER-encoded files are rejected.

## License

This project is licensed under the [MIT License](./LICENSE).
