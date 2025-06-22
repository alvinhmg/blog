import { api } from './index';

// 图片上传相关API
export const uploadAPI = {
  // 上传图片
  uploadImage: (formData) => {
    return api.post('/misc/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};