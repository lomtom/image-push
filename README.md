# Image Push

用于上传本地镜像到harbor或registry


**Parameter**
* address: registry or harbor server address,`required`
* file: docker image tar file path,`required`
* project: registry or harbor project name
* username: registry or harbor username.
* password: registry or harbor password.
* shipTls: skip ssl verify.
* chunkSize: default Monolithic Upload, if chunkSie is set,use Chunked Upload.

## Usage 1

### Run

```bash
# run
go run cmd/tool/main.go \
--address http://localhost:5000 \
--username admin \
--password admin@12345 \
--project libary \
--file ./resource/alpine.tar

go run cmd/tool/main.go \
--address http://localhost:5000 \
--file ./resource/alpine.tar
```


### Build(optional)

```bash
# build
go build -o bin/image-push cmd/tool/main.go

# run
./bin/image-push \
--address http://localhost:5000 \
--username admin \
--password admin@12345 \
--project libary \
--file ./resource/alpine.tar
```

## Usage 3

### Start Server

```bash
go run cmd/http/main.go
```

### Upload Image Tar

```bash
curl -v --request POST http://localhost:8080/upload \
-F "file=@./resource/alpine.tar" \
-F "address=http://localhost:5000" \
-F "username=admin" \
-F "password=admin@12345" \
-F "project=libary"
```

## Other

### Start Registry

```bash
docker run --name registry -d -p 5000:5000 registry:3.0.0
```
