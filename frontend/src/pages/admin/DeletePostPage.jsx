import React, { useState, useEffect } from 'react';

const DeletePostPage = () => {
  const [posts, setPosts] = useState([]);

  useEffect(() => {
    // Fetch posts from API
    fetch('/api/posts')
      .then(response => response.json())
      .then(data => setPosts(data));
  }, []);

  const handleDelete = (postId) => {
    fetch(`/api/posts/${postId}`, {
      method: 'DELETE',
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          setPosts(posts.filter(post => post.id !== postId));
        }
      });
  };

  return (
    <div>
      <h1>删除文章</h1>
      <ul>
        {posts.map(post => (
          <li key={post.id}>
            {post.title}
            <button onClick={() => handleDelete(post.id)}>删除</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default DeletePostPage;