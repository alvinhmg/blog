import { Layout } from 'antd';

const { Footer } = Layout;

const MainFooter = () => {
  return (
    <Footer style={{ textAlign: 'center' }}>
      Alvin的博客 ©{new Date().getFullYear()} Created by Alvin
    </Footer>
  );
};

export default MainFooter;