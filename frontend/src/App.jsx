import { useEffect } from 'react'; // Import useEffect
import { Outlet } from 'react-router-dom';
import { Layout, ConfigProvider, theme } from 'antd';
import './App.css';
import MainHeader from './layouts/MainHeader';
import MainFooter from './layouts/MainFooter';
import { useDispatch } from 'react-redux'; // Import useDispatch
import { setCredentials } from './store/authSlice'; // Import setCredentials
import { authAPI } from './api'; // Import authAPI

function App() {
  const dispatch = useDispatch(); // Get dispatch function

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        try {
          // Optionally verify token with backend or just fetch user data
          const response = await authAPI.getCurrentUser();
          if (response.code === 200 && response.data) {
            dispatch(setCredentials({ user: response.data, token }));
          } else {
            // Token might be invalid or expired
            localStorage.removeItem('token');
          }
        } catch (error) {
          console.error('Failed to fetch current user:', error);
          localStorage.removeItem('token'); // Clear invalid token
        }
      }
    };
    checkAuth();
  }, [dispatch]);

  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1890ff',
        },
      }}
    >
      <Layout className="app-container">
        <MainHeader />
        <Layout.Content className="main-content">
          <Outlet />
        </Layout.Content>
        <MainFooter />
      </Layout>
    </ConfigProvider>
  )
}

export default App
