import { Layout, Menu, Button, Space } from 'antd';
import { UserOutlined, HomeOutlined, ReadOutlined, LoginOutlined, LogoutOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';

const { Header } = Layout;

const MainHeader = () => {
  const navigate = useNavigate();
  // 这里将来需要从Redux或Context中获取用户登录状态
  const isLoggedIn = false;
  const isAdmin = false;

  const handleLogout = () => {
    // 这里将来需要实现登出逻辑
    console.log('用户登出');
    navigate('/login');
  };

  return (
    <Header className="header">
      <div className="logo">
        <Link to="/">Alvin的博客</Link>
      </div>
      <Menu
        theme="dark"
        mode="horizontal"
        defaultSelectedKeys={['home']}
        style={{ flex: 1, minWidth: 0 }}
        items={[
          { key: 'home', icon: <HomeOutlined />, label: <Link to="/">首页</Link> },
          { key: 'posts', icon: <ReadOutlined />, label: <Link to="/posts">文章</Link> },
          isAdmin ? { key: 'admin', icon: <UserOutlined />, label: <Link to="/admin">管理</Link> } : null,
        ].filter(Boolean)}
      />
      <Space>
        {isLoggedIn ? (
          <>
            <Button type="link" icon={<UserOutlined />} onClick={() => navigate('/profile')}>
              个人中心
            </Button>
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