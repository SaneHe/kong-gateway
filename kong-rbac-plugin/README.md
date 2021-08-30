### kong-api-plugin

- 要在Go中编写Kong插件，您需要：
    - 定义结构类型以保存配置
    - 编写一个New()函数来创建您的结构实例
    - 在该结构上添加方法以处理事件 (`Certificate` `Rewrite` `Access` `Preread` `Log`)
    - 使用编译 docker run --rm -v $(pwd):/plugins kong/go-plugin-tool build <source>
    - 将生成的库（.so文件）放入go_plugins_dir目录中
    - plugins = bundled,nick-rate-limiting
    - go_plugins_dir = /etc/kong/plugins
    - go_pluginserver_exe = /usr/local/bin/kong-rbac-plugin

- docker run --rm -e GO111MODULE=on -e GOPROXY=https://goproxy.io -v /$(pwd):/plugins kong/go-plugin-tool:2.0.4-alpine-latest build -o rbac.so rbac.go 

### tips

- 注意：这里是个 [demo](https://github.com/Kong/go-plugins),请查看
- [官方包](https://github.com/Kong/go-pdk)
