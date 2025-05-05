import React, { useState, useEffect } from 'react';
import { List, Typography, Spin, message, Divider, Empty } from 'antd';
import { Link } from 'react-router-dom';
import { miscAPI } from '../api'; // 引入 miscAPI
import dayjs from 'dayjs'; // Use dayjs for date manipulation
import 'dayjs/locale/zh-cn'; // Import Chinese locale for dayjs
dayjs.locale('zh-cn'); // Set locale globally

const { Title, Text } = Typography;

const ArchivePage = () => {
  const [loading, setLoading] = useState(true);
  const [archiveData, setArchiveData] = useState([]); // 存储后端返回的原始归档数据
  const [postsByYearMonthDisplay, setPostsByYearMonthDisplay] = useState({}); // 用于显示的按年月分组数据

  useEffect(() => {
    const fetchAllPosts = async () => {
      setLoading(true);
      try {
        // Fetch archive data from the backend
        const response = await miscAPI.getArchiveData();

        if (response && response.code === 200 && Array.isArray(response.data)) {
          setArchiveData(response.data); // Store the raw data

          // Process data for display (group by 'YYYY年MM月')
          const groupedForDisplay = response.data.reduce((acc, item) => {
            // Convert 'YYYY-MM' to 'YYYY年MM月' for display title
            const displayYearMonth = dayjs(item.year_month + '-01').format('YYYY年MM月'); 
            acc[displayYearMonth] = item.posts.map(post => ({
              ...post,
              // Ensure created_at is parsed correctly if needed for sorting within month (though backend already sorts)
              created_at: post.created_at 
            }));
            return acc;
          }, {});
          setPostsByYearMonthDisplay(groupedForDisplay);

        } else {
          console.warn('获取归档数据结构不正确或为空:', response);
          setArchiveData([]);
          setPostsByYearMonthDisplay({});
        }
      } catch (error) {
        console.error('获取归档数据失败:', error);
        message.error('加载归档数据失败');
        setArchiveData([]);
        setPostsByYearMonthDisplay({});
      } finally {
        setLoading(false);
      }
    };

    fetchAllPosts();
  }, []);

  return (
    <div style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Title level={2}>文章归档</Title>
      <Divider />
      <Spin spinning={loading}>
        {Object.keys(postsByYearMonthDisplay).length > 0 ? (
          Object.entries(postsByYearMonthDisplay).map(([yearMonth, posts]) => (
            <div key={yearMonth} style={{ marginBottom: '24px' }}>
              <Title level={4}>{yearMonth}</Title>
              <List
                itemLayout="horizontal"
                dataSource={posts}
                renderItem={post => (
                  <List.Item>
                    <List.Item.Meta
                      title={<Link to={`/posts/${post.id}`}>{post.title}</Link>}
                      description={<Text type="secondary">{dayjs(post.created_at).format('MM月DD日')}</Text>}
                    />
                  </List.Item>
                )}
              />
            </div>
          ))
        ) : (
          !loading && Object.keys(postsByYearMonthDisplay).length === 0 && <Empty description="暂无文章归档" />
        )}
      </Spin>
    </div>
  );
};

export default ArchivePage;