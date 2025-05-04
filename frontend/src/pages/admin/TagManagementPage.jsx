import React, { useState, useEffect } from 'react';

const TagManagementPage = () => {
  const [tags, setTags] = useState([]);
  const [newTag, setNewTag] = useState('');

  useEffect(() => {
    fetch('/api/tags')
      .then(response => response.json())
      .then(data => setTags(data));
  }, []);

  const handleAddTag = () => {
    fetch('/api/tags', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name: newTag }),
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          setTags([...tags, data.tag]);
          setNewTag('');
        }
      });
  };

  const handleDeleteTag = (tagId) => {
    fetch(`/api/tags/${tagId}`, {
      method: 'DELETE',
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          setTags(tags.filter(tag => tag.id !== tagId));
        }
      });
  };

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
        {tags.map(tag => (
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