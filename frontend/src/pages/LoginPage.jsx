import { useState } from 'react';
import { Form, Input, Button, Card, Typography, message, Divider } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../store/AuthContext'; // Import useAuth

const { Title, Paragraph } = Typography;

const LoginPage = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const auth = useAuth(); // Get auth context

  const onFinish = async (values) => {
    try {
      setLoading(true);
      // Call the actual login API
      const response = await fetch('/api/auth/login', { // Assuming backend runs on the same origin or proxy is set up
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || '登录失败');
      }

      const userData = await response.json(); // Assuming API returns user data (e.g., { user: {...}, token: '...' })
      
      // Call auth.login to update global state and store user info
      auth.login(userData.user); // Store user details from response

      message.success('登录成功！');
      navigate('/'); // Redirect to home page after successful login

    } catch (error) {
      console.error('Login failed:', error);
      message.error(error.message || '登录失败，请检查用户名和密码');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '40px auto', padding: '0 16px' }}>
      <Card bordered={false} style={{ boxShadow: '0 1px 2px rgba(0,0,0,0.1)' }}>
        <Title level={2} style={{ textAlign: 'center' }}>用户登录</Title>
        <Form
          name="login"
          initialValues={{ remember: true }}
          onFinish={onFinish}
          layout="vertical"
        >
          <Form.Item
            name="username" // Changed from email/username to just username based on API doc
            rules={[{ required: true, message: '请输入用户名!' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码!' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
              size="large"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block size="large">
              登录
            </Button>
          </Form.Item>
          
          <Divider plain>或者</Divider>
          
          <Paragraph style={{ textAlign: 'center' }}>
            还没有账号？ <Link to="/register">立即注册</Link>
          </Paragraph>
        </Form>
      </Card>
    </div>
  );
};

export default LoginPage;