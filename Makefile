gen_grpc_protoc:
	protoc \
	--go_out=grpc \
	--go_opt=paths=source_relative \
	--go-grpc_out=grpc \
	--go-grpc_opt=paths=source_relative \
	proto/*.proto
gen_key:
	openssl \
		req -x509 \
		-nodes \
		-days 365 \
		-newkey rsa:2048 \
		-keyout keys/server-payment/private.pem \
		-out keys/server-payment/public.pem \
		-config keys/server-payment/san.cfg
export_path:
	export PATH=$PATH:$(go env GOPATH)/bin