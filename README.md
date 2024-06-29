# whois

## Introduction

This is a simple program that connects to a whois server, specified in the CLI with the --host
option, with the port defaulting to the well known port 43 (but configurable with --port option)
and the domain to query being the --domain option.

An exercise in connecting to a simple TCP server and reading the response.

## Usage

In it's basic form, the code will default to the IANA whois server and will query for the name passed
as a CLI parameter:

- whois some.domain.name
- whois AS12345
- whois 1.2.3.4

The following options modify this behaviour:

* --host some.host.or.ip.address
  specifies a host to connect to instead of the default IANA server

* --port 43
  specify a port to connect to instead of the default 43

* --timeout 2
  specify a network connection timeout instead of the default 2 seconds

* --version
  output the code version

* --revision
  output more verbose code revision information

* --debug
  output debug information

## ToDo

* Add referral following

## Bug Reporting

Please use the issues feature, or, of course, offer to contribute.

## Licence and Copyright

This code is Copyright (c) 2024 Karl Dyson.

All rights reserved.

## Warranty

There's no warranty that this code is safe, secure, or fit for any purpose.

I'm not responsible if you don't read the code, check its suitability for your
use case, and wake up to find it's eaten your cat...
