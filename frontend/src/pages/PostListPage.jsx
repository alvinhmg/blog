import { useState, useEffect } from 'react';
import { List, Card, Typography, Tag, Space, Input, Spin, Empty, Pagination, message } from 'antd'; // 引入 message
import { SearchOutlined, CalendarOutlined, UserOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { postAPI } from '../api'; // 引入 postAPI

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

  // 移除 mockPosts
  // const mockPosts = [...];

  useEffect(() => {
    const fetchPosts = async () => {
      setLoading(true);
      try {
        // 根据是否有搜索文本决定调用哪个API
        let response;
        if (searchText) {
          response = await postAPI.searchPosts({
            q: searchText,
            page: pagination.current,
            page_size: pagination.pageSize // 后端 SearchPosts 使用 page_size
          });
        } else {
          response = await postAPI.getAllPosts({
            page: pagination.current,
            page_size: pagination.pageSize // 后端 GetPosts 使用 page_size
          });
        }

        if (response && response.data && Array.isArray(response.data.posts)) {
          setPosts(response.data.posts);
          setPagination(prev => ({
            ...prev,
            total: response.data.total || 0 // 使用后端返回的总数
          }));
        } else {
          console.warn('获取文章列表数据结构不正确或为空:', response);
          setPosts([]);
          setPagination(prev => ({ ...prev, total: 0 }));
        }
      } catch (error) {
        console.error('获取文章列表失败:', error);
        message.error('加载文章列表失败'); // 添加错误提示
        setPosts([]);
        setPagination(prev => ({ ...prev, total: 0 }));
      } finally {
        setLoading(false);
      }
    };

    fetchPosts();
  }, [pagination.current, pagination.pageSize, searchText]);

  const handleSearch = (value) => {
    setSearchText(value);
    setPagination(prev => ({ ...prev, current: 1 })); // 搜索时回到第一页
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
                    <Link to={`/posts/${post.id}`} style={{ display: 'block', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                      {post.title}
                    </Link>
                  }
                  style={{ height: '100%' }}
                  bodyStyle={{ height: 'calc(100% - 56px)', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}
                >
                  <div>
                    <Paragraph ellipsis={{ rows: 3 }} style={{ marginBottom: '16px' }}>
                      {post.excerpt || post.summary /* 优先使用 excerpt */} 
                    </Paragraph>
                  </div>

                  <div>
                    <Space wrap style={{ marginBottom: '8px' }}>
                      {/* 假设 post.tags 是存在的，如果不存在需要处理 */} 
                      {post.tags && post.tags.map(tag => (
                        <Tag key={tag.id || tag.name} color="blue">{tag.name}</Tag> // 使用 tag.name
                      ))}
                    </Space>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: 'rgba(0, 0, 0, 0.45)', fontSize: '12px' }}>
                      <span>
                        <CalendarOutlined style={{ marginRight: '4px' }} /> {new Date(post.created_at).toLocaleDateString('zh-CN')} {/* 格式化日期 */}
                      </span>
                      <span>
                        <UserOutlined style={{ marginRight: '4px' }} /> {post.author?.username || '未知作者'} {/* 处理 author 可能不存在的情况 */}
                      </span>
                    </div>
                  </div>
                </Card>
              </List.Item>
            )}
          />
        ) : (
          !loading && <Empty description="暂无文章" /> // 只有在非加载状态下显示 Empty
        )}

        {pagination.total > 0 && ( // 只有在有数据时显示分页
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