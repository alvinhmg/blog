import React, { useState, useEffect } from 'react';

const UpdatePostPage = ({ match }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [category, setCategory] = useState('');
  const [tags, setTags] = useState('');

  useEffect(() => {
    fetch(`/api/posts/${match.params.id}`)
      .then(response => response.json())
      .then(data => {
        setTitle(data.title);
        setContent(data.content);
        setCategory(data.category);
        setTags(data.tags);
      });
  }, [match.params.id]);

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch(`/api/posts/${match.params.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ title, content, category, tags }),
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          console.log('Post updated successfully');
        }
      });
  };

  return (
    <div>
      <h1>更新文章</h1>
      <form onSubmit={handleSubmit}>
        <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="标题" required />
        <textarea value={content} onChange={(e) => setContent(e.target.value)} placeholder="内容" required />
        <input type="text" value={category} onChange={(e) => setCategory(e.target.value)} placeholder="分类" />
        <input type="text" value={tags} onChange={(e) => setTags(e.target.value)} placeholder="标签" />
        <button type="submit">更新</button>
      </form>
    </div>
  );
};

export default UpdatePostPage;