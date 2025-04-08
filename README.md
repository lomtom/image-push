# Image Push

用于上传本地镜像压缩包到harbor或registry

**原理**
通过调用registry的api接口，实现上传镜像压缩包到registry或harbor

**参数**
* address: registry 或 harbor 服务地址，例如：http://localhost:5000,`required`
* file: 本地压缩包路径 ，例如：/resource/alpine.tar,`required`
* project: 项目地址，例如：library
* username: 用户名
* password: 密码
* skipTls: 是否跳过ssl认证，默认为false
* chunkSize: 默认整体上传，如果设置了chunkSie，请使用Chunked Upload

## Usage 1

### Run

```bash
# run
go run cmd/tool/main.go \
--address http://localhost:5000 \
--file ./resource/alpine.tar


go run cmd/tool/main.go \
--address http://localhost:5000 \
--username admin \
--password admin@12345 \
--project library \
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
--project library \
--file ./resource/alpine.tar
```

## Usage 2

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
-F "project=library"
```

## Other

### Start Registry

```bash
docker run --name registry -d -p 5000:5000 registry:3.0.0
```
