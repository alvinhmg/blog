import { useState, useEffect } from 'react';
import { List, Card, Typography, Tag, Space, Input, Spin, Empty, Pagination } from 'antd';
import { SearchOutlined, CalendarOutlined, UserOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';

const { Title, Paragraph } = Typography;
const { Search } = Input;

const PostListPage = () => {
  const [loading, setLoading] = useState(false);
  const [posts, setPosts] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });

  // 模拟文章数据，将来会从API获取
  const mockPosts = [
    { 
      id: 1, 
      title: '欢迎来到我的博客', 
      summary: '这是我的第一篇博客文章，介绍了这个网站的功能和特点。', 
      content: '详细内容...', 
      author: 'Alvin',
      createdAt: '2023-05-01',
      tags: ['介绍', '欢迎']
    },
    { 
      id: 2, 
      title: 'React 18新特性解析', 
      summary: '深入探讨React 18带来的新特性和改进，以及如何在项目中应用。', 
      content: '详细内容...', 
      author: 'Alvin',
      createdAt: '2023-05-15',
      tags: ['React', '前端', '技术']
    },
    { 
      id: 3, 
      title: '前端性能优化技巧', 
      summary: '分享一些实用的前端性能优化方法，帮助你的网站加载更快、运行更流畅。', 
      content: '详细内容...', 
      author: 'Alvin',
      createdAt: '2023-06-01',
      tags: ['性能优化', '前端', '技术']
    },
  ];

  useEffect(() => {
    // 模拟API请求
    const fetchPosts = async () => {
      setLoading(true);
      try {
        // 这里将来会调用实际的API
        // const response = await postAPI.getAllPosts({
        //   page: pagination.current,
        //   pageSize: pagination.pageSize,
        //   search: searchText
        // });
        
        // 模拟API响应
        setTimeout(() => {
          setPosts(mockPosts);
          setPagination(prev => ({
            ...prev,
            total: mockPosts.length
          }));
          setLoading(false);
        }, 500);
      } catch (error) {
        console.error('获取文章列表失败:', error);
        setLoading(false);
      }
    };

    fetchPosts();
  }, [pagination.current, pagination.pageSize, searchText]);

  const handleSearch = (value) => {
    setSearchText(value);
    setPagination(prev => ({ ...prev, current: 1 }));
  };

  const handlePageChange = (page, pageSize) => {
    setPagination(prev => ({ ...prev, current: page, pageSize }));
  };

  return (
    <div style={{ padding: '24px' }}>
      <Title level={2}>博客文章</Title>
      
      <div style={{ marginBottom: '24px' }}>
        <Search
          placeholder="搜索文章"
          allowClear
          enterButton={<SearchOutlined />}
          size="large"
          onSearch={handleSearch}
          style={{ maxWidth: '500px' }}
        />
      </div>
      
      <Spin spinning={loading}>
        {posts.length > 0 ? (
          <List
            grid={{ gutter: 16, xs: 1, sm: 1, md: 2, lg: 3, xl: 3, xxl: 4 }}
            dataSource={posts}
            renderItem={(post) => (
              <List.Item>
                <Card 
                  hoverable
                  title={
                    <Link to={`/posts/${post.id}`}>
                      {post.title}
                    </Link>
                  }
                >
                  <Paragraph ellipsis={{ rows: 3 }}>
                    {post.summary}
                  </Paragraph>
                  
                  <Space direction="vertical" style={{ width: '100%' }}>
                    <Space wrap>
                      {post.tags.map(tag => (
                        <Tag key={tag} color="blue">{tag}</Tag>
                      ))}
                    </Space>
                    
                    <Space split="|">  
                      <span>
                        <CalendarOutlined /> {post.createdAt}
                      </span>
                      <span>
                        <UserOutlined /> {post.author}
                      </span>
                    </Space>
                  </Space>
                </Card>
              </List.Item>
            )}
          />
        ) : (
          <Empty description="暂无文章" />
        )}
        
        {posts.length > 0 && (
          <div style={{ textAlign: 'center', marginTop: '24px' }}>
            <Pagination
              current={pagination.current}
              pageSize={pagination.pageSize}
              total={pagination.total}
              onChange={handlePageChange}
              showSizeChanger
              showQuickJumper
            />
          </div>
        )}
      </Spin>
    </div>
  );
};

export default PostListPage;