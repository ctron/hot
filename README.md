# HoT – Hono Test [![GitHub release](https://img.shields.io/github/release/ctron/hot.svg)](https://github.com/ctron/hot/releases)

This is a simple command line tool for testing Eclipse Hono™.

## Authentication

### Username and password

You can set the username and password for all operations using the `--username`
and `--password` parameters.

When publishing data the username is normally a combination of `<auth-id>@<tenant>`.
As you are already providing the tenant, you can use the `--auth-id` parameter
instead, which will internally generate the correct user name, by adding the
tenant suffix. 

<dl>
<dt><code>--auth-id,-a</code></dt>
<dd>The username to use for authenticating with the backend.</dd>
<dt><code>--username,-u</code></dt>
<dd>The full username to use for authenticating with the backend.</dd>
<dt><code>--password,-p</code></dt>
<dd>The password to use for authenticating with the backend.</dd>
</dl>

Assuming you have a tenant `foo` and an authentication id of `auth1`, then
you can use either:

    --username auth1@foo

Or:

    --auth-id auth1

### X.509 client certificates

It is possible to use X.509 client certificates, instead of
username/password authentication. For this you can use the parameters:

<dl>
<dt><code>--client-key</code></dt>
<dd>The path to an X.509 PKCS#8 encoded private key</dd>
<dt><code>--client-cert</code></dt>
<dd>The path to an file containing a PEM encoded client certificate chain</dd>
</dl>


## Start a test consumer

Fill in your connection information, and then execute the following command:

    hot consume amqps://my.server:443 tenant

Or:

    hot consume amqps://username@password:my.server:443 tenant

You can use the following flags:

<dl>

<dt><code>--disable-tls</code></dt>
<dd>Disable TLS negotiation on the AMQP connection</dd>

<dt><code>--insecure</code></dt>
<dd>Set to true to enable Insecure TLS connection</dd>

<dt><code>--cert</code></dt>
<dd>Path to the certificate bundle in PEM format (overrides system CA certs)</dd>

</dl>

### Telemetry & event

Running `consume` will consume both *telemetry* and *event* messages and
output them on the console.

### Command & control

You can enable command and control handling, by using the switch
`-c, --command`. You can pass an optional value to the switch,
which is the command name, `TEST` is being used by default.

When this feature is enabled and a message containing the `ttd` property
is received, it will try to get the next command, and forward it to
the device.

The source of the command can be specified with the `-r,--reader` argument.
The following readers are available:

* `ondemand` – When a command is required, it will show a prompt on the
  console, which can read the command payload. When the command requests
  times out, the prompt will get canceled.
* `prefill` – A prompt on the console allows putting in the command payload
  for the next command request.
* `static:<payload>` – Everything after the prefix `static:` will be used
  as the command payload. No interactive prompt is being presented. 

## Publish an HTTP message

Fill in your connection information, and then execute the following command:

    hot publish http telemety https://my.server tenant device payload --username auth --password password

The following additional flags are supported:

<dl>

<dt><code>--qos</code></dt>
<dd>Set the "Quality of Service". Defaults to <code>0</code>.</dd>

<dt><code>--ttd</code></dt>
<dd>Set the "time till disconnect", the amount of seconds the HTTP call will
wait for a command to the device</dd>

</dl>

## Publish an MQTT message

Fill in your connection information, and then execute the following command:

    hot publish mqtt telemety ssl://my.server tenant device payload

The following additional flags are supported:

<dl>

<dt><code>--qos</code></dt>
<dd>Set the "Quality of Service". Defaults to <code>0</code>.</dd>

<dt><code>--ttd</code></dt>
<dd>Set the "time till disconnect", the amount of seconds the MQTT call will
wait for a command to the device</dd>

</dl>

## Building

Building requires Go 1.13.x. You can build the binary by executing:

    GO111MODULE=on go build -o hot ./cmd

