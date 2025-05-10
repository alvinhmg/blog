import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Select, message, Typography, Card, Space, Divider, Spin } from 'antd';
import { useNavigate } from 'react-router-dom';
import MDEditor from '@uiw/react-md-editor';
import '@uiw/react-md-editor/markdown-editor.css';

import { postAPI, categoryAPI, tagAPI } from '../../api';

const { Title } = Typography;
const { TextArea } = Input;
const { Option } = Select;

const CreatePostPage = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState([]);
  const [tags, setTags] = useState([]);
  const [dataLoading, setDataLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      setDataLoading(true);
      try {
        const [categoriesRes, tagsRes] = await Promise.all([
          categoryAPI.getAllCategories(),
          tagAPI.getAllTags() // 假设有 tagAPI.getAllTags
        ]);
        
        // 确保返回的是数组
        setCategories(Array.isArray(categoriesRes) ? categoriesRes : (categoriesRes.data || []));
        setTags(Array.isArray(tagsRes) ? tagsRes : (tagsRes.data || []));
        
      } catch (error) {
        console.error('获取分类或标签失败:', error);
        message.error('获取分类或标签列表失败');
      } finally {
        setDataLoading(false);
      }
    };
    fetchData();
  }, []);

  // 这里应该从API获取分类和标签列表
  // const categories = [
  //   { id: 1, name: '技术' },
  //   { id: 2, name: '生活' },
  //   { id: 3, name: '思考' },
  // ];

  // const tags = [
  //   { id: 1, name: 'React' },
  //   { id: 2, name: 'JavaScript' },
  //   { id: 3, name: '前端' },
  //   { id: 4, name: '后端' },
  // ];

  const handleSubmit = async (values) => {
    try {
      setLoading(true);
      // 处理标签和分类的格式
      const postData = {
        title: values.title,
        content: values.content,
        excerpt: values.excerpt || values.content.substring(0, 150),
        status: 'published',
        categories: values.categories,
        tags: values.tags,
      };

      const response = await postAPI.createPost(postData);
      message.success('文章创建成功！');
      navigate('/admin');
    } catch (error) {
      console.error('创建文章失败:', error);
      message.error(error.response?.data?.message || '创建文章失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '20px' }}>
      <Card>
        <Title level={2}>创建新文章</Title>
        <Divider />
        {dataLoading ? (
          <div style={{ textAlign: 'center', padding: '50px' }}>
            <Spin tip="加载分类和标签..." />
          </div>
        ) : (
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{ status: 'published' }}
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
              <MDEditor
                height={400}
                data-color-mode="light" // Ensure light theme is applied
                previewOptions={{
                }}
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
                loading={dataLoading} // 添加加载状态
                disabled={dataLoading} // 禁用选择框直到数据加载完成
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
                loading={dataLoading} // 添加加载状态
                disabled={dataLoading} // 禁用选择框直到数据加载完成
              >
                {tags.map(tag => (
                  <Option key={tag.id} value={tag.id}>{tag.name}</Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit" loading={loading}>
                  发布文章
                </Button>
                <Button onClick={() => navigate('/admin')}>
                  取消
                </Button>
              </Space>
            </Form.Item>
          </Form>
        )}
      </Card>
    </div>
  );
};

export default CreatePostPage;