import { Layout, Menu, Button, Space } from 'antd';
import './MainHeader.css'; // 导入自定义样式
import { UserOutlined, HomeOutlined, ReadOutlined, LoginOutlined, LogoutOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux'; // Import useSelector and useDispatch
import { logout } from '../store/authSlice'; // Import logout action

const { Header } = Layout;

const MainHeader = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch(); // Get dispatch function
  const { user, isAuthenticated } = useSelector((state) => state.auth); // Get user and isAuthenticated from Redux store
  const isAdmin = user?.role === 'admin'; // Check if user has admin role

  const handleLogout = () => {
    dispatch(logout()); // Dispatch logout action
    console.log('用户登出');
    navigate('/login');
  };

  // No need for loading state from context anymore
  // if (loading) {
  //   return <Header className="header"></Header>; // Or a loading indicator
  // }

  return (
    <Header className="header">
      <div className="logo">
        <Link to="/">Alvin的博客</Link>
      </div>
      <Menu
        theme="dark"
        mode="horizontal"
        // Determine defaultSelectedKeys based on current path or keep it simple
        style={{ flex: 1, minWidth: 0 }}
        items={[
          { key: 'home', icon: <HomeOutlined />, label: <Link to="/">首页</Link> },
          { key: 'posts', icon: <ReadOutlined />, label: <Link to="/posts">文章</Link> },
          isAdmin ? { key: 'admin', icon: <UserOutlined />, label: <Link to="/admin">管理</Link> } : null,
        ].filter(Boolean)}
      />
      <Space>
        {isAuthenticated && user ? ( // Check isAuthenticated and user from Redux
          <>
            <span style={{ color: 'white', marginRight: '10px' }}>欢迎, {user.nickname || user.username}</span>
            <Button type="link" icon={<LogoutOutlined />} onClick={handleLogout}>
              登出
            </Button>
          </>
        ) : (
          <Button type="primary" icon={<LoginOutlined />} onClick={() => navigate('/login')}>
            登录/注册
          </Button>
        )}
      </Space>
    </Header>
  );
};

export default MainHeader;