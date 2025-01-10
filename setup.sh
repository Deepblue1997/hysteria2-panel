#!/bin/bash

# 颜色定义
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
PLAIN="\033[0m"

# 检查root权限
[[ $EUID -ne 0 ]] && echo -e "${RED}错误：请使用root用户运行此脚本${PLAIN}" && exit 1

# 设置工作目录
WORK_DIR="/etc/hysteria2-panel"

# 克隆项目
setup_project() {
    echo -e "${GREEN}正在克隆项目...${PLAIN}"
    
    # 安装git
    if ! command -v git &> /dev/null; then
        if [[ -f /etc/debian_version ]]; then
            apt update && apt install -y git
        else
            yum install -y git
        fi
    fi
    
    # 克隆项目
    mkdir -p ${WORK_DIR}
    git clone https://github.com/YOUR_USERNAME/hysteria2-panel.git ${WORK_DIR}
    
    # 创建必要的目录
    mkdir -p ${WORK_DIR}/{cert,certs,logs,configs}
}

# 安装Go环境
install_golang() {
    echo -e "${GREEN}正在安装Go环境...${PLAIN}"
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm -f go1.21.0.linux-amd64.tar.gz
    
    # 设置环境变量
    echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
    source /etc/profile.d/go.sh
}

# 编译项目
build_project() {
    echo -e "${GREEN}正在编译项目...${PLAIN}"
    cd ${WORK_DIR}/backend
    go mod download
    go mod tidy
    go build -o hysteria2-panel
}

# 创建配置文件
create_config() {
    echo -e "${GREEN}正在创建配置文件...${PLAIN}"
    
    # 生成随机密码
    DB_PASS=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 16)
    JWT_SECRET=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32)
    
    cat > ${WORK_DIR}/backend/configs/config.json << EOF
{
    "listen": ":8080",
    "tls_cert_path": "${WORK_DIR}/cert/cert.pem",
    "tls_key_path": "${WORK_DIR}/cert/key.pem",
    "domain": "",
    "email": "",
    "jwt_secret": "${JWT_SECRET}",
    "database": {
        "type": "mysql",
        "host": "localhost",
        "port": 3306,
        "user": "hysteria2",
        "password": "${DB_PASS}",
        "dbname": "hysteria2_panel"
    }
}
EOF
}

# 配置systemd服务
setup_service() {
    echo -e "${GREEN}正在配置系统服务...${PLAIN}"
    
    cat > /etc/systemd/system/hysteria2-panel.service << EOF
[Unit]
Description=Hysteria2 Panel Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=${WORK_DIR}/backend
ExecStart=${WORK_DIR}/backend/hysteria2-panel
Environment=GIN_MODE=release
Restart=always
RestartSec=5
StandardOutput=append:${WORK_DIR}/logs/access.log
StandardError=append:${WORK_DIR}/logs/error.log

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable hysteria2-panel
}

# 启动服务
start_service() {
    echo -e "${GREEN}正在启动服务...${PLAIN}"
    systemctl start hysteria2-panel
    
    # 检查服务状态
    if systemctl is-active --quiet hysteria2-panel; then
        echo -e "${GREEN}服务启动成功！${PLAIN}"
    else
        echo -e "${RED}服务启动失败，请检查日志${PLAIN}"
        journalctl -u hysteria2-panel -n 50
    fi
}

# 主函数
main() {
    echo -e "${GREEN}开始部署Hysteria2-Panel...${PLAIN}"
    
    setup_project
    install_golang
    build_project
    create_config
    setup_service
    start_service
    
    echo -e "${GREEN}部署完成！${PLAIN}"
    echo -e "${YELLOW}==================================${PLAIN}"
    echo -e "${YELLOW}请修改配置文件：${WORK_DIR}/backend/configs/config.json${PLAIN}"
    echo -e "${YELLOW}查看运行日志：${WORK_DIR}/logs/access.log${PLAIN}"
    echo -e "${YELLOW}查看错误日志：${WORK_DIR}/logs/error.log${PLAIN}"
    echo -e "${YELLOW}==================================${PLAIN}"
}

# 执行主函数
main 