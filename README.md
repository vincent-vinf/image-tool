# image-tool

## 特性

* 批量保存和加载镜像
* 支持自动在保存镜像时压缩镜像，并在加载时解压镜像
* 使用pigz代替gzip执行压缩和解压操作，由于[pigz](https://github.com/klauspost/pgzip)支持多线程，会大大提升压缩性能，并略微提升解压性能。
* 使用[container image库](github.com/containers/image)直接操作镜像，而不依赖docker daemon，在实际运维中，可以在无docker的跳板机执行load操作



## Quick Start

> 若镜像需要代理才能访问，可以在执行命令前配置环境变量
>
> export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7890

### image.txt

```
registry.cn-hangzhou.aliyuncs.com/qkd-system/workflow-controller:v3.5.3
# registry.cn-hangzhou.aliyuncs.com/qkd-system/vc-controller-manager:v1.8.1
registry.cn-hangzhou.aliyuncs.com/qkd-system/vc-scheduler:v1.8.1
# registry.cn-hangzhou.aliyuncs.com/qkd-system/vc-webhook-manager:v1.8.1
# registry.cn-hangzhou.aliyuncs.com/qkd-system/tensorboard:v0.1.0
```

按行分隔的镜像列表，使用#作为注释

### save

将镜像拉取下来，并作为tar包保存到指定目录当中。

支持的参数如下

```
  -d, --dir string        output dir (default "images")
  -i, --images string     images.txt path
  -p, --password string   docker registrypassword
      --platform string   image platform (default "linux/amd64")
  -u, --username string   docker registry username
      --zip               automatically compress image tar using pigz (default true)
```

#### 示例

```bash
image-tool save -d test-images -i images.txt -u <username> -p <passwd> --platform=linux/arm64
# 也可以不指定images.txt，而是使用参数直接传递镜像列表，若仓库可以公开读，则不需要指定用户名密码
image-tool save -d test-images image-tool:v0.1.0
# 也可以同时使用
image-tool save -d test-images -i images.txt image-tool:v0.1.0
```

程序会将images.txt列出的镜像下载到test-images目录，如果目录中已经有对应镜像的tar或tgz，则会跳过该镜像。

### sync

将image.txt中提及的镜像拉取下来，若已存在则跳过，若目录当中存在其他image.txt未提及的镜像，则会被删除。

支持的参数和save子命令相同，但不能通过参数直接传递镜像列表

#### 示例

```sh
image-tool sync -d test-images -i images.txt -u <username> -p <passwd> --platform=linux/arm64
```

### load

将目录中的镜像文件上传到指定的仓库

支持的参数为：

```
  -r, --registry string   example: harbor.qkd.cn:8443/library

  -d, --dir string        output dir (default "images")
  -i, --images string     images.txt path
  -p, --password string   docker registry password
  -u, --username string   docker registry username
```

* -r 指定需要推送的仓库，例如harbor.qkd.cn:8443/library，其中library是harbor的项目名称
* -d 指定镜像文件所在的目录
* -i image.txt文件，可以不指定
* -u 指定镜像仓库用户名
* -p 指定镜像仓库密码

### 示例

```sh
# 只有images.txt中指定的镜像会被推送到仓库
image-tool load -d ~/Downloads/test-images -i tmp/images.txt -r <registry>/library -u <username> -p <passwd>
# 也可以通过命令行参数直接指定需要推送到镜像，可以和images.txt一起使用，2者会被合并
image-tool load -d ~/Downloads/test-images -i tmp/images.txt -r <registry>/library -u <username> -p <passwd> nginx:latest
# 可以不指定任何镜像，程序会将目录中所有的tar文件都推送到仓库
image-tool load -d ~/Downloads/test-images -r <registry>/library -u <username> -p <passwd>
```

## todo

* containerd支持：containers-storage
* 校验tar



## 附录

* goreleaser测试：goreleaser release --snapshot --clean
* goreleaser发布：goreleaser release --skip=publish

### 使用的项目

* [goreleaser](https://github.com/goreleaser/goreleaser) 帮助构建和发布go程序的工具
* [skopeo](https://github.com/containers/skopeo) 实用镜像工具
* [cobra](https://github.com/spf13/cobra) 帮助命令行构建的库
* [pgzip](https://github.com/klauspost/pgzip) gzip的并发go实现