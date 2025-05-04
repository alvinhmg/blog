import React, { useState } from 'react';

const CreatePostPage = () => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [category, setCategory] = useState('');
  const [tags, setTags] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch('/api/posts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ title, content, category, tags }),
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          console.log('Post created successfully');
        }
      });
  };

  return (
    <div>
      <h1>创建文章</h1>
      <form onSubmit={handleSubmit}>
        <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="标题" required />
        <textarea value={content} onChange={(e) => setContent(e.target.value)} placeholder="内容" required />
        <input type="text" value={category} onChange={(e) => setCategory(e.target.value)} placeholder="分类" />
        <input type="text" value={tags} onChange={(e) => setTags(e.target.value)} placeholder="标签" />
        <button type="submit">创建</button>
      </form>
    </div>
  );
};

export default CreatePostPage;