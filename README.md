# Lumen IM 后端

## 1、简介

- 地址： [http://im.gzydong.club](http://im.gzydong.club)
- 账号： 18798272054 或 18798272055
- 密码： admin123

## 2、项目DEMO

- 地址： [http://im.gzydong.club](http://im.gzydong.club)
- 账号： 18798272054 或 18798272055
- 密码： admin123

## 3、环境部署

##### Nginx 后端代理

```nginx
# http 代理
upstream imhttp {
    server 127.0.0.1:8080;
}

# websocket 代理
upstream imwss {
    server 127.0.0.1:8080;
}

server {
    listen       443 ssl;
    server_name  api.xxxx.com;

    ssl_certificate             /etc/nginx/cert/www.domain.com/server.crt;
    ssl_certificate_key         /etc/nginx/cert/www.domain.com/server.key;
    ssl_session_cache           shared:SSL:1m;
    ssl_protocols               TLSv1.1 TLSv1.2;
    ssl_session_timeout         5m;
    ssl_ciphers                 ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;
    ssl_prefer_server_ciphers   on;

    # http 转发
    location / {
        client_max_body_size    20m;

        # 将客户端的 Host 和 IP 信息一并转发到对应节点
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        # 转发Cookie，设置 SameSite
        proxy_cookie_path / "/; secure; HttpOnly; SameSite=strict";

        # 执行代理访问真实服务器
        proxy_pass http://imhttp;
    }

    # Websocket 转发
    location /wss/ {
        # WebSocket Header
        proxy_http_version 1.1;
        proxy_set_header Upgrade websocket;
        proxy_set_header Connection "Upgrade";

        # 将客户端的 Host 和 IP 信息一并转发到对应节点
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        # 客户端与服务端无交互 60s 后自动断开连接，请根据实际业务场景设置
        proxy_read_timeout 180s;

        # 执行代理访问真实服务器
        proxy_pass http://imwss;
    }
}
```

##### 图片域名部署

```nginx
server {
    listen 80;
    server_name im-img.local-admin.com;
    index  index.html;

    location / {
        # 项目文件上传目录
        root /path/to/../../uploads;
    }

    # 私有目录禁止访问
    location /private {
        deny all;
    }

    location ~ .*\.(gif|jpg|jpeg|png|bmp|swf|flv)$ {
        # 设置缓存过期时间
        expires 30d;
        
        # 关闭访问日志
        access_log off;
    }
}
```

### 待开发功能

1. 群搜索功能 ok
4. 群模糊搜索和入群验证 ok
5. 消息免打扰 ok
6. 可以设置管理员，转让群主
7. 退出/解散群组 ok
8. 送达，可以显示聊天信息的送达状态
9. 阅读，聊天对象的阅读状态查看
10. 语音会议
11. 视频会议
12. 笔记转发好友
13. 加入搜索名字查找好友