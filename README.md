# UTProxy - An uptime monitor proxy for internal services

[![Go Report Card](https://goreportcard.com/badge/github.com/axllent/utproxy)](https://goreportcard.com/report/github.com/axllent/utproxy)

UTProxy is a HTTP(S) proxy service for uptime monitors to access internal services without having to directly expose those services to the internet. It provides different internal checks (HTTP, TCP, MySQL or a command) and returns a HTTP response and status to the uptime monitor.


## Configuration

You have to set up a configuration file, see [`contrib/utproxy.yaml`](contrib/utproxy.yaml) for an example. Save this configuration in `/etc/utproxy.yaml`, or alternatively use the `-c` flag to specify a different configuration location.

The configuration has two main sections, firstly the service configuration:

```yaml
listen: 0.0.0.0:3500                                    # interface and port to listen on
#sslcert: /etc/letsencrypt/live/example.com/cert.pem    # SSL certificate (optional)
#sslkey: /etc/letsencrypt/live/example.com/privkey.pem  # SSL key (optional)
#log: /var/log/utproxy.log                              # log file (optional)
```

If both `sslcert` and `sslkey` are set, then UTProxy should be accessed via `https://`, otherwise `http://`. Inn this example we would be accessing the proxy via `http://example.com:3500`. UTProxy does not register or renew SSL certificates, so the service should be restarted manually if you update the certificates.

And then secondly the services you wish to test. Each service is added as a array to the `services:` section.

```yaml
services:
  # an array of services to test, see below
```

Each service must contain a unique "check key" (only a-z0-9 and - characters allowed), which will correspond to the URL on our UTProxy for the uptime monitor, eg: `http://example.com:3500/intranet`, `http://example.com:3500/smtp` etc.

Checks can be set up with one of the following types:

### `http`

A check for a HTTP response.

```yaml
services:
  intranet:                         # check key
    type: http                      # check type
    endpoint: http://192.168.0.10   # check url
    status: 200                     # expected response, default 200
    method: HEAD                    # request type (HEAD, GET, POST), default HEAD
```

### `tcp`

A check for a TCP connection.

```yaml
services:
  smtp:                             # check key
    type: tcp                       # check type
    endpoint: localhost:25          # check <destination>:<port>
```

### `mysql`

A check for a MySQL connection.

```yaml
services:
  database:                         # check key
    type: mysql                     # check type
    endpoint: localhost:3306        # mysql <destination>:<port (TCP only, no sockets)
    user: secretuser                # MySQL username 
    pass: secretpass                # MySQL password 
```

### `exec`

A check to run a command. The command should exit with a `0` status (success).

The following example is how to ping an internal machine:

```yaml
services:
  printer:                          # check key
    type: exec                      # check type
    command: ping                   # command to run
    args:                           # an optional array of command arguments
      - "-c"
      - "1"
      - "-W"
      - "2"
      - "-q"
      - "192.168.0.100"
```

Your `exec` check can be any command that the UTProxy daemon is allowed to run.


## Testing

You can test all your configured services by running `utproxy test`


## Setting up an uptime monitor

There are plenty of uptime monitors you can use both free and commercial, so pick one you are happy with. Some examples are:

- [HetrixTools](https://hetrix.tools/u-625253) (free plans)
- [Uptime Robot](https://uptimerobot.com/) (free plans)
- [Pingdom](https://www.pingdom.com/)

You need to set up your uptime monitors to monitor the HTTP status of each of your server checks.

`http://example.com:3500/intranet`, `http://example.com:3500/smtp`, `http://example.com:3500/database`, `http://example.com:3500/printer` etc

Checks should return a `200` status, else they are failing.


## Running as a systemd service

See the example [`utproxy.service`](contrib/utproxy.service).