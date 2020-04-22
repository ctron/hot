# HoT – Hono Test [![GitHub release](https://img.shields.io/github/release/ctron/hot.svg)](https://github.com/ctron/hot/releases)

This is a simple command line tool for testing Eclipse Hono™.

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

<dt><code>--username</code></dt>
<dd>Tenant username (if required)</dd>

<dt><code>--password</code></dt>
<dd>Tenant password (if required)</dd>

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

    hot publish http telemety https://my.server tenant device auth password payload

The following flags are supported:

<dl>

<dt><code>--qos</code></dt>
<dd>Set the "Quality of Service". Defaults to <code>0</code>.</dd>

<dt><code>--ttd</code></dt>
<dd>Set the "time till disconnect", the amount of seconds the HTTP call will
wait for a command to the device</dd>

</dl>

## Publish an MQTT message

Fill in your connection information, and then execute the following command:

    hot publish mqtt telemety ssl://my.server tenant device auth password payload

The following flags are supported:

<dl>

<dt><code>--qos</code></dt>
<dd>Set the "Quality of Service". Defaults to <code>0</code>.</dd>

</dl>


## Building

Building requires Go 1.12.x. You can build the binary by executing:

    GO111MODULE=on go build -o hot ./cmd

