# Update Home IP

Keep an A record in Route53 up to date with your public IP.

## Installation

```
go install github.com/jbaikge/update-home-ip
```

## Running With Systemd

Copy the included `update-home-ip.service` to `/etc/systemd/system` and modify with the correct path to the binary and domain name parameter. Create a new file, `/etc/defaults/update-home-ip` with the correct values:

```shell
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=
AWS_HOSTED_ZONE_ID=
```

Copy the included `update-home-ip.timer` to `/etc/systemd/system` and modify if desired.
