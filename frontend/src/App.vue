<template>
  <div class="app">
    <div class="bg-gradient"></div>
    <canvas ref="bgCanvas" class="bg-canvas"></canvas>

    <!-- 顶部装饰区 -->
    <header class="top-banner">
      <div class="banner-content">
        <h1 class="site-title">My Little Yard</h1>
        <p class="site-subtitle">简单一点，快一点</p>
      </div>
    </header>

    <!-- 导航栏 -->
    <nav class="navbar">
      <div class="nav-inner">
        <a class="nav-link" :class="{ active: route.path === '/' }" @click="$router.push('/')">Home</a>
        <a class="nav-link" :class="{ active: route.path.startsWith('/edit') }" @click="$router.push('/edit')">Write</a>
      </div>
    </nav>

    <!-- 主体内容 -->
    <div class="layout">
      <main class="content">
        <router-view />
      </main>
      <aside class="sidebar">
        <div class="sidebar-card">
          <h3 class="sidebar-title">POPULAR POSTS</h3>
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

        <div class="sidebar-card">
          <h3 class="sidebar-title">FOLLOW ME</h3>
          <p class="sidebar-text">Stay connected</p>
        </div>

        <div class="sidebar-card">
          <h3 class="sidebar-title">SUBSCRIBE</h3>
          <p class="sidebar-text">rss feeds</p>
        </div>

        <div class="sidebar-card">
          <h3 class="sidebar-title">CHECK THIS OUT!</h3>
          <ul class="check-list">
            <li>CSS Templates</li>
            <li>Flash Templates</li>
            <li>Flash Layouts</li>
            <li>Flash Website Gallery</li>
          </ul>
        </div>
      </aside>
    </div>

    <ChatBot />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ChatBot from './components/ChatBot.vue'
import { listArticles } from './api/article'

const route = useRoute()
const popularArticles = ref([])
const bgCanvas = ref(null)

// 蓝色粒子星空背景
const initStarfield = () => {
  const canvas = bgCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')

  let width = (canvas.width = window.innerWidth)
  let height = (canvas.height = window.innerHeight)

  const count = Math.min(120, Math.max(60, Math.floor((width * height) / 15000)))
  const stars = Array.from({ length: count }, () => ({
    x: Math.random() * width,
    y: Math.random() * height,
    vx: (Math.random() - 0.5) * 0.4,
    vy: (Math.random() - 0.5) * 0.4,
    r: Math.random() * 1.5 + 0.5,
  }))

  const draw = () => {
    ctx.clearRect(0, 0, width, height)
    for (const s of stars) {
      s.x += s.vx
      s.y += s.vy
      if (s.x < 0 || s.x > width) s.vx *= -1
      if (s.y < 0 || s.y > height) s.vy *= -1

      ctx.beginPath()
      ctx.arc(s.x, s.y, s.r, 0, Math.PI * 2)
      ctx.fillStyle = 'rgba(58, 130, 255, 0.55)'
      ctx.fill()
    }
    // 粒子连线
    for (let i = 0; i < stars.length; i++) {
      for (let j = i + 1; j < stars.length; j++) {
        const dx = stars[i].x - stars[j].x
        const dy = stars[i].y - stars[j].y
        const dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < 120) {
          ctx.beginPath()
          ctx.moveTo(stars[i].x, stars[i].y)
          ctx.lineTo(stars[j].x, stars[j].y)
          ctx.strokeStyle = `rgba(58, 130, 255, ${(1 - dist / 120) * 0.22})`
          ctx.lineWidth = 1
          ctx.stroke()
        }
      }
    }
    requestAnimationFrame(draw)
  }
  draw()

  window.addEventListener('resize', () => {
    width = canvas.width = window.innerWidth
    height = canvas.height = window.innerHeight
  })
}

onMounted(async () => {
  initStarfield()
  try {
    const res = await listArticles()
    popularArticles.value = (res.data.data || []).slice(0, 3)
  } catch {}
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  background-color: #0c1120;
  color: #e2e8f0;
  font-family: 'Georgia', 'Times New Roman', serif;
}

.app {
  position: relative;
  z-index: 1;
  min-height: 100%;
  background-color: transparent;
}

/* 动态粒子背景层 */
.bg-canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}

/* 流动蓝色渐变光晕层 */
.bg-gradient {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 20% 30%, rgba(58, 130, 255, 0.12), transparent 45%),
    radial-gradient(circle at 80% 70%, rgba(45, 100, 220, 0.12), transparent 45%),
    radial-gradient(circle at 55% 15%, rgba(90, 160, 255, 0.10), transparent 50%),
    radial-gradient(circle at 10% 85%, rgba(40, 90, 200, 0.10), transparent 50%);
  animation: aurora 18s ease-in-out infinite;
}

@keyframes aurora {
  0%, 100% {
    opacity: 0.6;
    transform: translate(0, 0) scale(1);
  }
  25% {
    opacity: 0.9;
    transform: translate(40px, -30px) scale(1.1);
  }
  50% {
    opacity: 1;
    transform: translate(-30px, 20px) scale(1.2);
  }
  75% {
    opacity: 0.8;
    transform: translate(20px, 40px) scale(1.05);
  }
}

/* 顶部装饰区 */
.top-banner {
  background: linear-gradient(135deg, #0a0e1a 0%, #141c33 50%, #0a0e1a 100%);
  padding: 40px 0 30px;
  text-align: center;
  border-bottom: 1px solid #1e293b;
}

.site-title {
  font-size: 36px;
  color: #3a82ff;
  font-weight: bold;
  letter-spacing: 2px;
}

.site-subtitle {
  font-size: 14px;
  color: #8895a7;
  margin-top: 6px;
}

/* 导航栏 */
.navbar {
  background-color: #0d1324;
  border-bottom: 1px solid #1e293b;
}

.nav-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  gap: 0;
}

.nav-link {
  padding: 12px 24px;
  color: #94a3b8;
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
  transition: color 0.2s, background 0.2s;
}

.nav-link:hover, .nav-link.active {
  color: #3a82ff;
  background-color: #141c33;
}

/* 双栏布局 */
.layout {
  max-width: 1100px;
  margin: 20px auto;
  display: flex;
  gap: 20px;
  padding: 0 20px;
}

.content {
  flex: 1;
  min-width: 0;
}

.sidebar {
  width: 280px;
  flex-shrink: 0;
}

/* 侧边栏卡片 */
.sidebar-card {
  background-color: #141c33;
  border: 1px solid #1e293b;
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;
}

.sidebar-title {
  font-size: 14px;
  font-weight: bold;
  color: #3a82ff;
  letter-spacing: 1px;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #24304a;
}

.sidebar-text {
  font-size: 13px;
  color: #8895a7;
}

.popular-list {
  list-style: none;
}

.popular-item {
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #24304a;
}

.popular-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.popular-link {
  color: #3a82ff;
  cursor: pointer;
  font-size: 13px;
  text-decoration: underline;
}

.popular-link:hover {
  color: #6aa1ff;
}

.popular-desc {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
  line-height: 1.4;
}

.check-list {
  list-style: none;
}

.check-list li {
  padding: 6px 0;
  font-size: 13px;
  color: #94a3b8;
  border-bottom: 1px solid #24304a;
}

.check-list li:last-child {
  border-bottom: none;
}
</style>
