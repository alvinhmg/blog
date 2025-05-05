import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Select, message, Typography, Card, Space, Divider, Spin } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { postAPI } from '../../api';

const { Title } = Typography;
const { TextArea } = Input;
const { Option } = Select;

const UpdatePostPage = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [fetchLoading, setFetchLoading] = useState(true);
  const navigate = useNavigate();
  const { id } = useParams();

  // 这里应该从API获取分类和标签列表
  const categories = [
    { id: 1, name: '技术' },
    { id: 2, name: '生活' },
    { id: 3, name: '思考' },
  ];

  const tags = [
    { id: 1, name: 'React' },
    { id: 2, name: 'JavaScript' },
    { id: 3, name: '前端' },
    { id: 4, name: '后端' },
  ];

  useEffect(() => {
    const fetchPost = async () => {
      try {
        setFetchLoading(true);
        const post = await postAPI.getPostById(id);
        
        // 确保post数据存在且格式正确
        if (!post || typeof post !== 'object') {
          throw new Error('获取文章数据格式不正确');
        }
        
        // 处理可能的嵌套数据结构
        const postData = post.data || post;
        
        // 设置表单初始值
        form.setFieldsValue({
          title: postData.title || '',
          content: postData.content || '',
          excerpt: postData.excerpt || '',
          categories: Array.isArray(postData.categories) ? postData.categories.map(c => c.id) : [],
          tags: Array.isArray(postData.tags) ? postData.tags.map(t => t.id) : []
        });
      } catch (error) {
        console.error('获取文章失败:', error);
        message.error('获取文章失败，请稍后重试');
        navigate('/admin');
      } finally {
        setFetchLoading(false);
      }
    };

    if (id) {
      fetchPost();
    }
  }, [id, form, navigate]);

  const handleSubmit = async (values) => {
    try {
      setLoading(true);
      // 处理标签和分类的格式
      const postData = {
        title: values.title,
        content: values.content,
        excerpt: values.excerpt || values.content.substring(0, 150),
        status: 'published',
        categories: values.categories || [],
        tags: values.tags || [],
      };

      // 添加错误处理和重试逻辑
      try {
        // 确保使用正确的API路径和参数
        const response = await postAPI.updatePost(id, postData);
        if (response && (response.code === 200 || response.success)) {
          message.success('文章更新成功！');
          navigate('/admin/delete-post'); // 导航到文章管理页面
        } else {
          throw new Error('更新文章失败');
        }
      } catch (apiError) {
        console.error('API调用失败，尝试直接发送请求:', apiError);
        // 如果API调用失败，尝试直接使用fetch
        try {
          const fetchResponse = await fetch(`http://localhost:8080/api/posts/${id}`, {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify(postData)
          });

          if (!fetchResponse.ok) {
            const errorData = await fetchResponse.text(); // 尝试获取错误详情
            throw new Error(`HTTP error! status: ${fetchResponse.status}, message: ${errorData}`);
          }

          // 假设 fetch 成功意味着更新成功
          message.success('文章更新成功！(通过备用请求)');
          navigate('/admin/delete-post'); // 导航到文章管理页面
        } catch (fetchError) {
          console.error('直接发送请求失败:', fetchError);
          // 显示具体的 fetch 错误信息
          message.error(`通过备用请求更新文章也失败了: ${fetchError.message}`);
          // 此处不再向上抛出错误，避免被外层 catch 再次处理，除非有特定需求
        }
      }
    } catch (error) {
      console.error('更新文章失败:', error);
      message.error(error.response?.data?.message || error.message || '更新文章失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  if (fetchLoading) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <Spin size="large" tip="加载文章中..." />
      </div>
    );
  }

  return (
    <div style={{ padding: '20px' }}>
      <Card>
        <Title level={2}>编辑文章</Title>
        <Divider />
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="title"
            label="文章标题"
            rules={[{ required: true, message: '请输入文章标题' }]}
          >
            <Input placeholder="请输入文章标题" />
          </Form.Item>

          <Form.Item
            name="excerpt"
            label="文章摘要"
            extra="如不填写，将自动截取正文前150个字符"
          >
            <TextArea 
              placeholder="请输入文章摘要" 
              autoSize={{ minRows: 2, maxRows: 4 }}
            />
          </Form.Item>

          <Form.Item
            name="content"
            label="文章内容"
            rules={[{ required: true, message: '请输入文章内容' }]}
          >
            <TextArea 
              placeholder="请输入文章内容，支持Markdown格式" 
              autoSize={{ minRows: 10, maxRows: 20 }}
            />
          </Form.Item>

          <Form.Item
            name="categories"
            label="文章分类"
          >
            <Select
              mode="multiple"
              placeholder="请选择文章分类"
              style={{ width: '100%' }}
            >
              {categories.map(category => (
                <Option key={category.id} value={category.id}>{category.name}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="tags"
            label="文章标签"
          >
            <Select
              mode="multiple"
              placeholder="请选择文章标签"
              style={{ width: '100%' }}
            >
              {tags.map(tag => (
                <Option key={tag.id} value={tag.id}>{tag.name}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={loading}>
                更新文章
              </Button>
              <Button onClick={() => navigate('/admin')}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default UpdatePostPage;