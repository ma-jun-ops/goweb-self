import { createRouter, createWebHistory } from 'vue-router'
import ArticleList from '../views/ArticleList.vue'

const routes = [
  { path: '/', name: 'list', component: ArticleList },
  { path: '/article/:id', name: 'detail', component: () => import('../views/ArticleDetail.vue') },
  { path: '/edit/:id?', name: 'edit', component: () => import('../views/ArticleEdit.vue') },
]

export default createRouter({ history: createWebHistory(), routes })
