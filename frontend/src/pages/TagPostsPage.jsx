import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { List, Typography, Spin, message, Card, Breadcrumb, Tag as AntTag } from 'antd';
import { HomeOutlined, TagOutlined } from '@ant-design/icons';
import { postAPI, tagAPI } from '../api';

const { Title, Text } = Typography;

const TagPostsPage = () => {
  const [posts, setPosts] = useState([]);
  const [tag, setTag] = useState(null);
  const [loading, setLoading] = useState(true);
  const { id: tagId } = useParams();

  useEffect(() => {
    const fetchTagAndPosts = async () => {
      if (!tagId) return;
      setLoading(true);
      try {
        const tagRes = await tagAPI.getTagById(tagId); // 假设有 getTagById API
        if (tagRes && (tagRes.data || tagRes.id)) {
          setTag(tagRes.data || tagRes);
        } else {
          message.error('未找到该标签');
          setTag(null);
        }

        const postsRes = await postAPI.getPostsByTag(tagId); // 假设有 getPostsByTag API
        if (postsRes && postsRes.data && Array.isArray(postsRes.data.posts)) {
          setPosts(postsRes.data.posts);
        } else if (postsRes && Array.isArray(postsRes)) { // Fallback for direct array response, though less likely for paginated data
          setPosts(postsRes);
        } else {
          setPosts([]);
        }
      } catch (error) {
        console.error('获取标签文章失败:', error);
        message.error('获取标签文章失败');
        setPosts([]);
        setTag(null);
      } finally {
        setLoading(false);
      }
    };

    fetchTagAndPosts();
  }, [tagId]);

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
            <TagOutlined />
            <span>标签</span>
          </Breadcrumb.Item>
          {tag && <Breadcrumb.Item>{tag.name}</Breadcrumb.Item>}
        </Breadcrumb>

        {tag ? (
          <Title level={2}>标签: <AntTag color="blue">{tag.name}</AntTag></Title>
        ) : (
          <Title level={2}>标签文章</Title>
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
                  <Text type="secondary">发布于: {new Date(post.createdAt).toLocaleDateString()}</Text>,
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
          <Text>该标签下暂无文章。</Text>
        )}
      </Card>
    </div>
  );
};

export default TagPostsPage;