import React from 'react';
import { Link } from 'react-router-dom';

const MainPage = () => {
  return (
    <div>
      <h1>管理界面</h1>
      <ul>
        <li><Link to="/admin/create-post">创建文章</Link></li>
        <li><Link to="/admin/update-post">更新文章</Link></li>
        <li><Link to="/admin/delete-post">删除文章</Link></li>
        <li><Link to="/admin/category-management">分类管理</Link></li>
        <li><Link to="/admin/tag-management">标签管理</Link></li>
      </ul>
    </div>
  );
};

export default MainPage;