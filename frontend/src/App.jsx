import { Outlet } from 'react-router-dom';
import { Layout, ConfigProvider, theme } from 'antd';
import './App.css';
import MainHeader from './layouts/MainHeader';
import MainFooter from './layouts/MainFooter';

function App() {
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
