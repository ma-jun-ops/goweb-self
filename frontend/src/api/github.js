import axios from 'axios'

const http = axios.create({
  baseURL: '/api',
  timeout: 25000,
})

export const listGithubFriends = () => http.get('/github/friends')
