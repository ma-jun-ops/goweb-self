<template>
  <div class="app">
    <div class="bg-wash"></div>
    <div class="bg-orbs" aria-hidden="true">
      <span class="orb orb-a"></span>
      <span class="orb orb-b"></span>
      <span class="orb orb-c"></span>
    </div>
    <div class="bg-grain"></div>

    <header class="top-banner">
      <div class="banner-content">
        <p class="site-kicker">a quiet place to write</p>
        <h1 class="site-title" @click="$router.push('/')">My Little Yard</h1>
        <p class="site-subtitle">简单一点，快一点</p>
      </div>
    </header>

    <nav class="navbar">
      <div class="nav-inner" :class="{ wide: showSidebar }">
        <div class="nav-links">
          <a class="nav-link" :class="{ active: route.path === '/' }" @click="$router.push('/')">首页</a>
          <a class="nav-link" :class="{ active: route.path === '/articles' || route.path.startsWith('/article/') }" @click="$router.push('/articles')">文章</a>
          <a class="nav-link" :class="{ active: route.path === '/friends' }" @click="$router.push('/friends')">友邻</a>
        </div>
        <div class="nav-tools">
          <div class="nav-search">
            <input
              v-model="keyword"
              class="nav-search-input"
              type="text"
              placeholder="搜索文章"
              @keyup.enter="onSearch"
            />
            <button class="nav-search-btn" @click="onSearch">搜索</button>
          </div>
          <button class="nav-write" type="button" @click="$router.push('/edit')">写文章</button>
        </div>
      </div>
    </nav>

    <div class="layout" :class="{ 'has-sidebar': showSidebar }">
      <main class="content">
        <router-view />
      </main>
      <aside v-if="showSidebar" class="sidebar">
        <div class="sidebar-card">
          <h3 class="sidebar-title">热门文章</h3>
          <ul class="popular-list">
            <li v-for="a in popularArticles" :key="a.id" class="popular-item">
              <a class="popular-link" @click="$router.push(`/article/${a.id}`)">{{ a.title }}</a>
              <p class="popular-desc">{{ (a.content || '').slice(0, 60) }}...</p>
            </li>
            <li v-if="!popularArticles.length" class="popular-item">
              <p class="popular-desc">暂无文章</p>
            </li>
          </ul>
        </div>
      </aside>
    </div>

    <ChatBot />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ChatBot from './components/ChatBot.vue'
import { listArticles } from './api/article'

const route = useRoute()
const router = useRouter()
const popularArticles = ref([])
const keyword = ref('')

const showSidebar = computed(() => {
  const p = route.path
  return p === '/articles' || p.startsWith('/article/') || p.startsWith('/search')
})

const onSearch = () => {
  const q = keyword.value.trim()
  if (!q) return
  router.push({ path: '/search', query: { q } })
}

onMounted(async () => {
  try {
    const res = await listArticles()
    popularArticles.value = (res.data.data || []).slice(0, 3)
  } catch {}
})
</script>

<style>
:root {
  --ink: #2c2a26;
  --muted: #6f675c;
  --faint: #9a9186;
  --paper: #fffaf3;
  --paper-edge: rgba(70, 58, 42, 0.1);
  --leaf: #3f6b4f;
  --leaf-deep: #2f513c;
  --clay: #c46a4a;
  --font-sans: 'Source Sans 3', 'PingFang SC', 'Microsoft YaHei', 'Noto Sans SC', sans-serif;
  --font-serif: 'Cormorant Garamond', 'Songti SC', 'STSong', 'Noto Serif SC', Georgia, serif;
  --shadow-soft: 0 10px 30px rgba(80, 58, 32, 0.07);
  --shadow-lift: 0 16px 36px rgba(80, 58, 32, 0.11);

  /* Element Plus 主题色跟随页面配色 */
  --el-color-primary: #3f6b4f;
  --el-color-primary-light-3: #6b8f78;
  --el-color-primary-light-5: #9bb5a3;
  --el-color-primary-light-7: #c5d4c9;
  --el-color-primary-light-8: #d9e3db;
  --el-color-primary-light-9: #edf2ee;
  --el-color-primary-dark-2: #2f513c;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body,
#app {
  min-height: 100%;
  color: var(--ink);
  font-family: var(--font-sans);
  background-color: #efe6d4;
}

body {
  background:
    radial-gradient(1100px 720px at 8% -10%, rgba(214, 176, 132, 0.45), transparent 58%),
    radial-gradient(900px 680px at 96% 8%, rgba(156, 186, 154, 0.38), transparent 52%),
    radial-gradient(800px 620px at 78% 108%, rgba(232, 196, 164, 0.42), transparent 55%),
    linear-gradient(180deg, #f6eee0 0%, #efe4d0 46%, #e7eedf 100%);
}

.app {
  position: relative;
  z-index: 1;
  min-height: 100%;
}

.bg-wash,
.bg-orbs,
.bg-grain {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

.bg-orbs {
  overflow: hidden;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(40px);
  opacity: 0.55;
  animation: drift 22s ease-in-out infinite;
}

.orb-a {
  width: 420px;
  height: 420px;
  left: -80px;
  top: 80px;
  background: #d9c19a;
}

.orb-b {
  width: 360px;
  height: 360px;
  right: -60px;
  top: 220px;
  background: #b7cbb4;
  animation-delay: -7s;
}

.orb-c {
  width: 480px;
  height: 480px;
  left: 32%;
  bottom: -180px;
  background: #e8c7b0;
  animation-delay: -13s;
}

@keyframes drift {
  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(28px, -24px) scale(1.08);
  }
}

.bg-grain {
  opacity: 0.22;
  mix-blend-mode: multiply;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='160' height='160'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='.55'/%3E%3C/svg%3E");
}

.top-banner {
  position: relative;
  z-index: 1;
  padding: 40px 20px 20px;
  text-align: center;
}

.site-kicker {
  font-size: 12px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--faint);
  margin-bottom: 10px;
}

