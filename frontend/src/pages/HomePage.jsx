import React, { useState, useEffect } from 'react';
import { Typography, Row, Col, Card, Space, Button, Divider, Spin, message } from 'antd';
import { ReadOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { postAPI } from '../api'; // 引入 postAPI
import './HomePage.css'; // 引入自定义 CSS

const { Title, Paragraph, Text } = Typography;

const HomePage = () => {
  const [latestPosts, setLatestPosts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLatestPosts = async () => {
      try {
        setLoading(true);
        // 获取最新的几篇文章，例如最新的3篇
        const response = await postAPI.getAllPosts({ page: 1, limit: 3, sort: 'created_at', order: 'desc' });
        if (response && response.data && Array.isArray(response.data.posts)) {
          setLatestPosts(response.data.posts);
        } else {
          console.warn('获取最新文章数据结构不正确或为空:', response);
          setLatestPosts([]); // 设置为空数组以避免渲染错误
        }
      } catch (error) {
        console.error('获取最新文章失败:', error);
        message.error('加载最新文章失败');
        setLatestPosts([]); // 出错时也设置为空数组
      } finally {
        setLoading(false);
      }
    };

    fetchLatestPosts();
  }, []);

  return (
    <div className="home-page-container">
      {/* Hero Section */}
      <Row justify="center" align="middle" className="hero-section">
        <Col span={24} style={{ textAlign: 'center' }}>
          <Title level={1} className="hero-title">欢迎来到 Alvin 的博客</Title>
          <Paragraph className="hero-subtitle">
            这是一个分享技术、思考和创意的个人空间。
          </Paragraph>
          <Button type="primary" size="large" icon={<ReadOutlined />} className="hero-button">
            <Link to="/posts" style={{ color: 'inherit', textDecoration: 'none' }}>浏览文章</Link>
          </Button>
        </Col>
      </Row>

      <Divider />

      {/* Latest Posts Section */}
      <Row justify="center" className="latest-posts-section">
        <Col xs={24} lg={20} xl={18}>
          <Title level={2} style={{ textAlign: 'center', marginBottom: '2rem' }}>最新文章</Title>
          <Spin spinning={loading}>
            {latestPosts.length > 0 ? (
              <Row gutter={[24, 24]}>
                {latestPosts.map(post => (
                  <Col xs={24} sm={12} md={8} key={post.id}>
                    <Card
                      hoverable
                      className="post-card"
                      title={<Link to={`/posts/${post.id}`} className="post-card-title">{post.title}</Link>}
                      extra={<Link to={`/posts/${post.id}`}>阅读更多</Link>}
                    >
                      <Paragraph ellipsis={{ rows: 3 }} className="post-card-summary">{post.excerpt || post.summary /* 优先使用 excerpt */}</Paragraph>
                      <Text type="secondary">发布于: {new Date(post.created_at).toLocaleDateString('zh-CN')}</Text>
                    </Card>
                  </Col>
                ))}
              </Row>
            ) : (
              !loading && <Paragraph style={{ textAlign: 'center' }}>暂无最新文章</Paragraph>
            )}
          </Spin>
        </Col>
      </Row>
    </div>
  );
};

export default HomePage;