

## 特性

* 批量保存和加载镜像
* 支持自动在保存镜像时压缩镜像，并在加载时解压镜像
* 使用pigz代替gzip执行压缩和解压操作，由于[pigz](https://github.com/klauspost/pgzip)支持多线程，会大大提升压缩性能，并略微提升解压性能。
* 使用[container image库](github.com/containers/image)直接操作镜像，而不依赖docker daemon，在实际运维中，可以在无docker的跳板机执行load操作



## 支持的命令

### save

将镜像拉取下来，并作为tar包保存到指定目录当中。

支持的参数如下

```
  -d, --dir string        output dir (default "images")
  -i, --images string     images.txt path
  -p, --password string   password
      --platform string   image platform (default "linux/amd64")
  -u, --username string   username
      --zip               automatically compress image tar using pigz (default true)
```

#### 示例



### sync

### load

### clean



todo

* containerd支持：containers-storage
* 



## 附录

* goreleaser测试：goreleaser release --snapshot --clean
* goreleaser发布：