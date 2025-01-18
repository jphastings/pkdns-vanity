# PKDNS Vanity domain name generator

A tool that abuses your CPU to generate vanity [PKDNS domains](https://github.com/pubky/pkdns) starting (or ending) with characters of your choosing.

More than 4 characters will take minutes on a fast modern machine, more than 5 can take hours. More than 6 and you'll only generate heat.

```sh
$ pkdns-vanity pkdns
Looking for a key with 'pkdns' at the start…
  pkdnsuds3gpyn4c9rsjypbts7u38fd44skonpf53bpked9b7dngo <- yyyu4xsmmw4bbxjjsqccz6mspwq5yydpz475m8pkqjixfibpis6y
  Public Key <- Private Key
```

## Installing

On a mac or linux you can install with [homebrew](https://brew.sh):

```sh
brew install jphastings/tools/pkdns-vanity
```

If you have [go](https://golang.org) installed you can:

```sh
go install github.com/jphastings/pkdns-vanity@latest
```

Or download binaries [from this repo](https://github.com/jphastings/pkdns-vanity/releases/latest).
