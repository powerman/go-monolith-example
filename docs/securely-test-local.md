# How to securely test local/staging HTTPS project

Modern projects often support HTTPS and HTTP/2, moreover they can use `Strict-Transport-Security:` and
`Content-Security-Policy:` headers which result in different behaviour for HTTP and HTTPS versions, or
even completely forbid HTTP version. To develop and test such project locally, on CI, and at staging server
we either have to provide a way to access it using HTTP in non-production environments (bad idea) or
somehow make it work with HTTPS everywhere.

> HTTP in non-production environments is a bad idea because we'll test not the same thing which will runs
on production, and because there is a chance to occasionally keep HTTP enabled on production too.

Quick and dirty way to provide HTTPS everywhere is to either create self-signed certificate or company's
own CA and use it to sign certificate for a project, and then include this certificate in a project's repo.
This also doesn't work well because self-signed certificates are very inconvenient to use (endless browser
warnings for everyone), while company's CA has to be installed into and trusted by every employee's browser,
which opens possibility to use that CA to issue certificate for any website and use MitM attack on any
employee to analyse/modify all his traffic. This became even worse in case some developers/testers are freelancers,
who works on projects for many different companies, and for sure won't like to install each company's CA
into browser on their own workstation.

So, let's do it right way:
- Each developer/tester who wants to run project **locally will use his own CA**.
- CI will use it's own CA.
- **Staging will runs on public domain with real certificate** (e.g., using free Let's Encrypt).

## How to get certificate to run project locally

I'll use `~/.easyrsa/` here, but feel free to change it to whatever you like to.

### Once, on each developer's/tester's workstation

#### Install EasyRSA tool
This command may need GNU tar on OSX, or you can unpack archive in any other way you like.
```sh
mkdir -p ~/.easyrsa &&
  curl -L https://github.com/OpenVPN/easy-rsa/releases/download/v3.1.2/EasyRSA-3.1.2.tgz |
  tar xzvf - --strip-components=1 -C ~/.easyrsa
```
Same command can be used to upgrade EasyRSA later (just change version in url).

#### Create your CA
```sh
cd ~/.easyrsa
./easyrsa init-pki
echo Local CA $(hostname -f) | ./easyrsa build-ca nopass
```
If you like to password-protect your CA to ensure it won't be used to issue certificates (trusted by
your browser) if someone manage to steal it from your workstaion, then remove `nopass`.

Now import `~/.easyrsa/pki/ca.crt` into your browser's CA list.

#### Create certificate for localhost
This is optional, but local projects often run on https://localhost:projectport/, and they all may
use this certificate, so let's create it, just in case.
```sh
./easyrsa --days=3650 "--subject-alt-name=IP:127.0.0.1,DNS:localhost,DNS:*.localhost" build-server-full localhost nopass
```
I've included `*.localhost` to let you add something like `www.localhost` to your `/etc/hosts` and
run project on https://www.localhost instead of https://localhost - this may be important because
browsers handle cookies differently for domains with/without dots in some cases.
**Update:** Looks like this doesn't work for `*.localhost` for some reason (at least in Chromium - probably
it handle it like `*.com` and reject as too unsafe), but it does work for `*.project.localhost`.

### For each project running on localhost
```sh
cp ~/.easyrsa/pki/issued/localhost.crt /path/to/project.crt
cp ~/.easyrsa/pki/private/localhost.key /path/to/project.key
```
If you use docker to run project then you can bind-mount certificate instead of copying it:
```sh
docker run --name nginx -d -p 8080:80 \
  -v /path/to/your/nginx/conf.d:/etc/nginx/conf.d:ro \
  -v ~/.easyrsa/pki/issued/localhost.crt:/etc/nginx/ssl/server.crt:ro \
  -v ~/.easyrsa/pki/private/localhost.key:/etc/nginx/ssl/server.key:ro \
  nginx:alpine
```

### For each project running on unique local domain
Replace `project.home.arpa` with your local domain used to run this project.
```sh
cd ~/.easyrsa
./easyrsa build-server-full project.home.arpa nopass
cp ~/.easyrsa/pki/issued/project.home.arpa.crt /path/to/project.crt
cp ~/.easyrsa/pki/private/project.home.arpa.key /path/to/project.key
```
Now you may need to add project.home.arpa to your `/etc/hosts` or local DNS.
You can use this to provide access to local project for another devices in your LAN (like smartphone).
If you'll want to use your LAN's IP 192.168.0.42 instead of domain (editing `/etc/hosts` on smartphone
may not be easy, and running local DNS too) then create certificate this way:
```sh
./easyrsa --days=3650 "--subject-alt-name=IP:192.168.0.42,DNS:project.home.arpa" build-server-full project.home.arpa nopass
```

## How to setup CI

Create CA and certificate for project in same way as local.
(Either create them each time you run build in CI, or create just once and set them as CI's environment variables.)

Next, you'll have to make your `~/.easyrsa/pki/ca.crt` trusted by OS. How this should be done depends on your
linux distributive (e.g., on Ubuntu you'll need to `cp ca.crt /usr/local/share/ca-certificates/ca.crt && update-ca-certificates`).

## How to restrict access to staging environment

There is no easy way to hide staging website's domain name and get real certificate for it, so just deal with it.

It is possible to get wildcard certificate for public website https://example.com and use it on
https://staging.example.com, which resolve to internal IP in company's LAN and thus restrict public access.
This also isn't the best solution, because you'll have to copy very powerful wildcard certificate for your
main website to less secure staging server, and because you won't be able to provide access to staging from
the internet if you'll need this in the future.

Best option is to run staging on public domain and real public IP, accessible from the internet, and then
restrict access using your webserver (e.g., nginx) configuration. This way you can keep access to path
`http://staging.example.com/.well-known/acme-challenge/` (used by Let's Encrypt to issue certificates) open
from the internet, but allow (e.g., by IP or using HTTP Basic auth) access to the rest of staging website to
your employees only.
