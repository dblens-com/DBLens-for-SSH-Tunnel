# DBLens for SSH Tunnel Manager

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/ssh-tunnel-manager)](https://github.com/dblens-com/DBLens-for-SSH-Tunnel)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

专业的SSH隧道管理工具，支持多端口转发和可视化控制


[![Go Report Card](https://github.com/dblens-com/DBLens-for-SSH-Tunnel/blob/main/image/demo.png?raw=true)]
## 功能特性

✅ **多隧道管理** - 同时管理多个SSH隧道  
✅ **交互式CLI** - 友好的命令行操作界面  
✅ **实时监控** - 显示隧道运行状态和时长  
✅ **自动重连** - 支持隧道重启和故障恢复  
✅ **安全传输** - 基于SSH协议的加密通信  
✅ **跨平台** - 支持Windows/Linux/macOS

## 安装指南

### 前置要求
- Go 1.18+ 环境
- SSH服务器访问权限

### 快速安装
```bash
go install github.com/dblens-com/DBLens-for-SSH-Tunnel@latest
```
## 手动构建
```bash
git clone https://github.com/dblens-com/DBLens-for-SSH-Tunnel.git
cd DBLens-for-SSH-Tunnel
go build -o DBLens-for-SSH-Tunnel main.go
```

# 使用说明
### 启动程序

```bash
DBLens-for-SSH-Tunnel
```

### 首次配置
1. 输入SSH服务器地址（格式：host:port）
2. 提供SSH用户名和密码
3. 验证连接成功后进入主菜单

### 管理隧道
```
[主菜单]
1. 添加新隧道    - 创建新的端口转发规则
2. 删除隧道      - 移除指定本地端口的隧道
3. 重启隧道      - 重新建立指定隧道连接
4. 删除所有隧道  - 清除所有转发规则
5. 退出程序      - 安全关闭所有连接
```
### 隧道配置示例
```text
本地端口: 8080
远程目标: internal-app:80

本地端口: 3306
远程目标: database:3306
```

### 操作演示

```ascii

SSH Tunnel Manager v2.0
------------------------------------------
当前隧道状态:
本地端口: 8080
远程目标: internal-app:80
运行状态: 运行中
运行时长: 2h35m

本地端口: 3306
远程目标: database:3306
运行状态: 运行中
运行时长: 1h12m
------------------------------------------
1. 添加新隧道
2. 删除隧道
3. 重启隧道
4. 删除所有隧道
5. 退出程序
```

### 注意事项
1. 确保本地防火墙允许监听端口
2. SSH服务器需开启TCP转发功能
3. 长时间连接建议使用SSH密钥认证
4. 程序退出时会自动清理所有隧道

### 贡献指南
欢迎通过Issue提交问题或PR贡献代码