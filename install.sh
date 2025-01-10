#!/bin/bash

# 颜色定义
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
PLAIN="\033[0m"

# 检查root权限
[[ $EUID -ne 0 ]] && echo -e "${RED}错误：请使用root用户运行此脚本${PLAIN}" && exit 1

# 一键安装函数
install_panel() {
    echo -e "${GREEN}开始安装 Hysteria2 Panel...${PLAIN}"
    
    # 1. 安装基础依赖
    if [[ -f /etc/debian_version ]]; then
        apt update
        apt install -y git golang-go mysql-server nginx curl wget
    else
        yum install -y git golang mysql-server nginx curl wget
    fi
    
    # 2. 创建工作目录
    mkdir -p /etc/hysteria2-panel/{cert,logs,configs}
    cd /etc/hysteria2-panel
    
    # 3. 克隆项目
    git clone https://github.com/Deepblue1997/hysteria2-panel.git .
    
    # 4. 配置MySQL
    systemctl start mysql
    systemctl enable mysql
    
    DB_PASS=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 16)
    mysql -e "CREATE DATABASE IF NOT EXISTS hysteria2_panel DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    mysql -e "CREATE USER IF NOT EXISTS 'hysteria2'@'localhost' IDENTIFIED BY '${DB_PASS}';"
    mysql -e "GRANT ALL PRIVILEGES ON hysteria2_panel.* TO 'hysteria2'@'localhost';"
    mysql -e "FLUSH PRIVILEGES;"
    
    # 5. 创建配置文件
    cat > backend/configs/config.json << EOF
{
    "listen": ":8080",
    "tls_cert_path": "/etc/hysteria2-panel/cert/cert.pem",
    "tls_key_path": "/etc/hysteria2-panel/cert/key.pem",
    "domain": "",
    "email": "",
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
    
    # 6. 编译后端
    cd backend
    go mod download
    go mod tidy
    go build -o hysteria2-panel
    
    # 7. 配置系统服务
    cat > /etc/systemd/system/hysteria2-panel.service << EOF
[Unit]
Description=Hysteria2 Panel Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/etc/hysteria2-panel/backend
ExecStart=/etc/hysteria2-panel/backend/hysteria2-panel
Restart=always
RestartSec=5
StandardOutput=append:/etc/hysteria2-panel/logs/access.log
StandardError=append:/etc/hysteria2-panel/logs/error.log

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable hysteria2-panel
    systemctl start hysteria2-panel
    
    # 8. 配置Nginx
    cat > /etc/nginx/conf.d/hysteria2-panel.conf << EOF
server {
    listen 80;
    server_name _;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF

    systemctl restart nginx
    
    # 9. 输出安装信息
    echo -e "${GREEN}安装完成！${PLAIN}"
    echo -e "${YELLOW}==================================${PLAIN}"
    echo -e "${YELLOW}MySQL密码: ${DB_PASS}${PLAIN}"
    echo -e "${YELLOW}面板地址: http://服务器IP${PLAIN}"
    echo -e "${YELLOW}==================================${PLAIN}"
}

# 执行安装
install_panel 