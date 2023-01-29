
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto <br>
protoc -I=. --go_out=. ./greeter.proto <br>
./protoc.exe  -I=. --go_out=. --go-grpc_out=. greeter.proto <br>    

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest   <br>
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  <br>
//注意你的GOOS不能错误
