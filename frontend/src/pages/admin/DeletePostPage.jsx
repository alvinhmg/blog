import React, { useState, useEffect } from 'react';
import { Table, Button, message, Popconfirm, Input, Typography, Space, Card, Divider } from 'antd'; // 添加 Typography, Space, Card, Divider
import { EyeOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'; // 添加图标导入
import axios from 'axios';
import { Link } from 'react-router-dom'; // 添加这行导入
import { useNavigate } from 'react-router-dom';
import { postAPI } from '../../api';

const { Title } = Typography;

const DeletePostPage = () => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const fetchPosts = async () => {
    try {
      setLoading(true);
      // 添加查询参数以确保获取所有文章
      const response = await postAPI.getAllPosts({limit: 100});
      console.log('API响应 (经过拦截器处理后):', response);

      // 根据网络请求截图，数据在 response.data.posts 中
      let postsArray = [];
      if (response && response.data && Array.isArray(response.data.posts)) {
        postsArray = response.data.posts;
      } else {
        console.warn('未识别的文章数据结构或数据为空:', response);
        // 如果 postsArray 仍然为空，可以给用户一个更友好的提示
        if (postsArray.length === 0) {
            message.info('暂无文章数据');
        } else {
            message.warn('获取到的文章数据格式不正确，请检查API响应');
        }
      }

      // 确保每个post对象都有唯一的key
      const postsWithKeys = postsArray.map(post => ({
        ...post,
        key: post.id || `fallback-${Math.random().toString(36).substr(2, 9)}` // 添加前缀避免潜在冲突
      }));
      
      console.log('处理后的文章数据:', postsWithKeys);
      setPosts(postsWithKeys);
    } catch (error) {
      console.error('获取文章列表失败:', error);
      message.error('获取文章列表失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handleDelete = async (postId) => {
    try {
      // 添加错误处理和重试逻辑
      try {
        await postAPI.deletePost(postId);
        message.success('文章删除成功！');
      } catch (apiError) {
        console.error('API调用失败，尝试直接发送请求:', apiError);
        // 如果API调用失败，尝试直接使用fetch
        const response = await fetch(`http://localhost:8080/api/posts/${postId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        message.success('文章删除成功！');
      }
      
      // 重新获取文章列表
      fetchPosts();
    } catch (error) {
      console.error('删除文章失败:', error);
      message.error(error.response?.data?.message || error.message || '删除文章失败，请稍后重试');
    }
  };

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (text, record) => <Link to={`/posts/${record.id}`}>{text}</Link>,
    },
    {
      title: '作者',
      dataIndex: ['author', 'username'],
      key: 'author',
    },
    {
      title: '发布日期',
      dataIndex: 'created_at', // Corrected dataIndex from 'createdAt' to 'created_at'
      key: 'created_at', // Also update key for consistency
      render: (text) => {
        if (!text) return '无日期'; // Handle null or undefined dates
        const date = new Date(text);
        return isNaN(date.getTime()) ? '无效日期' : date.toLocaleDateString('zh-CN'); // Check for invalid date after parsing
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button 
            icon={<EyeOutlined />} 
            size="small"
            onClick={() => navigate(`/posts/${record.id}`)} // 使用 navigate 进行跳转
          >
            查看
          </Button>
          <Button 
            icon={<EditOutlined />} 
            size="small"
            onClick={() => navigate(`/admin/update-post/${record.id}`)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这篇文章吗？"
            description="删除后将无法恢复！"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button icon={<DeleteOutlined />} danger size="small">删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '20px' }}>
      <Card>
        <Title level={2}>文章管理</Title>
        <Divider />
        <Table 
          columns={columns} 
          dataSource={Array.isArray(posts) ? posts : []} 
          rowKey="id" 
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
      </Card>
    </div>
  );
};

export default DeletePostPage;