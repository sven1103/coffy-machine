# Welcome to Coffy Machine

This is the documentation for Coffy's web service.

## Install

### Binaries

*Reference binaries for different OS architectures here.*

### Compile from source

Clone the repo first:

```bash
git clone git@github.com:sven1103/coffy-machine.git
```

Then check the required Go version in ``./go.mod`` (e.g. 'go 1.23'). 
Make sure you have the required Go version installed:

```bash
> go version  
go version go1.23.4 darwin/arm64
```

If not, visit the [Go website](https://go.dev/) and do so.

Coffy Machine uses swag for automated Swagger 2.0 documentation. So please install ``swag`` first:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

and then run in the projects root directory:

```bash
swag init .
```

This should create a ``./docs`` folder with this content:

```bash
.docs/
   docs.go
   swagger.json
   swagger.yaml
```

Then compile the source code with:

```bash
go build -o coffy-machine main.go
```

Then make it executable if not yet, e.g. under macOS or Linux:

```bash
chmod +x coffy-machine
```

### Verify installation

You can just try to run the web service with the example configuration ``example_config.yaml``:

```
./coffy-machine -c example_config.yaml
```
That's it &#x1F680;


