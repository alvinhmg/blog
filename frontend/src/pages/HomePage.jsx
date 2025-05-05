import React, { useState, useEffect } from 'react';
import { Typography, Row, Col, Card, Space, Button, Divider, Spin, message } from 'antd';
import { ReadOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { miscAPI } from '../api'; // 引入 miscAPI
import PopularTagsCategories from '../components/PopularTagsCategories'; // 导入新组件
import './HomePage.css'; // 引入自定义 CSS

const { Title, Paragraph, Text } = Typography;

const HomePage = () => {
  const [homeData, setHomeData] = useState({ latest_posts: [], hot_posts: [], hot_categories: [], hot_tags: [] });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchHomePageData = async () => {
      try {
        setLoading(true);
        const response = await miscAPI.getHomePageData();
        if (response && response.code === 200 && response.data) {
          setHomeData(response.data);
        } else {
          console.warn('获取首页数据结构不正确或为空:', response);
          setHomeData({ latest_posts: [], hot_posts: [], hot_categories: [], hot_tags: [] });
        }
      } catch (error) {
        console.error('获取首页数据失败:', error);
        message.error('加载首页数据失败');
        setHomeData({ latest_posts: [], hot_posts: [], hot_categories: [], hot_tags: [] });
      } finally {
        setLoading(false);
      }
    };

    fetchHomePageData();
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

      {/* Main Content Area */}
      <Row justify="center" gutter={[24, 24]} style={{ padding: '0 24px' }}>
        {/* Latest Posts Section */}
        <Col xs={24} md={16} lg={18} className="latest-posts-section">
          <Title level={2} style={{ marginBottom: '2rem' }}>最新文章</Title>
          <Spin spinning={loading}>
            {homeData.latest_posts.length > 0 ? (
              <Row gutter={[24, 24]}>
                {homeData.latest_posts.map(post => (
                  <Col xs={24} sm={12} xl={8} key={post.id}> {/* Adjusted grid for better layout */}
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

        {/* Sidebar Section */}
        <Col xs={24} md={8} lg={6}>
          <PopularTagsCategories /> {/* 在侧边栏渲染热门标签和分类 */}
          {/* 可以添加其他侧边栏内容，例如热门文章 */}
        </Col>
      </Row>
    </div>
  );
};

export default HomePage;