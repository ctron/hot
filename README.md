# HoT – Hono Test [![GitHub release](https://img.shields.io/github/release/ctron/hot.svg)](https://github.com/ctron/hot/releases)

This is a simple command line tool for testing Eclipse Hono™.

## Start a test consumer

Fill in your connection information, and then execute the following command:

    hot consume telemetry amqps://my.server:443 tenant

You can use the following flags:

<dl>

<dt><code>--tlsConfig</code></dt>
<dd>0: To Skip the TLS verification. (Default)</dd>
<dd>1: To enable insecure TLS connection.</dd>
<dd>2: To enable secure TLS connection.</dd>
<dt><code>--tlsPath</code></dt>
<dd>Set to path of trusted store file </dd>
<dt><code>--clientUsername</code></dt>
<dd>Tenant Username (If Applicaple)</dd>
<dt><code>--clientPassword</code></dt>
<dd>Tenant Password (If Applicaple)</dd>


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

