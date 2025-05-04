import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Layout, Menu, Button, Dropdown } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { logout } from '../store/authSlice';

const { Header: AntHeader } = Layout;

const Header = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { user, isAuthenticated } = useSelector((state) => state.auth);

  // Removed useEffect hook that fetched user info

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  const userMenu = [
    {
      key: 'profile',
      label: '个人资料',
      onClick: () => navigate('/profile'), // Assuming a profile page exists or will be created
    },
    {
      key: 'logout',
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return (
    <AntHeader style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 20px' }}>
      <div className="logo" style={{ color: '#fff', fontSize: '18px' }}>
        Alvin的博客
      </div>
      <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['1']} style={{ flex: 1, minWidth: 0 }}>
        <Menu.Item key="1" onClick={() => navigate('/')}>首页</Menu.Item>
        <Menu.Item key="2" onClick={() => navigate('/posts')}>文章</Menu.Item>
      </Menu>
      <div style={{ marginLeft: 'auto' }}>
        {isAuthenticated && user ? (
          <Dropdown
            menu={{ items: userMenu }}
            placement="bottomRight"
          >
            <Button type="text" style={{ color: '#fff' }}>
              <UserOutlined />
              {/* Display nickname if available, otherwise username */}
              {user.nickname || user.username}
            </Button>
          </Dropdown>
        ) : (
          <Button type="primary" onClick={() => navigate('/login')}>
            登录
          </Button>
        )}
      </div>
    </AntHeader>
  );
};

export default Header;