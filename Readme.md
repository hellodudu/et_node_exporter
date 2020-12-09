# 物理主机监控组件
    该组件暴露物理主机监控指标给监控中间件，默认启用端口9200，windows和linux主机启动服务不同，所用配置文件也不同

## 目录结构

- 依赖组件:
    * `docker` 
    * `docker-compose` 

- 需要主机开启对外端口：
    * 9200 (`node_exporter`)

- node_exporter windows服务配置文件:
    * `config/node_exporter/windows_config.yml`

- node_exporter linux服务配置文件:
    * `config/node_exporter/docker-compose.yml`


## 开启node_exporter服务
1. 修改对应主机系统的配置文件，`windows`系统修改`config/node_exporter/windows_config.yml`文件中的`service-where`字段，`linux`系统修改`docker-compose.yml`文件中的`hostname`字段，标识此主机名字
2. 执行不同系统对应的`start`脚本，脚本会将`txt`配置文件合并转换为`node_exporter`可读取的配置文件，并且开启`node_exporter`的`docker`镜像
3. 配置文件有更改后，直接运行`start`脚本，服务将重启，不会对其他线上服务造成任何影响