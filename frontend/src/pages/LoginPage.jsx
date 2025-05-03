import { useState } from 'react';
import { Form, Input, Button, Card, Typography, message, Divider } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';

const { Title, Paragraph } = Typography;

const LoginPage = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    try {
      setLoading(true);
      // 这里将来会实现实际的登录API调用
      console.log('登录信息:', values);
      
      // 模拟登录成功
      setTimeout(() => {
        message.success('登录成功！');
        navigate('/');
        setLoading(false);
      }, 1000);
    } catch (error) {
      message.error('登录失败，请检查用户名和密码');
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
            name="username"
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