# apim-response-tester
Simple config to deploy apim backends with load balanced policy compiled from c#. Go services used to replicate internal api behaviour for testing within policies.

with qemu based lima, `qemu-user-static` package required for cross platform target compilation

login
```
podman login docker.io
```
build
```
make all VERSION=v0.0.x
make push VERSION=v0.0.x
```
apim
```
azure-apim-policy-compiler --s src/ --o target/
```
