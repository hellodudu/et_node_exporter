#!/bin/bash
echo "停止node_exporter服务..."
docker-compose -f ./config/node_exporter/docker-compose.yml down

echo "转换txt配置文件为docker-config.yml..."
./bin/exporter_mac -import_path=./config/scene/ -export_path=./config/node_exporter/

echo "启动nodex_exporter容器..."
docker-compose -f ./config/node_exporter/docker-compose.yml up -d

echo "success..."