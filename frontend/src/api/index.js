import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api', // 后端API的基础URL
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器 - 添加认证token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response && error.response.status === 401) {
      // 未授权，清除token并重定向到登录页
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 认证相关API
export const authAPI = {
  login: (username, password) => api.post('/auth/login', { username, password }),
  register: (userData) => api.post('/auth/register', userData),
  getCurrentUser: () => api.get('/auth/user'),
  logout: () => {
    // 先调用后端登出接口，然后清除本地存储的token
    return api.post('/auth/logout')
      .then(response => {
        localStorage.removeItem('token');
        return response;
      })
      .catch(error => {
        // 即使后端请求失败，也清除本地token
        localStorage.removeItem('token');
        return Promise.reject(error);
      });
  },
};

// 文章相关API
export const postAPI = {
  getAllPosts: (params) => api.get('/posts', { params }),
  getPostById: (id) => api.get(`/posts/${id}`),
  createPost: (postData) => api.post('/posts', postData),
  updatePost: (id, postData) => api.put(`/posts/${id}`, postData),
  deletePost: (id) => api.delete(`/posts/${id}`),
  // 评论相关API
  addComment: (postId, content) => api.post(`/posts/${postId}/comments`, { content }),
  getComments: (postId) => api.get(`/posts/${postId}/comments`),
  deleteComment: (commentId) => api.delete(`/comments/${commentId}`),
};

export default api;