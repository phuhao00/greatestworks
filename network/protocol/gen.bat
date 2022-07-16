@echo  protocol generate
@echo current dir: %~dp0
set cDir= %~dp0
set SRC_DIR=%cDir%
set DST_DIR="%cDir%gen\"
@echo "SRC_DIR:"%SRC_DIR%
@echo "DST_DIR:"%DST_DIR%

rem protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto

rem protoc --doc_out=./doc --doc_opt=html,index.html proto/*.proto
goto r
protoc \
      -I . \
      -I ${GOPATH}/src \
      -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
      --go_out=":../generated" \
      --validate_out="lang=go:../generated" \
      example.proto
:r

mkdir gen

::protoc -I=./proto/ --go_out=./gen/ proto/*.proto

mkdir doc

protoc --validate_out="lang=go:./gen" --go_out=./gen/ --doc_out=./doc --doc_opt=html,index.html proto/*.proto

pause
