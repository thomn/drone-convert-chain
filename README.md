# drone-convert-chain

A conversion extension to chain conversion plugins. _Please note this project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=TARGET_ENDPOINTS=https://convert1:8433,https://convert2:8443 \
  --env=TARGET_SECRETS=4fce87fb3a28d462331dff6e1bb9c98b,2fd99e9939e4774ff9dfda816cbf93a9 \
  --env=TARGET_SKIP_VERIFIES=false,false \
  --restart=always \
  --name=converter thomn/drone-convert-chain
```

Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_CONVERT_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
