import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import ArticleList from '../views/ArticleList.vue'

const routes = [
  { path: '/', name: 'home', component: Home },
  { path: '/articles', name: 'list', component: ArticleList },
  { path: '/article/:id', name: 'detail', component: () => import('../views/ArticleDetail.vue') },
  { path: '/edit/:id?', name: 'edit', component: () => import('../views/ArticleEdit.vue') },
  { path: '/search', name: 'search', component: () => import('../views/SearchResults.vue') },
  { path: '/friends', name: 'friends', component: () => import('../views/Friends.vue') },
]

export default createRouter({ history: createWebHistory(), routes })
