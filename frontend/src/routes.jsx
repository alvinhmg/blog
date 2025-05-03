import { createBrowserRouter } from 'react-router-dom';
import App from './App';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import PostListPage from './pages/PostListPage';
import PostDetailPage from './pages/PostDetailPage';
import UserProfilePage from './pages/UserProfilePage';
import AdminDashboard from './pages/admin/AdminDashboard';
import PostEditor from './pages/admin/PostEditor';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'login', element: <LoginPage /> },
      { path: 'register', element: <RegisterPage /> },
      { path: 'posts', element: <PostListPage /> },
      { path: 'posts/:id', element: <PostDetailPage /> },
      { path: 'profile', element: <UserProfilePage /> },
      { path: 'admin', element: <AdminDashboard /> },
      { path: 'admin/posts/edit/:id', element: <PostEditor /> },
      { path: 'admin/posts/new', element: <PostEditor /> },
    ],
  },
]);

export default router;