# HoT – Hono Test [![GitHub release](https://img.shields.io/github/release/ctron/hot.svg)](https://github.com/ctron/hot/releases)

This is a simple command line tool for testing Eclipse Hono™.

## Start a test consumer

Fill in your connection information, and then execute the following command:

    hot consume telemetry amqps://my.server:443 tenant

You can optionally use the following flags to configure the connection:

<dl>


<dt><code>--insecure</code></dt>
<dd>Set to true to enable Insecure TLS connection</dd>
<dt><code>--tlsPath</code></dt>
<dd>Set to path of trusted store file </dd>
<dt><code>--clientUsername</code></dt>
<dd>Tenant Username</dd>
<dt><code>--clientPassword</code></dt>
<dd>Tenant Password</dd>

*NOTE: if neither --insecure nor --tlsPath are set the AMQP client TLS default is used

</dl>

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

