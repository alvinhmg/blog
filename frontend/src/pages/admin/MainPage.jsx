import React from 'react';
import { Link } from 'react-router-dom';

const MainPage = () => {
  return (
    <div>
      <h1>管理界面</h1>
      <ul>
        <li><Link to="/admin/create-post">创建文章</Link></li>
        <li><Link to="/admin/delete-post">文章管理</Link></li> {/* 保留这个作为文章列表/管理入口 */}
        {/* <li><Link to="/admin/delete-post">删除文章</Link></li>  删除重复且指向相同的链接 */}
        <li><Link to="/admin/category-management">分类管理</Link></li>
        <li><Link to="/admin/tag-management">标签管理</Link></li>
      </ul>
    </div>
  );
};

export default MainPage;