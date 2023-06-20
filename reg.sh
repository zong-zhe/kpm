#!/bin/bash

# 安装 Apache 的 htpasswd
brew install httpd

# 创建存储用户密码的文件目录
mkdir -p registry_auth

# 使用 htpasswd 创建加密文件
htpasswd -Bbn test 1234 > registry_auth/htpasswd

# 检查是否已经有一个名为 registry 的容器存在
if [ "$(docker ps -aq -f name=registry)" ]; then
    # 停止并删除名为 registry 的容器
    docker stop registry
    docker rm registry
fi

# 启动带认证的 Docker Registry
docker run -p 5001:5000 \
--restart=always \
--name registry \
-v /var/lib/registry:/var/lib/registry \
-v $PWD/registry_auth/:/auth/ \
-e "REGISTRY_AUTH=htpasswd" \
-e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
-e "REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd" \
-d registry
