#!/bin/bash

# 设置编译参数
PROJECT_NAME="dblens-for-ssh-tunnel"      # 更改为你的项目名称
VERSION="1.0.0"           # 设置你的版本号
OUTPUT_DIR="build"        # 输出目录

# 支持的平台和架构列表
PLATFORMS=(
    "windows/amd64"
    "darwin/amd64"        # Intel Mac
    "darwin/arm64"        # M1/M2 Mac
    "linux/amd64"
)

# 创建输出目录
mkdir -p $OUTPUT_DIR

# 遍历所有平台进行编译
for platform in "${PLATFORMS[@]}"
do
    # 分割平台配置
    GOOS=${platform%/*}
    GOARCH=${platform#*/}

    # 设置输出文件名
    OUTPUT_NAME=$PROJECT_NAME-$VERSION-$GOOS-$GOARCH

    # 设置文件扩展名
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi

    # 设置编译参数
    echo "正在编译 $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o $OUTPUT_DIR/$OUTPUT_NAME .

    # 检查编译结果
    if [ $? -ne 0 ]; then
        echo "编译 $GOOS/$GOARCH 失败!"
        exit 1
    else
        echo "编译 $GOOS/$GOARCH 完成!"
    fi
done

echo "所有平台编译完成！输出目录：$OUTPUT_DIR"