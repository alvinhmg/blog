import React, { useState, useEffect } from 'react';
import { message, Spin, Input, Button, List, Popconfirm, Card, Typography, Space } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { categoryAPI } from '../../api';

const { Title } = Typography;

// 简单的 slug 生成函数
const slugify = (text) => {
  if (!text) return '';
  return text
    .toString()
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-') // Replace spaces with -
    .replace(/[^\w-]+/g, '') // Remove all non-word chars
    .replace(/--+/g, '-'); // Replace multiple - with single -
};

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
    
    const categorySlug = slugify(newCategory);
    if (!categorySlug) {
      message.warning('无法为分类名称生成有效的Slug');
      return;
    }

    try {
      // 同时发送 name 和 slug
      const response = await categoryAPI.createCategory({ name: newCategory, slug: categorySlug });
      if (response.success || response.code === 201 || response.code === 200) { // code 201 for created
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
    <div style={{ padding: '20px' }}>
      <Card>
        <Title level={2}>分类管理</Title>
        <Space style={{ marginBottom: 16 }}>
          <Input
            placeholder="新分类名称"
            value={newCategory}
            onChange={(e) => setNewCategory(e.target.value)}
            onPressEnter={handleAddCategory}
            style={{ width: 200 }}
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddCategory}>
            添加分类
          </Button>
        </Space>
        <List
          bordered
          dataSource={categories}
          renderItem={item => (
            <List.Item
              actions={[
                <Popconfirm
                  title="确定要删除这个分类吗？"
                  onConfirm={() => handleDeleteCategory(item.id)}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
                </Popconfirm>
              ]}
            >
              {item.name}
            </List.Item>
          )}
        />
      </Card>
    </div>
  );
};

export default CategoryManagementPage;