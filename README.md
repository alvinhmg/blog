# Alvin's Blog

这是一个基于Golang(Hertz框架)和MySQL的个人博客网站项目。

## 项目结构

```
/
├── backend/         # 后端Go代码
│   ├── api/        # API处理函数
│   ├── config/     # 配置文件
│   ├── middleware/ # 中间件
│   ├── models/     # 数据模型
│   ├── routes/     # 路由定义
│   ├── utils/      # 工具函数
│   └── main.go     # 入口文件
├── frontend/       # 前端代码
│   ├── public/     # 静态资源
│   └── src/        # 源代码
└── scripts/        # 部署脚本
```

## 技术栈

- 后端：Golang + Hertz框架
- 数据库：MySQL
- 前端：React + Ant Design
- 部署：Docker + Nginx

## 功能特性

- 博客文章的发布、编辑、删除
- 文章分类和标签管理
- 用户注册、登录和权限管理
- 评论和点赞功能
- 响应式设计，支持移动端和桌面端
- 全文搜索功能

## 部署信息

网站将部署在alvinhmg.top域名下。