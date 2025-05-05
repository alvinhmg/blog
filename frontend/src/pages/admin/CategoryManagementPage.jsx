import React, { useState, useEffect } from 'react';
import { message, Spin } from 'antd';
import { categoryAPI } from '../../api';

const CategoryManagementPage = () => {
  const [categories, setCategories] = useState([]);
  const [newCategory, setNewCategory] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCategories = async () => {
      setLoading(true);
      try {
        const response = await categoryAPI.getAllCategories();
        // 确保categories是一个数组
        const categoriesArray = Array.isArray(response) ? response : (response.data || []);
        setCategories(categoriesArray);
      } catch (error) {
        console.error('获取分类列表失败:', error);
        message.error('获取分类列表失败，请稍后重试');
      } finally {
        setLoading(false);
      }
    };
    
    fetchCategories();
  }, []);

  const handleAddCategory = async () => {
    if (!newCategory.trim()) {
      message.warning('分类名称不能为空');
      return;
    }
    
    try {
      const response = await categoryAPI.createCategory({ name: newCategory });
      if (response.success || response.code === 200) {
        const newCategoryData = response.category || response.data;
        setCategories(prev => Array.isArray(prev) ? [...prev, newCategoryData] : [newCategoryData]);
        setNewCategory('');
        message.success('添加分类成功');
      }
    } catch (error) {
      console.error('添加分类失败:', error);
      message.error('添加分类失败，请稍后重试');
    }
  };

  const handleDeleteCategory = async (categoryId) => {
    try {
      const response = await categoryAPI.deleteCategory(categoryId);
      if (response.success || response.code === 200) {
        setCategories(prev => Array.isArray(prev) ? prev.filter(category => category.id !== categoryId) : []);
        message.success('删除分类成功');
      }
    } catch (error) {
      console.error('删除分类失败:', error);
      message.error('删除分类失败，请稍后重试');
    }
  };

  if (loading) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <Spin size="large" tip="加载分类中..." />
      </div>
    );
  }

  return (
    <div>
      <h1>分类管理</h1>
      <input
        type="text"
        value={newCategory}
        onChange={(e) => setNewCategory(e.target.value)}
        placeholder="新分类名称"
      />
      <button onClick={handleAddCategory}>添加分类</button>
      <ul>
        {Array.isArray(categories) && categories.map(category => (
          <li key={category.id}>
            {category.name}
            <button onClick={() => handleDeleteCategory(category.id)}>删除</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default CategoryManagementPage;