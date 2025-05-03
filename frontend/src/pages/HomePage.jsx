import { Typography, Row, Col, Card, Space, Button } from 'antd';
import { ReadOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';

const { Title, Paragraph } = Typography;

const HomePage = () => {
  // 这里将来会从API获取最新文章
  const latestPosts = [
    { id: 1, title: '欢迎来到我的博客', summary: '这是我的第一篇博客文章，介绍了这个网站的功能和特点。', createdAt: '2023-05-01' },
    { id: 2, title: 'React 18新特性解析', summary: '深入探讨React 18带来的新特性和改进，以及如何在项目中应用。', createdAt: '2023-05-15' },
    { id: 3, title: '前端性能优化技巧', summary: '分享一些实用的前端性能优化方法，帮助你的网站加载更快、运行更流畅。', createdAt: '2023-06-01' },
  ];

  return (
    <div className="home-container" style={{ padding: '2rem' }}>
      <Row gutter={[24, 24]}>
        <Col span={24}>
          <div className="hero-section" style={{ textAlign: 'center', marginBottom: '2rem' }}>
            <Title>欢迎来到Alvin的博客</Title>
            <Paragraph>
              这是一个分享技术、思考和创意的个人空间。
            </Paragraph>
            <Button type="primary" size="large" icon={<ReadOutlined />}>
              <Link to="/posts">浏览文章</Link>
            </Button>
          </div>
        </Col>

        <Col span={24}>
          <Title level={2}>最新文章</Title>
          <Row gutter={[16, 16]}>
            {latestPosts.map(post => (
              <Col xs={24} sm={12} md={8} key={post.id}>
                <Card 
                  hoverable 
                  title={post.title}
                  extra={<Link to={`/posts/${post.id}`}>阅读更多</Link>}
                >
                  <Paragraph ellipsis={{ rows: 3 }}>{post.summary}</Paragraph>
                  <Space>
                    <span>发布于: {post.createdAt}</span>
                  </Space>
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