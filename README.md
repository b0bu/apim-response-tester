# apim-response-tester
Simple config to deploy apim backends with load balanced policy compiled from c#. Go services used to replicate internal api behaviour for testing within policies.

with qemu based lima, `qemu-user-static` package required for cross platform target compilation

login
```bash
podman login docker.io
```
build
```bash
make all VERSION=v0.0.x
make push VERSION=v0.0.x
```
apim
```bash
azure-apim-policy-compiler --s src/ --o target/
```
debug
```bash
. .env
bash debug.sh # get trace token
```
add to the request when tracing
```
Apim-Debug-Authorization: traceToken
Ocp-Apim-Subscription-Key: subToken
```
