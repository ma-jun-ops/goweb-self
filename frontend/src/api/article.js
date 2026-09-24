import axios from 'axios'

const http = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
})

export const listArticles = () => http.get('/articles')
export const getArticle = (id) => http.get(`/articles/${id}`)
export const createArticle = (data) => http.post('/articles', data)
export const updateArticle = (id, data) => http.put(`/articles/${id}`, data)
export const deleteArticle = (id) => http.delete(`/articles/${id}`)

// 上传 Word 文件
export const uploadFile = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return http.post('/files/upload', formData)
}
