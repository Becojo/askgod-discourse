module github.com/nsec/askgod-discourse

replace google.golang.org/grpc/naming => google.golang.org/grpc v1.29.1

go 1.16

require (
	github.com/gorilla/websocket v1.5.0
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/nsec/askgod v0.0.0-20220427021641-f0b2a15e2395
	github.com/open-policy-agent/opa v0.46.1
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli/v2 v2.23.5
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