.site-title {
  font-family: var(--font-serif);
  font-size: clamp(40px, 6vw, 56px);
  font-weight: 600;
  color: var(--leaf-deep);
  letter-spacing: 0.02em;
  cursor: pointer;
}

.site-subtitle {
  font-size: 15px;
  color: var(--muted);
  margin-top: 8px;
}

.navbar {
  position: relative;
  z-index: 1;
  border-top: 1px solid var(--paper-edge);
  border-bottom: 1px solid var(--paper-edge);
  background: rgba(255, 250, 243, 0.62);
  backdrop-filter: blur(14px);
}

.nav-inner {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 24px;
}

.nav-inner.wide {
  max-width: 1040px;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-link {
  padding: 8px 14px;
  color: var(--muted);
  cursor: pointer;
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
  border-radius: 999px;
  transition: color 0.2s, background 0.2s;
}

.nav-link:hover,
.nav-link.active {
  color: var(--leaf-deep);
  background: rgba(63, 107, 79, 0.1);
}

.nav-tools {
  display: flex;
  align-items: center;
  gap: 10px;
}

.nav-search {
  display: flex;
  align-items: center;
  height: 38px;
  background: var(--paper);
  border: 1px solid var(--paper-edge);
  border-radius: 999px;
  overflow: hidden;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.nav-search:focus-within {
  border-color: rgba(63, 107, 79, 0.45);
  box-shadow: 0 0 0 3px rgba(63, 107, 79, 0.1);
}

.nav-search-input {
  background: transparent;
  border: none;
  color: var(--ink);
  padding: 0 16px;
  height: 100%;
  font-size: 13px;
  outline: none;
  width: 168px;
  font-family: inherit;
}

.nav-search-input::placeholder {
  color: var(--faint);
}

.nav-search-btn {
  background: transparent;
  border: none;
  border-left: 1px solid var(--paper-edge);
  color: var(--leaf-deep);
  padding: 0 16px;
  height: 100%;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: background 0.2s;
}

.nav-search-btn:hover {
  background: rgba(63, 107, 79, 0.08);
}

.nav-write {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 38px;
  padding: 0 20px;
  background: var(--leaf);
  border: none;
  border-radius: 999px;
  color: #fffaf3;
  font-size: 14px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s, transform 0.15s, box-shadow 0.2s;
}

.nav-write:hover {
  background: var(--leaf-deep);
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(63, 107, 79, 0.22);
}

.nav-write:active {
  transform: translateY(0);
}

.layout {
  position: relative;
  z-index: 1;
  max-width: 800px;
  margin: 36px auto 72px;
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 28px;
  padding: 0 24px;
}

.layout.has-sidebar {
  max-width: 1040px;
  grid-template-columns: minmax(0, 1fr) 260px;
  align-items: start;
}

.content {
  min-width: 0;
}

.sidebar {
  width: auto;
  position: sticky;
  top: 24px;
}

.sidebar-card {
  background: rgba(255, 250, 243, 0.78);
  border: 1px solid var(--paper-edge);
  border-radius: 18px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 10px 30px rgba(80, 58, 32, 0.06);
  backdrop-filter: blur(10px);
}

.sidebar-title {
  font-family: var(--font-serif);
  font-size: 20px;
  font-weight: 600;
  color: var(--leaf-deep);
  margin-bottom: 14px;
}

.sidebar-text {
  font-size: 14px;
  color: var(--muted);
  line-height: 1.6;
}

.popular-list {
  list-style: none;
}

.popular-item {
  margin-bottom: 14px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--paper-edge);
}

.popular-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.popular-link {
  color: var(--ink);
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
}

.popular-link:hover {
  color: var(--leaf);
}

.popular-desc {
  font-size: 13px;
  color: var(--muted);
  margin-top: 6px;
  line-height: 1.5;
}

button:focus-visible,
a:focus-visible,
input:focus-visible {
  outline: 2px solid rgba(63, 107, 79, 0.5);
  outline-offset: 2px;
}

@media (max-width: 860px) {
  .layout,
  .layout.has-sidebar {
    grid-template-columns: 1fr;
    max-width: 760px;
  }

  .sidebar {
    position: static;
  }

  .nav-inner {
    flex-wrap: wrap;
    justify-content: center;
    padding: 12px 16px;
    gap: 10px;
  }

  .nav-tools {
    width: 100%;
    justify-content: center;
    flex-wrap: nowrap;
  }

  .nav-search {
    flex: 1;
    min-width: 0;
  }

  .nav-search-input {
    width: 100%;
    min-width: 0;
  }

  .nav-write {
    flex-shrink: 0;
  }
}

@media (max-width: 520px) {
  .top-banner {
    padding: 32px 16px 16px;
  }

  .layout,
  .layout.has-sidebar {
    padding: 0 16px;
    margin: 28px auto 56px;
  }
}
</style>
