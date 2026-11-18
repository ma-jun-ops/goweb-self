<template>
  <div class="home">
    <section class="hero">
      <p class="hero-desc">记录生活、工作与学习中的感悟与随笔</p>
    </section>

    <section class="latest-section">
      <div class="section-head">
        <h2 class="section-title">最近文章</h2>
        <button v-if="latestArticles.length" class="text-btn" type="button" @click="$router.push('/articles')">查看全部</button>
      </div>
      <div v-if="latestArticles.length" class="latest-list">
        <article v-for="a in latestArticles" :key="a.id" class="latest-card" @click="$router.push(`/article/${a.id}`)">
          <div class="latest-card-body">
            <h3 class="latest-card-title">{{ a.title }}</h3>
            <p class="latest-card-meta">{{ formatTime(a.createdAt) }}</p>
            <p class="latest-card-desc">{{ (a.content || '').slice(0, 120) }}{{ (a.content || '').length > 120 ? '...' : '' }}</p>
            <div class="latest-card-tags">
              <span v-for="t in a.tags" :key="t" class="tag">{{ t }}</span>
            </div>
          </div>
        </article>
      </div>
      <div v-else class="empty-state">
        <p>还没有文章</p>
        <button class="hero-btn primary" type="button" @click="$router.push('/edit')">写下第一篇</button>
      </div>
    </section>
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
.hero {
  text-align: left;
  margin-bottom: 32px;
}

.hero-desc {
  color: var(--muted);
  font-size: 17px;
  line-height: 1.7;
  margin: 0;
  max-width: 560px;
}

.hero-actions {
  display: flex;
  justify-content: flex-start;
  gap: 12px;
}

.hero-btn {
  min-width: 128px;
  padding: 10px 22px;
  border: 1px solid var(--paper-edge);
  border-radius: 999px;
  color: var(--ink);
  font-size: 14px;
  font-family: inherit;
  cursor: pointer;
  background: var(--paper);
}

.hero-btn.primary {
  background: var(--leaf);
  border-color: var(--leaf);
  color: #fffaf3;
  font-weight: 600;
}

.hero-btn.primary:hover {
  background: var(--leaf-deep);
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.section-title {
  font-family: var(--font-serif);
  font-size: 26px;
  color: var(--leaf-deep);
}

.text-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 16px;
  border: 1px solid var(--paper-edge);
  border-radius: 999px;
  background: rgba(255, 250, 243, 0.55);
  color: var(--leaf-deep);
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s, transform 0.15s;
}

.text-btn:hover {
  background: var(--paper);
  border-color: rgba(63, 107, 79, 0.35);
  transform: translateY(-1px);
}

.text-btn::after {
  content: '→';
  font-size: 13px;
}

.latest-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.latest-card {
  background: rgba(255, 250, 243, 0.82);
  border: 1px solid var(--paper-edge);
  border-radius: 16px;
  padding: 22px 24px;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s, border-color 0.2s;
}

.latest-card:hover {
  box-shadow: var(--shadow-lift);
  transform: translateY(-2px);
  border-color: rgba(63, 107, 79, 0.22);
}

.latest-card-title {
  font-family: var(--font-serif);
  font-size: 23px;
  line-height: 1.35;
  color: var(--ink);
  margin-bottom: 6px;
}

.latest-card-meta {
  font-size: 12px;
  color: var(--faint);
  margin-bottom: 10px;
}

.latest-card-desc {
  font-size: 15px;
  color: var(--muted);
  line-height: 1.65;
  margin-bottom: 12px;
}

.latest-card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.tag {
  font-size: 12px;
  color: var(--leaf-deep);
  background: rgba(63, 107, 79, 0.1);
  padding: 3px 10px;
  border-radius: 999px;
}

.empty-state {
  text-align: center;
  padding: 48px 20px;
  color: var(--muted);
  background: rgba(255, 250, 243, 0.7);
  border-radius: 18px;
  border: 1px dashed var(--paper-edge);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}
</style>
