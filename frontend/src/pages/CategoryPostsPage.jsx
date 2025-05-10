import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { List, Typography, Spin, message, Card, Breadcrumb } from 'antd';
import { HomeOutlined, FolderOutlined } from '@ant-design/icons';
import { postAPI, categoryAPI } from '../api';

const { Title, Text } = Typography;

const CategoryPostsPage = () => {
  const [posts, setPosts] = useState([]);
  const [category, setCategory] = useState(null);
  const [loading, setLoading] = useState(true);
  const { id: categoryId } = useParams();

  useEffect(() => {
    const fetchCategoryAndPosts = async () => {
      if (!categoryId) return;
      setLoading(true);
      try {
        const categoryRes = await categoryAPI.getCategoryById(categoryId);
        if (categoryRes && (categoryRes.data || categoryRes.id)) {
          setCategory(categoryRes.data || categoryRes);
        } else {
          message.error('未找到该分类');
          setCategory(null);
        }

        const postsRes = await postAPI.getPostsByCategory(categoryId);
        if (postsRes && postsRes.data && Array.isArray(postsRes.data.posts)) {
          setPosts(postsRes.data.posts);
        } else if (postsRes && Array.isArray(postsRes)) { // Fallback for direct array response
          setPosts(postsRes);
        } else {
          setPosts([]);
        }
      } catch (error) {
        console.error('获取分类文章失败:', error);
        message.error('获取分类文章失败');
        setPosts([]);
        setCategory(null);
      } finally {
        setLoading(false);
      }
    };

    fetchCategoryAndPosts();
  }, [categoryId]);

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '200px' }}>
        <Spin tip="加载中..." />
      </div>
    );
  }

  return (
    <div style={{ padding: '20px' }}>
      <Card>
        <Breadcrumb style={{ marginBottom: '16px' }}>
          <Breadcrumb.Item href="/">
            <HomeOutlined />
            <span>首页</span>
          </Breadcrumb.Item>
          <Breadcrumb.Item>
            <FolderOutlined />
            <span>分类</span>
          </Breadcrumb.Item>
          {category && <Breadcrumb.Item>{category.name}</Breadcrumb.Item>}
        </Breadcrumb>

        {category ? (
          <Title level={2}>分类: {category.name}</Title>
        ) : (
          <Title level={2}>分类文章</Title>
        )}

        {posts.length > 0 ? (
          <List
            itemLayout="vertical"
            size="large"
            dataSource={posts}
            renderItem={post => (
              <List.Item
                key={post.id}
                actions={[
                  <Text type="secondary">
                    发布于: {post.created_at ? new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) : '未知日期'}
                  </Text>,
                ]}
                extra={post.cover_image_url && <img width={272} alt="logo" src={post.cover_image_url} />}
              >
                <List.Item.Meta
                  title={<Link to={`/posts/${post.id}`}>{post.title}</Link>}
                  description={post.excerpt || `${post.content.substring(0, 150)}...`}
                />
              </List.Item>
            )}
          />
        ) : (
          <Text>该分类下暂无文章。</Text>
        )}
      </Card>
    </div>
  );
};

export default CategoryPostsPage;