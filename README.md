# HOT - Hono Test

This is a simple command line tool for testing Eclipse Hono™.

## Start a test consumer

Fill in your connection information, and then execute the following command:

    hot consume telemetry amqps://my.server:443 tenant

You can use the following flags:

<dl>

<dt><code>--insecure</code></dt>
<dd>Skip the TLS verification.</dd>

</dl>

## Publish an HTTP message

Fill in your connection information, and then execute the following command:

    hot publish http telemety https://my.server tenant device auth password payload

The following flags are supported:

<dl>

<dt><code>--qos</code></dt>
<dd>Set the "Quality of Service". Defaults to `0`.</dd>

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
<dd>Set the "Quality of Service". Defaults to `0`.</dd>

</dl>


## Building

Building requires Go 1.12.x. You can build the binary by executing:

    GO111MODULE=on go build -o hot ./cmd

