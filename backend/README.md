# Alvin's Blog 后端

## 项目概述

这是一个基于Golang的Hertz框架和MySQL的个人博客网站后端API服务。

## 当前进度

目前已完成的功能：

1. 基础项目结构搭建
2. 数据库模型设计与定义
   - 用户模型
   - 文章模型
   - 分类模型
   - 标签模型
   - 评论模型
3. 用户认证功能
   - 用户注册
   - 用户登录
   - JWT认证中间件
4. 基础文章API
   - 获取文章列表
   - 获取文章详情

## 项目结构

```
/backend
├── api/          # API处理函数
│   ├── auth.go   # 认证相关API
│   └── post.go   # 文章相关API
├── config/       # 配置文件
│   └── database.go # 数据库配置
├── middleware/   # 中间件
│   └── auth.go   # 认证中间件
├── models/       # 数据模型
│   └── models.go # 数据库模型定义
├── routes/       # 路由定义
│   └── routes.go # API路由
├── utils/        # 工具函数（待实现）
├── .env          # 环境变量配置
├── go.mod        # Go模块定义
└── main.go       # 入口文件
```

## 如何运行

1. 确保已安装Go 1.20或更高版本
2. 确保MySQL服务已启动
3. 创建数据库：

```sql
CREATE DATABASE alvin_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. 修改`.env`文件中的数据库配置
5. 运行服务：

```bash
cd backend
go run main.go
```

服务将在`http://localhost:8080`启动

## API文档

### 认证API

#### 注册用户

```
POST /api/auth/register
```

请求体：
```json
{
  "username": "用户名",
  "email": "邮箱",
  "password": "密码",
  "nickname": "昵称"
}
```

#### 用户登录

```
POST /api/auth/login
```

请求体：
```json
{
  "username": "用户名",
  "password": "密码"
}
```

### 文章API

#### 获取文章列表

```
GET /api/posts?page=1&page_size=10
```

#### 获取文章详情

```
GET /api/posts/:id
```

## 下一步计划

1. 实现文章的创建、更新和删除功能
2. 实现分类和标签管理功能
3. 实现评论功能
4. 实现搜索功能
5. 开发前端界面
6. 部署到服务器