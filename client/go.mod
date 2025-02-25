module crypto-chad-client

go 1.22.1

require (
	crypto-chad-lib v0.0.0
	github.com/golang/protobuf v1.5.4
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
)

replace crypto-chad-lib => ../lib
