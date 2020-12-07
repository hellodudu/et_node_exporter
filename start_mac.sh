#!/bin/bash
echo "停止node_exporter服务..."
docker-compose down

echo "启动nodex_exporter容器..."
docker-compose up -d

echo "success..."