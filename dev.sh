#!/bin/sh
(trap 'kill 0' SIGINT; $(cd client;yarn serve) & go run . & caddy run --config Caddyfile)
