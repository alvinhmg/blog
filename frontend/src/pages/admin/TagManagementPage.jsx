import React, { useState, useEffect } from 'react';
import { message, Spin, Input, Button, List, Popconfirm, Card, Typography, Space } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { tagAPI } from '../../api';

const { Title } = Typography;

const TagManagementPage = () => {
  const [tags, setTags] = useState([]);
  const [newTag, setNewTag] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTags = async () => {
      setLoading(true);
      try {
        const response = await tagAPI.getAllTags();
        // 确保tags是一个数组
        const tagsArray = Array.isArray(response) ? response : (response.data || []);
        setTags(tagsArray);
      } catch (error) {
        console.error('获取标签列表失败:', error);
        message.error('获取标签列表失败，请稍后重试');
      } finally {
        setLoading(false);
      }
    };
    
    fetchTags();
  }, []);

  const handleAddTag = async () => {
    if (!newTag.trim()) {
      message.warning('标签名称不能为空');
      return;
    }
    
    try {
      const response = await tagAPI.createTag({ name: newTag });
      if (response.success || response.code === 200) {
        const newTagData = response.tag || response.data;
        setTags(prev => Array.isArray(prev) ? [...prev, newTagData] : [newTagData]);
        setNewTag('');
        message.success('添加标签成功');
      }
    } catch (error) {
      console.error('添加标签失败:', error);
      message.error('添加标签失败，请稍后重试');
    }
  };

  const handleDeleteTag = async (tagId) => {
    try {
      const response = await tagAPI.deleteTag(tagId);
      if (response.success || response.code === 200) {
        setTags(prev => Array.isArray(prev) ? prev.filter(tag => tag.id !== tagId) : []);
        message.success('删除标签成功');
      }
    } catch (error) {
      console.error('删除标签失败:', error);
      message.error('删除标签失败，请稍后重试');
    }
  };

  if (loading) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <Spin size="large" tip="加载标签中..." />
      </div>
    );
  }

  return (
    <div style={{ padding: '20px' }}>
      <Card>
        <Title level={2}>标签管理</Title>
        <Space style={{ marginBottom: 16 }}>
          <Input
            placeholder="新标签名称"
            value={newTag}
            onChange={(e) => setNewTag(e.target.value)}
            onPressEnter={handleAddTag}
            style={{ width: 200 }}
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddTag}>
            添加标签
          </Button>
        </Space>
        <List
          bordered
          dataSource={tags}
          renderItem={item => (
            <List.Item
              actions={[
                <Popconfirm
                  title="确定要删除这个标签吗？"
                  onConfirm={() => handleDeleteTag(item.id)}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
                </Popconfirm>
              ]}
            >
              {item.name}
            </List.Item>
          )}
        />
      </Card>
    </div>
  );
};

export default TagManagementPage;