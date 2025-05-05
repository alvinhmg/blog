import React, { useState, useEffect } from 'react';
import { message, Spin } from 'antd';
import { tagAPI } from '../../api';

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
    <div>
      <h1>标签管理</h1>
      <input
        type="text"
        value={newTag}
        onChange={(e) => setNewTag(e.target.value)}
        placeholder="新标签名称"
      />
      <button onClick={handleAddTag}>添加标签</button>
      <ul>
        {Array.isArray(tags) && tags.map(tag => (
          <li key={tag.id}>
            {tag.name}
            <button onClick={() => handleDeleteTag(tag.id)}>删除</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TagManagementPage;