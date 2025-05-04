import React, { useState, useEffect } from 'react';

const CategoryManagementPage = () => {
  const [categories, setCategories] = useState([]);
  const [newCategory, setNewCategory] = useState('');

  useEffect(() => {
    fetch('/api/categories')
      .then(response => response.json())
      .then(data => setCategories(data));
  }, []);

  const handleAddCategory = () => {
    fetch('/api/categories', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name: newCategory }),
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          setCategories([...categories, data.category]);
          setNewCategory('');
        }
      });
  };

  const handleDeleteCategory = (categoryId) => {
    fetch(`/api/categories/${categoryId}`, {
      method: 'DELETE',
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          setCategories(categories.filter(category => category.id !== categoryId));
        }
      });
  };

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
        {categories.map(category => (
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