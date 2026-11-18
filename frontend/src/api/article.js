import axios from 'axios'

const http = axios.create({
  // 相对路径：生产环境由 Nginx 反代到后端，开发环境由 Vite proxy 转发
  baseURL: '/api',
  timeout: 10000,
})

export const listArticles = () => http.get('/articles')
export const getArticle = (id) => http.get(`/articles/${id}`)
export const createArticle = (data) => http.post('/articles', data)
export const updateArticle = (id, data) => http.put(`/articles/${id}`, data)
export const deleteArticle = (id) => http.delete(`/articles/${id}`)

// 全文搜索
export const searchArticles = (q) => http.get('/search', { params: { q } })

// 上传 Word 文件
export const uploadFile = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return http.post('/files/upload', formData)
}
