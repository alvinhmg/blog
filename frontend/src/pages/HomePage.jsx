import { Typography, Row, Col, Card, Space, Button, Divider } from 'antd';
import { ReadOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import './HomePage.css'; // 引入自定义 CSS

const { Title, Paragraph, Text } = Typography;

const HomePage = () => {
  // 这里将来会从API获取最新文章
  const latestPosts = [
    { id: 1, title: '欢迎来到我的博客', summary: '这是我的第一篇博客文章，介绍了这个网站的功能和特点。', createdAt: '2023-05-01' },
    { id: 2, title: 'React 18新特性解析', summary: '深入探讨React 18带来的新特性和改进，以及如何在项目中应用。', createdAt: '2023-05-15' },
    { id: 3, title: '前端性能优化技巧', summary: '分享一些实用的前端性能优化方法，帮助你的网站加载更快、运行更流畅。', createdAt: '2023-06-01' },
  ];

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
          <Row gutter={[24, 24]}>
            {latestPosts.map(post => (
              <Col xs={24} sm={12} md={8} key={post.id}>
                <Card
                  hoverable
                  className="post-card"
                  title={<Link to={`/posts/${post.id}`} className="post-card-title">{post.title}</Link>}
                  extra={<Link to={`/posts/${post.id}`}>阅读更多</Link>}
                >
                  <Paragraph ellipsis={{ rows: 3 }} className="post-card-summary">{post.summary}</Paragraph>
                  <Text type="secondary">发布于: {post.createdAt}</Text>
                </Card>
              </Col>
            ))}
          </Row>
        </Col>
      </Row>
    </div>
  );
};

export default HomePage;