#/bin/bash

protoc --go_out=../ --go_opt=paths=source_relative --go-grpc_out=../ --go-grpc_opt=paths=source_relative *.service.proto
protoc --go_out=../ --go_opt=paths=source_relative gameserver.proto

# delete json omitempty
cd ../
ls *.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'