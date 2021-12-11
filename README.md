Lumen IM 后端

## 目录结构说明

```
    app/                    应用目录
       cache/               缓存处理
       http/                http服务
            handler/        handler 处理
            middleware/     中间件
            response/       响应
            request/        请求
            router/         路由
       model/               model定义
       dao/                 dao定义
       service/             服务层 
       pkg/                 包
   config/                  配置文件
   resource/                资源目录
   runtime/                 运行目录，存放日志
```

### 图片域名部署

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