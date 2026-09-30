<template>
  <div class="home">
    <section class="hero">
      <h1 class="hero-title">Welcome to My Little Yard</h1>
      <p class="hero-desc">
        记录生活、工作与学习中的感悟与随笔
      </p>
      <div class="hero-actions">
        <a class="hero-btn primary" @click="$router.push('/articles')">Browse Articles</a>
        <a class="hero-btn" @click="$router.push('/edit')">Write Article</a>
      </div>
    </section>

    <div class="latest-section">
      <h2 class="section-title">Latest Posts</h2>
      <div v-if="latestArticles.length" class="latest-grid">
        <div v-for="a in latestArticles" :key="a.id" class="latest-card" @click="$router.push(`/article/${a.id}`)">
          <h3 class="latest-card-title">{{ a.title }}</h3>
          <p class="latest-card-meta">{{ formatTime(a.createdAt) }}</p>
          <p class="latest-card-desc">{{ (a.content || '').slice(0, 120) }}{{ (a.content || '').length > 120 ? '...' : '' }}</p>
          <div class="latest-card-tags">
            <span v-for="t in a.tags" :key="t" class="tag">{{ t }}</span>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">
        <p>No articles yet. Click "Write Article" to create your first post.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listArticles } from '../api/article'

const latestArticles = ref([])

const load = async () => {
  try {
    const res = await listArticles()
    latestArticles.value = (res.data.data || []).slice(0, 6)
  } catch {}
}

const formatTime = (s) => new Date(s).toLocaleString()

onMounted(load)
</script>

<style scoped>
.home {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 20px 48px;
}

.hero {
  text-align: center;
  margin-bottom: 48px;
  padding: 40px 20px;
}

.hero-title {
  font-size: clamp(28px, 5vw, 40px);
  font-weight: 700;
  color: #3a82ff;
  letter-spacing: 1px;
  margin-bottom: 12px;
}

.hero-desc {
  color: #8895a7;
  font-size: 16px;
  line-height: 1.6;
  max-width: 480px;
  margin: 0 auto 28px;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}

.hero-btn {
  display: inline-block;
  padding: 10px 28px;
  border: 1px solid #2a3654;
  border-radius: 4px;
  color: #cbd5e1;
  font-size: 14px;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
}

.hero-btn:hover {
  border-color: #3a82ff;
  color: #3a82ff;
}

.hero-btn.primary {
  background-color: #3a82ff;
  border-color: #3a82ff;
  color: #0c1120;
  font-weight: bold;
}

.hero-btn.primary:hover {
  background-color: #6aa1ff;
  border-color: #6aa1ff;
  color: #0c1120;
}

/* Latest Posts */
.latest-section {
  margin-top: 20px;
}

.section-title {
  font-size: 18px;
  color: #3a82ff;
  letter-spacing: 1px;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #1e293b;
}

.latest-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.latest-card {
  background-color: #141c33;
  border: 1px solid #1e293b;
  border-radius: 4px;
  padding: 20px;
  cursor: pointer;
  transition: border-color 0.2s, transform 0.2s;
}

.latest-card:hover {
  border-color: #2a3654;
  transform: translateY(-2px);
}

.latest-card-title {
  font-size: 17px;
  color: #3a82ff;
  margin-bottom: 8px;
}

.latest-card-meta {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 10px;
}

.latest-card-desc {
  font-size: 14px;
  color: #cbd5e1;
  line-height: 1.6;
  margin-bottom: 12px;
}

.latest-card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.tag {
  font-size: 12px;
  color: #94a3b8;
  background-color: #1e293b;
  padding: 2px 8px;
  border-radius: 3px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #64748b;
  font-size: 15px;
}
</style>
