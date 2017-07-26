echo off & color 0A
:: 项目名称
set PROJECT=openssl
:: 版本标签 github上可查 :https://github.com/openssl/openssl/releases
set VESION=OpenSSL_1_1_0-stable
:: 项目路径
set PROJECT_PATH=%cd%
:: 代码存放路径
set CODE_PATH="%PROJECT_PATH%\%PROJECT%_%VESION%"
:: github openssl 项目网址
set OPENSSL_GIT_URL=https://github.com/openssl/openssl.git
::安装路径
set OPENSSL_INSTALL_DIR=%cd%

::从github上按照指定版本拉取源码
::需要已安装git
if not exist "%CODE_PATH%" (
git clone -b %VESION% https://github.com/openssl/openssl.git %CODE_PATH%
)

cd /d "%CODE_PATH%"

::通过perl脚本根据配置生成makefile
::需要已安装好perl（strawberry）
perl Configure VC-WIN32 --prefix=%OPENSSL_INSTALL_DIR% no-asm

:: 设置VS工具集目录,取决于电脑中VS安装路径
set VS_DEV_CMD="D:\Program Files (x86)\Microsoft Visual Studio 14.0\Common7\Tools\VsDevCmd.bat"
call %VS_DEV_CMD%
:: 编译
nmake -f makefile
:: 测试(可选)
nmake test
:: 安装
nmake install

pause