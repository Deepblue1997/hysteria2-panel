#!/bin/bash

# 颜色定义
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
PLAIN="\033[0m"

# 检查是否为root用户
if [[ $EUID -ne 0 ]]; then
    echo -e "${RED}错误：请使用root用户运行此脚本${PLAIN}"
    exit 1
fi

# 检查系统类型
if [[ ! -f /etc/os-release ]]; then
    echo -e "${RED}错误：不支持的操作系统${PLAIN}"
    exit 1
fi

# 读取系统版本信息
source /etc/os-release

# 安装基础依赖
install_base() {
    echo -e "${GREEN}开始安装基础依赖...${PLAIN}"
    if [[ ${ID} == "debian" || ${ID} == "ubuntu" ]]; then
        apt update
        apt install -y curl wget git nginx mysql-server build-essential
    elif [[ ${ID} == "centos" ]]; then
        yum install -y epel-release
        yum install -y curl wget git nginx mysql-server gcc
    else
        echo -e "${RED}错误：不支持的操作系统${PLAIN}"
        exit 1
    fi
}

# 安装Go环境
install_golang() {
    echo -e "${GREEN}开始安装Go环境...${PLAIN}"
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm -f go1.21.0.linux-amd64.tar.gz
    
    # 设置Go环境变量
    echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
    source /etc/profile.d/go.sh
}

# 安装acme.sh
install_acme() {
    echo -e "${GREEN}开始安装acme.sh...${PLAIN}"
    curl https://get.acme.sh | sh
    mkdir -p /etc/hysteria2-panel/cert
}

# 安装Hysteria2-Panel
install_panel() {
    echo -e "${GREEN}开始安装Hysteria2-Panel...${PLAIN}"
    
    # 创建工作目录
    mkdir -p /etc/hysteria2-panel
    cd /etc/hysteria2-panel
    
    # 创建必要的目录
    mkdir -p {configs,cert,logs}
    
    # 下载源码
    git clone https://github.com/your-repo/hysteria2-panel.git .
    
    # 编译后端
    cd backend
    go mod download
    go mod tidy
    go build -o hysteria2-panel
    
    # 创建配置文件
    cat > configs/config.json << EOF
{
    "listen": ":8080",
    "tls_cert_path": "/etc/hysteria2-panel/cert/cert.pem",
    "tls_key_path": "/etc/hysteria2-panel/cert/key.pem",
    "domain": "",
    "email": "",
    "jwt_secret": "$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32)",
    "database": {
        "type": "mysql",
        "host": "localhost",
        "port": 3306,
        "user": "hysteria2",
        "password": "$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 16)",
        "dbname": "hysteria2_panel"
    }
}
EOF
}

# 配置MySQL
setup_mysql() {
    echo -e "${GREEN}开始配置MySQL...${PLAIN}"
    
    # 启动MySQL服务
    systemctl start mysql
    systemctl enable mysql
    
    # 创建数据库和用户
    DB_PASS=$(grep -oP '"password": "\K[^"]+' /etc/hysteria2-panel/configs/config.json)
    mysql -e "CREATE DATABASE IF NOT EXISTS hysteria2_panel DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    mysql -e "CREATE USER IF NOT EXISTS 'hysteria2'@'localhost' IDENTIFIED BY '${DB_PASS}';"
    mysql -e "GRANT ALL PRIVILEGES ON hysteria2_panel.* TO 'hysteria2'@'localhost';"
    mysql -e "FLUSH PRIVILEGES;"
}

# 配置systemd服务
setup_service() {
    echo -e "${GREEN}开始配置系统服务...${PLAIN}"
    
    cat > /etc/systemd/system/hysteria2-panel.service << EOF
[Unit]
Description=Hysteria2 Panel Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/etc/hysteria2-panel/backend
ExecStart=/etc/hysteria2-panel/backend/hysteria2-panel
Environment=GIN_MODE=release
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
}

# 配置Nginx
setup_nginx() {
    echo -e "${GREEN}开始配置Nginx...${PLAIN}"
    
    cat > /etc/nginx/conf.d/hysteria2-panel.conf << EOF
server {
    listen 80;
    server_name _;
    
    # 静态文件缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
        expires 7d;
        proxy_pass http://127.0.0.1:8080;
    }
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_read_timeout 60s;
        proxy_send_timeout 60s;
    }
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src * data: 'unsafe-eval' 'unsafe-inline'" always;
}
EOF

    systemctl restart nginx
}

# 创建管理员账户
create_admin() {
    echo -e "${GREEN}创建管理员账户...${PLAIN}"
    
    ADMIN_PASS=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 12)
    
    mysql -e "USE hysteria2_panel; INSERT INTO users (email, password, role, status, created_at, updated_at) VALUES ('admin@example.com', '$(echo -n "${ADMIN_PASS}" | sha256sum | cut -d' ' -f1)', 'admin', 1, NOW(), NOW());"
    
    echo -e "${YELLOW}管理员账户创建成功：${PLAIN}"
    echo -e "${YELLOW}邮箱: admin@example.com${PLAIN}"
    echo -e "${YELLOW}密码: ${ADMIN_PASS}${PLAIN}"
}

# 主函数
main() {
    echo -e "${GREEN}开始安装Hysteria2-Panel...${PLAIN}"
    
    install_base
    install_golang
    install_acme
    install_panel
    setup_mysql
    setup_service
    setup_nginx
    create_admin
    
    echo -e "${GREEN}安装完成！${PLAIN}"
    echo -e "${YELLOW}==================================${PLAIN}"
    echo -e "${YELLOW}请修改 /etc/hysteria2-panel/configs/config.json 中的配置信息${PLAIN}"
    echo -e "${YELLOW}请及时修改默认管理员密码${PLAIN}"
    echo -e "${YELLOW}然后重启服务：systemctl restart hysteria2-panel${PLAIN}"
    echo -e "${YELLOW}==================================${PLAIN}"
}

# 执行主函数
main 