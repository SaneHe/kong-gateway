.PHONY: all

all: rbac go-pluginserver

rbac: rbac.go
	go mod tidy && go build -o rbac .

rbac.so: rbac.so
	go build -buildmode=plugin -o rbac.so rbac.go

go-pluginserver:
	go get -d -v github.com/Kong/go-pluginserver && go build github.com/Kong/go-pluginserver