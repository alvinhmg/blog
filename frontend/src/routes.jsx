import { createBrowserRouter } from 'react-router-dom';
import App from './App';
import HomePage from './pages/HomePage';
import Login from './pages/Login';
import RegisterPage from './pages/RegisterPage';
import PostListPage from './pages/PostListPage';
import PostDetailPage from './pages/PostDetailPage';
import ArchivePage from './pages/ArchivePage'; // 导入归档页面组件
import MainPage from './pages/admin/MainPage';
import CreatePostPage from './pages/admin/CreatePostPage';
import UpdatePostPage from './pages/admin/UpdatePostPage';
import DeletePostPage from './pages/admin/DeletePostPage';
import CategoryManagementPage from './pages/admin/CategoryManagementPage';
import TagManagementPage from './pages/admin/TagManagementPage';
import CategoryPostsPage from './pages/CategoryPostsPage'; // 新增：分类文章列表页面
import TagPostsPage from './pages/TagPostsPage'; // 新增：标签文章列表页面

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'login', element: <Login /> },
      { path: 'register', element: <RegisterPage /> },
      { path: 'posts', element: <PostListPage /> },
      { path: 'posts/:id', element: <PostDetailPage /> },
      { path: 'archive', element: <ArchivePage /> }, // 添加归档页面路由
      { path: 'admin', element: <MainPage /> },
      { path: 'admin/create-post', element: <CreatePostPage /> },
      { path: 'admin/update-post/:id', element: <UpdatePostPage /> },
      { path: 'admin/delete-post', element: <DeletePostPage /> },
      { path: 'admin/category-management', element: <CategoryManagementPage /> },
      { path: 'admin/tag-management', element: <TagManagementPage /> },
      { path: 'categories/:id', element: <CategoryPostsPage /> }, // 新增：分类文章路由
      { path: 'tags/:id', element: <TagPostsPage /> }, // 新增：标签文章路由
    ],
  },
]);

export default router;