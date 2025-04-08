## 介绍

在容器化应用的开发与部署过程中，Docker 镜像是构建、传输和交付的核心载体。我们通常使用 docker push 命令将镜像推送到远程镜像仓库（如 Docker Registry 或 Harbor）。但在某些特定场景下，比如离线环境部署、镜像跨环境迁移，往往会使用已保存为 .tar 格式的镜像包。

这时，如果希望跳过 Docker 守护进程，**直接将 tar 包上传到镜像仓库**，image-push 就是一个非常实用的工具。

## **image-push 是什么？**

[image-push](https://github.com/lomtom/image-push) 是一个使用 Go 语言开发的轻量级 CLI 工具，用于将本地的 Docker 镜像 tar 包直接上传至 Harbor 或 Docker Registry，**无需解压加载到本地 Docker 引擎，也无需手动处理 Registry API 的复杂细节**。

其核心优势在于：**直接与镜像仓库 API 交互，实现镜像分层上传与 manifest 推送，完整还原镜像结构。**

## **功能特点**

- ✅ **支持认证**：通过用户名和密码登录受保护的镜像仓库。
- 📁 **支持项目指定**：镜像可上传至指定项目（如 Harbor 的 project）。
- 🔐 **跳过 TLS 校验**：适用于使用自签名证书的测试环境。
- 📦 **支持分块上传**：按配置的块大小进行上传，提升大镜像传输的稳定性。

## **快速开始**

### **📥 安装方式**

下载编译好的二进制包：[Releases 页面](https://github.com/lomtom/image-push/releases/tag/dev)

也可从源码自行构建：

### **🔧 源码编译**

```bash
git clone https://github.com/lomtom/image-push.git
cd image-push
go build -o image-push cmd/tool/main.go
```

### **🚀 使用示例**

```bash
./image-push \
  --address http://your-registry-address:5000 \
  --username your-username \
  --password your-password \
  --project your-project \
  --file /path/to/your-image.tar
```



**参数说明：**

| **参数**    | **说明**                                                   |
| ----------- | ---------------------------------------------------------- |
| --address   | 目标镜像仓库地址，例如 http://localhost:5000               |
| --username  | 登录用户名（可选）                                         |
| --password  | 登录密码（可选）                                           |
| --project   | 上传目标项目名（可选）                                     |
| --file      | 本地镜像 tar 包路径（可选）                                |
| --skipTls   | 是否跳过 TLS 校验，默认为 false（可选）                    |
| --chunkSize | 分块上传大小（单位：字节），如需开启分块上传则必填（可选） |

## **注意事项**

- 🔑 **权限校验**：确保提供的账号具备上传镜像的权限。
- 🌐 **网络可达**：本地环境需能正常访问目标 Registry。
- 🔄 **版本兼容性**：目标仓库需兼容 Docker Registry API v2。

在很多企业级应用场景中，我们可能无法依赖 Docker 守护进程来完成镜像的加载与推送操作。通过 image-push 工具，可以**大幅简化镜像 tar 包上传流程，提升自动化部署效率**，是离线部署、镜像迁移等场景下的优雅解决方案。
