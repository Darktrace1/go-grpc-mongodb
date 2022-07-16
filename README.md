# Go-gRPC-MongoDB

## 기본설정
### Install Protoc Compiler (protoc v3.21.1)<br>
<ol>
  <li> Install https://github.com/protocolbuffers/protobuf/releases/tag/v21.1 and decompress
  <li> cd ~/Downloads/protobuf-3.21.1
  <li> Typing<br>

```sh
$ ./configure
$ make
$ make check
$ make install
```
  <li> Confirm protoc version

```console
protoc --version
```
</ol>

### Go plugins for the protocol compiler
<p><strong>Go plugins</strong> for the protocol compiler:</p>
<ol>
  <li>
    <p>Install the protocol compiler plugins for Go using the following commands:</p>

```go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
  </li>
  <li>
    <p>Update your <code>PATH</code> so that the protoc compiler can find the plugins:</p>

```sh
    $ export PATH="$PATH:$(go env GOPATH)/bin"
```
</ol>


Go module 생성

```Go
go mod init github.com/Darktrace1/go-grpc-mongodb
```


