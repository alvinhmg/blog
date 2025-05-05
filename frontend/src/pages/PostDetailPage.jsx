import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Typography, Card, Tag, Space, Divider, Spin, Button, Form, Input, List, Avatar, message } from 'antd';
import { ArrowLeftOutlined, CalendarOutlined, UserOutlined, CommentOutlined, SendOutlined, LikeOutlined, LikeFilled } from '@ant-design/icons'; // Import Like icons
import ReactMarkdown from 'react-markdown';
import { postAPI } from '../api';

const { Title, Paragraph } = Typography;

const PostDetailPage = () => {
  const { id } = useParams();
  const [loading, setLoading] = useState(true);
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [submitting, setSubmitting] = useState(false);
  const [liked, setLiked] = useState(false); // State for like status
  const [likeCount, setLikeCount] = useState(0); // State for like count
  const [liking, setLiking] = useState(false); // State for like action loading
  const [form] = Form.useForm();

  useEffect(() => {
    // 获取文章详情
    const fetchPostDetail = async () => {
      setLoading(true);
      try {
        // 调用实际的API
        const response = await postAPI.getPostById(id);
        if (response.code === 200 && response.data) {
          setPost(response.data);
          setComments(response.data.comments || []);
          setLikeCount(response.data.like_count || 0); // Assuming backend provides like_count
          // TODO: Check if current user has liked this post (requires user info and backend logic)
          // setLiked(response.data.user_has_liked || false);
        } else {
          message.error('获取文章详情失败');
        }
        setLoading(false);
      } catch (error) {
        console.error('获取文章详情失败:', error);
        message.error('获取文章详情失败，请稍后再试');
        setLoading(false);
      }
    };

    fetchPostDetail();
  }, [id]);
  
  // 提交评论
  const handleSubmitComment = async (values) => {
    if (!values.content || values.content.trim() === '') {
      message.warning('评论内容不能为空');
      return;
    }
    
    setSubmitting(true);
    try {
      const response = await postAPI.addComment(id, values.content);
      if (response.code === 200 && response.data) {
        // 添加新评论到列表
        setComments([...comments, response.data]);
        message.success('评论发布成功');
        form.resetFields();
      } else {
        message.error(response.message || '发布评论失败');
      }
    } catch (error) {
      console.error('发布评论失败:', error);
      message.error('发布评论失败，请稍后再试');
    } finally {
      setSubmitting(false);
    }
  };

  // 处理点赞
  const handleLike = async () => {
    setLiking(true);
    try {
      // TODO: Implement actual API call for liking/unliking
      // const response = await postAPI.toggleLikePost(id);
      // if (response.code === 200) {
      //   setLiked(!liked);
      //   setLikeCount(liked ? likeCount - 1 : likeCount + 1);
      //   message.success(liked ? '取消点赞成功' : '点赞成功');
      // } else {
      //   message.error('操作失败');
      // }

      // Placeholder logic:
      await new Promise(resolve => setTimeout(resolve, 500)); // Simulate API call
      setLiked(!liked);
      setLikeCount(liked ? likeCount - 1 : likeCount + 1);
      message.success(liked ? '取消点赞成功 (模拟)' : '点赞成功 (模拟)');

    } catch (error) {
      console.error('点赞/取消点赞失败:', error);
      message.error('操作失败，请稍后再试');
    } finally {
      setLiking(false);
    }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '50px' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (!post) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Title level={3}>文章不存在或已被删除</Title>
        <Button type="primary">
          <Link to="/posts"><ArrowLeftOutlined /> 返回文章列表</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="post-detail-container" style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Button type="link" style={{ marginBottom: '16px', padding: 0 }}>
        <Link to="/posts"><ArrowLeftOutlined /> 返回文章列表</Link>
      </Button>
      
      <Card bordered={false}>
        <Title>{post.title}</Title>
        
        <Space split="|" style={{ marginBottom: '16px', color: 'rgba(0, 0, 0, 0.45)' }}>
          <span><CalendarOutlined style={{ marginRight: '4px' }} />{new Date(post.created_at).toLocaleDateString('zh-CN')}</span>
          <span><UserOutlined style={{ marginRight: '4px' }} />{post.author?.username || '未知作者'}</span>
        </Space>
        
        <Space wrap style={{ marginBottom: '24px' }}>
          {post.tags && post.tags.map(tag => (
            <Tag key={tag.id || tag.name} color="blue">
              {tag.name}
            </Tag>
          ))}
        </Space>
        
        <Divider />
        
        <div className="markdown-content" style={{ marginBottom: '24px' }}>
          <ReactMarkdown>{post.content}</ReactMarkdown>
        </div>

        {/* Like Button */}
        <div style={{ marginBottom: '24px', textAlign: 'center' }}>
          <Button 
            type={liked ? "primary" : "default"} 
            icon={liked ? <LikeFilled /> : <LikeOutlined />} 
            onClick={handleLike}
            loading={liking}
            size="large"
          >
            {liked ? '已赞' : '点赞'} ({likeCount})
          </Button>
        </div>
        
        <Divider orientation="left">
          <Space>
            <CommentOutlined />
            评论 ({comments.length})
          </Space>
        </Divider>
        
        {/* 评论列表 */}
        <List
          itemLayout="horizontal"
          dataSource={comments}
          renderItem={comment => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar src={comment.user?.avatar || 'https://joeschmoe.io/api/v1/random'} />}
                title={
                  <Space split="|">
                    <span>{comment.user?.username || '匿名用户'}</span>
                    <span style={{ fontSize: '12px', color: '#8c8c8c' }}>
                      {new Date(comment.created_at).toLocaleString('zh-CN')}
                    </span>
                  </Space>
                }
                description={comment.content}
              />
            </List.Item>
          )}
        />
        
        {/* 评论表单 */}
        <div style={{ marginTop: '24px' }}>
          <Form
            form={form}
            onFinish={handleSubmitComment}
            layout="vertical"
          >
            <Form.Item
              name="content"
              rules={[{ required: true, message: '请输入评论内容' }]}
            >
              <Input.TextArea 
                rows={4} 
                placeholder="写下你的评论..."
                disabled={submitting}
              />
            </Form.Item>
            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                icon={<SendOutlined />}
                loading={submitting}
              >
                发布评论
              </Button>
            </Form.Item>
          </Form>
        </div>
      </Card>
    </div>
  );
};

export default PostDetailPage;