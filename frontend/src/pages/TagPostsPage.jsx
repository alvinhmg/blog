import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { List, Typography, Spin, message, Card, Breadcrumb, Tag as AntTag, Avatar } from 'antd';
import { HomeOutlined, TagOutlined, CalendarOutlined, EyeOutlined, LikeOutlined, MessageOutlined } from '@ant-design/icons';
import { postAPI, tagAPI } from '../api';
import './TagPostsPage.css'; // 引入自定义CSS文件

const { Title, Text, Paragraph } = Typography;

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
      <div className="loading-container">
        <Spin size="large" tip="精彩内容加载中..." />
      </div>
    );
  }

  return (
    <div className="tag-posts-container">
      <Card className="main-card">
        <Breadcrumb className="custom-breadcrumb">
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
          <div className="tag-header">
            <Title level={2}>
              <TagOutlined className="tag-icon" />
              <span>标签:</span> 
              <AntTag className="custom-tag" color="#1890ff">{tag.name}</AntTag>
            </Title>
            <Paragraph className="tag-description">
              {tag.description || `与 "${tag.name}" 相关的所有文章`}
            </Paragraph>
          </div>
        ) : (
          <Title level={2} className="page-title">标签文章</Title>
        )}

        {posts.length > 0 ? (
          <List
            className="posts-list"
            itemLayout="vertical"
            size="large"
            dataSource={posts}
            renderItem={post => (
              <List.Item
                key={post.id}
                className="post-item"
                actions={[
                  <Text type="secondary" key={`action-date-${post.id}`} className="post-meta-item">
                    <CalendarOutlined className="post-meta-icon" />
                    {post.created_at ? new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) : '未知日期'}
                  </Text>,
                  <Text type="secondary" key={`action-views-${post.id}`} className="post-meta-item">
                    <EyeOutlined className="post-meta-icon" />
                    {post.views || 0} 阅读
                  </Text>,
                  <Text type="secondary" key={`action-likes-${post.id}`} className="post-meta-item">
                    <LikeOutlined className="post-meta-icon" />
                    {post.likes || 0} 喜欢
                  </Text>,
                  <Text type="secondary" key={`action-comments-${post.id}`} className="post-meta-item">
                    <MessageOutlined className="post-meta-icon" />
                    {post.comments_count || 0} 评论
                  </Text>,
                ]}
                extra={
                  post.cover_image_url ? (
                    <div className="post-cover-container">
                      <img
                        alt={post.title || 'post cover'}
                        src={post.cover_image_url}
                        className="post-cover-image"
                      />
                    </div>
                  ) : null
                }
              >
                <List.Item.Meta
                  avatar={post.author && post.author.avatar_url && <Avatar src={post.author.avatar_url} />}
                  title={<Link to={`/posts/${post.id}`} className="post-title">{post.title}</Link>}
                  description={
                    <div className="post-description">
                      <Text type="secondary" className="post-excerpt">
                        {post.excerpt || `${post.content.substring(0, 150)}...`}
                      </Text>
                      {post.tags && post.tags.length > 0 && (
                        <div className="post-tags">
                          {post.tags.map(tag => (
                            <Link to={`/tags/${tag.id}`} key={tag.id}>
                              <AntTag className="post-tag" color="#e6f7ff">{tag.name}</AntTag>
                            </Link>
                          ))}
                        </div>
                      )}
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <div className="empty-state">
            <TagOutlined className="empty-icon" />
            <Text>该标签下暂无文章。</Text>
          </div>
        )}
      </Card>
    </div>
  );
};

export default TagPostsPage;