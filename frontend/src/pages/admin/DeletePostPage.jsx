import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Popconfirm, message, Typography, Card, Divider } from 'antd';
import { DeleteOutlined, EditOutlined, EyeOutlined } from '@ant-design/icons';
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
      console.log('获取到的文章数据:', response);
      
      // 更健壮的数据处理
      let postsArray = [];
      if (Array.isArray(response)) {
        postsArray = response;
      } else if (response && typeof response === 'object') {
        if (Array.isArray(response.data)) {
          postsArray = response.data;
        } else if (response.posts && Array.isArray(response.posts)) {
          postsArray = response.posts;
        }
      }
      
      // 确保每个post对象都有唯一的key
      const postsWithKeys = postsArray.map(post => ({
        ...post,
        key: post.id || Math.random().toString(36).substr(2, 9)
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
      render: (text, record) => <a href={`/posts/${record.id}`} target="_blank" rel="noopener noreferrer">{text}</a>,
    },
    {
      title: '作者',
      dataIndex: ['author', 'username'],
      key: 'author',
    },
    {
      title: '发布日期',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (text) => new Date(text).toLocaleDateString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button 
            icon={<EyeOutlined />} 
            size="small"
            onClick={() => window.open(`/posts/${record.id}`, '_blank')}
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