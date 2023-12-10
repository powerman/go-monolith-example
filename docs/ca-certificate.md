# Create local CA to issue localhost HTTPS certificates

You can check [How to securely test local/staging HTTPS
project](securely-test-local.md)
for more details about required setup or just follow instructions below.

**WARNING:** You'll need to run these commands just once, don't run them
again if you already did this before for some other project.

MacOS users should first prepare OpenSSL package:
```
brew install openssl
export EASYRSA_OPENSSL="$(ls -1 $(brew --prefix)/bin/openssl | sort -n -t/ -k6 | tail -n1)"
```

Install EasyRSA into `~/.easyrsa/` to generate local CA and website
certificates:
```
mkdir -p ~/.easyrsa &&
  curl -L https://github.com/OpenVPN/easy-rsa/releases/download/v3.1.6/EasyRSA-3.1.6.tgz |
  tar xzvf - --strip-components=1 -C ~/.easyrsa
```

Create local CA for signing certificates for local websites plus
Diffie-Hellman parameter for DHE cipher suites:
```
cd ~/.easyrsa
./easyrsa init-pki
echo Local CA $(hostname -f) | ./easyrsa build-ca nopass
openssl dhparam 2048 | install -m 0600 /dev/stdin pki/private/dhparam2048.pem
```

Now import local CA certificate `~/.easyrsa/pki/ca.crt` into your browser:

- MacOS: You can easily add the certificate as a trusted certificate
  authority for the currently logged in user:
  `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.easyrsa/pki/ca.crt`
- Linux:
    - Chrome-based browsers: go to chrome://settings/certificates,
      AUTHORITIES, IMPORT, select file, check "Trust this certificate for
      identifying websites", OK.
    - Firefox, command-line tools (curl, etc.):
      `sudo mkdir -p /usr/local/share/ca-certificates && sudo cp ~/.easyrsa/pki/ca.crt /usr/local/share/ca-certificates/ && sudo chmod 0644 /usr/local/share/ca-certificates/ca.crt && sudo update-ca-certificates`
