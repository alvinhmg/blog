import React, { useState, useEffect } from 'react';
import { Card, Tag, List, Typography, Spin, message } from 'antd';
import { Link } from 'react-router-dom';
import { categoryAPI, tagAPI, miscAPI } from '../api'; // Import miscAPI as well if needed elsewhere, or just categoryAPI and tagAPI

const { Title } = Typography;

const PopularTagsCategories = () => {
  const [popularTags, setPopularTags] = useState([]);
  const [popularCategories, setPopularCategories] = useState([]);
  const [loadingTags, setLoadingTags] = useState(true);
  const [loadingCategories, setLoadingCategories] = useState(true);

  useEffect(() => {
    const fetchPopularTags = async () => {
      try {
        setLoadingTags(true);
        // Use the actual API call for hot tags
        const response = await tagAPI.getHotTags({ limit: 10 }); 
        // Adjust response handling based on the actual API structure (assuming { code, message, data })
        if (response && response.code === 200 && Array.isArray(response.data)) {
          setPopularTags(response.data);
        } else if (response && Array.isArray(response)) { // Keep handling direct array for flexibility
           console.warn('Received direct array for tags, expected standard response format.');
           setPopularTags(response);
        } else {
          console.warn('获取热门标签数据结构不正确或为空:', response);
          setPopularTags([]);
        }
      } catch (error) {
        console.error('获取热门标签失败:', error);
        message.error('加载热门标签失败');
        setPopularTags([]);
      } finally {
        setLoadingTags(false);
      }
    };

    const fetchPopularCategories = async () => {
      try {
        setLoadingCategories(true);
        // Use the actual API call for hot categories
        const response = await categoryAPI.getHotCategories({ limit: 5 });
        // Adjust response handling based on the actual API structure (assuming { code, message, data })
        if (response && response.code === 200 && Array.isArray(response.data)) {
          setPopularCategories(response.data);
        } else if (response && Array.isArray(response)) { // Keep handling direct array for flexibility, though less likely now
           console.warn('Received direct array for categories, expected standard response format.');
           setPopularCategories(response);
        } else {
          console.warn('获取热门分类数据结构不正确或为空:', response);
          setPopularCategories([]);
        }
      } catch (error) {
        console.error('获取热门分类失败:', error);
        message.error('加载热门分类失败');
        setPopularCategories([]);
      } finally {
        setLoadingCategories(false);
      }
    };

    fetchPopularTags();
    fetchPopularCategories();
  }, []);

  return (
    <Card title="热门内容" style={{ marginBottom: '24px' }}>
      <Spin spinning={loadingCategories}>
        <Title level={5}>热门分类</Title>
        {popularCategories.length > 0 ? (
          <List
            size="small"
            dataSource={popularCategories}
            renderItem={category => (
              <List.Item>
                {/* TODO: Link to category page if exists */}
                <Link to={`/categories/${category.id}`}>{category.name}</Link>
              </List.Item>
            )}
          />
        ) : (
          !loadingCategories && <Typography.Text type="secondary">暂无热门分类</Typography.Text>
        )}
      </Spin>

      <Spin spinning={loadingTags} style={{ marginTop: '16px' }}>
        <Title level={5} style={{ marginTop: '16px' }}>热门标签</Title>
        {popularTags.length > 0 ? (
          <div>
            {popularTags.map(tag => (
              // TODO: Link to tag page if exists
              <Tag key={tag.id} color="blue" style={{ margin: '4px' }}>
                <Link to={`/tags/${tag.id}`} style={{ color: 'inherit' }}>{tag.name}</Link>
              </Tag>
            ))}
          </div>
        ) : (
          !loadingTags && <Typography.Text type="secondary">暂无热门标签</Typography.Text>
        )}
      </Spin>
    </Card>
  );
};

export default PopularTagsCategories;