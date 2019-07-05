# HOT - Hono Test

This is a simple command line tool for testing Eclipse Hono™.

## Start a test consumer

    hot consume telemetry amqps://my.server:443 tenant

You can use `--insecure` in case you want to skip TLS verification.

## Publish an HTTP message

    hot publish http telemety https://my.server tenant device auth password payload

You can use the following flags:

<dl>

<dt><code>--qos</code></dt>
<dd>Set the "Quality of Service". Defaults to `0`.</dd>

<dt><code>--ttd</code></dt>
<dd>Set the "time till disconnect", the amount of seconds the HTTP call will
wait for a command to the device</dd>

</dl>