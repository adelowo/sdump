# Deployment

Here is a quick run through of how I have deployed the version that
you can test online today.

## HTTP server

I already have a K8s cluster I use to run my side projects and
[Ayinke Ventures](https://ayinke.ventures) so I have deployed
the http component there.

> K8s is overkill for this by the way. Just take the binary
and run it please

- `k8s/infisical.yml`: I use Infisical to manage all my secrets, I find
it easier to selfhost than Hashicorp Vault. This syncs all the
env value to the namespace and store in a secret called `managed_secret`

- `k8s/http-server.yml`: Creates a deployment and service

- `k8s/ingress.yml`: set ups Nginx ingress and tls termination
for the service created above

- `k8s/update_deployment.sh`: Takes an image ID and updates the image of the
deployment.
Eg: `./deploy/k8s/update_deployment.sh e84a5c5f3b8724072d48f8b96f7794fb`

## SSH server

I run the SSH command on a small ec2 instance I use for miscellaneous things.
I use `screen` to run this.

```sh

apt install -y screen
screen -dmS ssh_server ./sdump ssh

```

If you'd rather go fancy, you can use Systemd as described below or even
K8s 👿👿👿👿

### Systemd?

If you want to run this over systemd, this config should work ideally

```sh
sudo vi /etc/systemd/system/sdump.service
```

```text

[Unit]
Description=sdump
After=network.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/
ExecStart=server ssh
Restart=on-failure

[Install]
WantedBy=multi-user.target

```

```sh
sudo systemctl daemon-reload
sudo systemctl start sdump
```

I already use Caddy so all i needed to do was extend the config as below:

```json
{
  "logging": {
    "logs": {
      "": {
        "level": "debug"
      }
    }
  },
  "apps": {
    "layer4": {
      "servers": {
        "sdump-ssh": {
          "listen": [
            "0.0.0.0:2222"
          ],
          "routes": [
            {
              "match": [
                {
                  "ssh": {}
                }
              ],
              "handle": [
                {
                  "handler": "proxy",
                  "upstreams": [
                    {
                      "dial": [
                        "localhost:3333"
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      }
    }
  }
}
```

You need to build caddy to get L4 support to make the above config work.
Here is an example that should work:

```sh
xcaddy build \
    --with github.com/mholt/caddy-l4/layer4 \
    --with github.com/mholt/caddy-l4/modules/l4tls \
    --with github.com/mholt/caddy-l4/modules/l4subroute \
    --with github.com/mholt/caddy-l4/modules/l4http \
    --with github.com/mholt/caddy-l4/modules/l4ssh \
    --with github.com/mholt/caddy-l4/modules/l4proxy \
    --with github.com/caddy-dns/duckdns
```

See [the documentation](https://github.com/mholt/caddy-l4) for more details
