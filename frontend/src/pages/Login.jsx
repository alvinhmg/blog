import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { setCredentials } from '../store/authSlice';
import { Button, Form, Input, message, Card, Typography, Divider } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { authAPI } from '../api';
import './Login.css';

const { Title, Paragraph } = Typography;

const Login = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values) => {
    try {
      setLoading(true);
      const response = await authAPI.login(values.username, values.password);
      
      if (response.code === 200) {
        const { user, token } = response.data;
        localStorage.setItem('token', token);
        dispatch(setCredentials({ user, token }));
        message.success('登录成功');
        navigate('/');
      } else {
        message.error(response.message || '登录失败');
      }
    } catch (error) {
      message.error(error.response?.data?.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div className="login-page-background"></div>
      <div className="login-container">
        <Card bordered={false} className="login-card">
          <Title level={2} className="login-title">用户登录</Title>
          <Form
            name="login"
            onFinish={onFinish}
            layout="vertical"
            autoComplete="off"
            className="login-form"
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
            <Button 
              type="primary" 
              htmlType="submit" 
              loading={loading} 
              block 
              size="large"
              className="login-button"
            >
              登录
            </Button>
          </Form.Item>
          
          <Divider plain className="login-divider">或者</Divider>
          
          <Paragraph className="register-link">
            还没有账号？ <Link to="/register">立即注册</Link>
          </Paragraph>
        </Form>
      </Card>
    </div>
    </>

  );
};

export default Login;