自动生成读取服务端所有toml配置文件的go语言脚本:
1. 读取toml配置文件的go脚本, 被生成在目录ffGameConfig
2. 获取到目录ffGameConfig下的所有go文件
3. 生成读取所有toml配置文件的go脚本, 存放在ffGameConfig/ffGameConfigReader/main.go
4. 生成ffGameConfigReader
5. 拷贝ffGameConfigReader到可读取到toml配置文件的目录
6. 运行ffGameConfigReader, 查看运行状态