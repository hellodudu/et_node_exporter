@echo off
REM 后续命令使用的是：UTF-8编码
chcp 65001
echo .
echo 停止node_exporter服务...
taskkill /F /IM windows_exporter.exe

REM "转换txt配置文件为windows_config.yml..."
.\bin\exporter_win.exe -import_path=.\config\scene\ -export_path=.\config\node_exporter\

echo .
echo 启动node_exporter服务...
.\bin\windows_exporter.exe --config.file=.\config\node_exporter\windows_config.yml

echo .
echo success...
pause