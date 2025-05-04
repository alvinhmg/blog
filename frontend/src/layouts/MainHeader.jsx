import { Layout, Menu, Button, Space } from 'antd';
import './MainHeader.css'; // 导入自定义样式
import { UserOutlined, HomeOutlined, ReadOutlined, LoginOutlined, LogoutOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../store/AuthContext'; // Import useAuth

const { Header } = Layout;

const MainHeader = () => {
  const navigate = useNavigate();
  const { user, logout, loading } = useAuth(); // Get user, logout function, and loading state from context
  // const isLoggedIn = false; // Replaced by user state from context
  const isAdmin = user?.role === 'admin'; // Example: Check if user has admin role

  const handleLogout = () => {
    logout(); // Call logout from context
    console.log('用户登出');
    navigate('/login');
  };

  // Don't render header content until auth state is loaded
  if (loading) {
    return <Header className="header"></Header>; // Or a loading indicator
  }

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
        {user ? ( // Check if user object exists
          <>
            <span style={{ color: 'white', marginRight: '10px' }}>欢迎, {user.nickname || user.username}</span>
            {/* <Button type="link" icon={<UserOutlined />} onClick={() => navigate('/profile')}>
              个人中心
            </Button> */}
            <Button type="link" icon={<LogoutOutlined />} onClick={handleLogout}>
              登出
            </Button>
          </>
        ) : (
          <Button type="primary" icon={<LoginOutlined />} onClick={() => navigate('/login')}>
            登录
          </Button>
        )}
      </Space>
    </Header>
  );
};

export default MainHeader;