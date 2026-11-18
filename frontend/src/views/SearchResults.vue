<template>
  <div>
    <h2 class="search-heading">搜索 "{{ q }}" 的结果</h2>
    <div v-for="a in articles" :key="a.id" class="article-card">
      <h3 class="article-title" @click="$router.push(`/article/${a.id}`)">{{ a.title }}</h3>
      <div class="article-preview">
        {{ (a.content || '').slice(0, 200) }}{{ (a.content || '').length > 200 ? '...' : '' }}
      </div>
      <div class="article-tags">
        <el-tag v-for="t in a.tags" :key="t" size="small" class="tag-dark">{{ t }}</el-tag>
      </div>
    </div>
    <div v-if="loaded && !articles.length" class="empty-state">
      <p>没有找到与 "{{ q }}" 相关的文章</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { searchArticles } from '../api/article'

const route = useRoute()
const q = ref(route.query.q || '')
const articles = ref([])
const loaded = ref(false)

const doSearch = async () => {
  if (!q.value) {
    loaded.value = true
    return
  }
  try {
    const res = await searchArticles(q.value)
    articles.value = res.data.data || []
  } catch (e) {
    console.error('搜索失败', e)
    articles.value = []
  } finally {
    loaded.value = true
  }
}

// 路由 query 变化时（例如再次搜索）重新执行
watch(() => route.query.q, (val) => {
  q.value = val || ''
  doSearch()
})

doSearch()
</script>

<style scoped>
.search-heading {
  font-family: var(--font-serif);
  font-size: 26px;
  color: var(--leaf-deep);
  margin-bottom: 18px;
}

.article-card {
  background: rgba(255, 250, 243, 0.82);
  border: 1px solid var(--paper-edge);
  border-radius: 18px;
  padding: 24px 26px;
  margin-bottom: 16px;
  transition: box-shadow 0.2s, transform 0.2s, border-color 0.2s;
}

.article-card:hover {
  box-shadow: var(--shadow-lift);
  transform: translateY(-2px);
  border-color: rgba(63, 107, 79, 0.2);
}

.article-title {
  font-family: var(--font-serif);
  font-size: 24px;
  line-height: 1.35;
  color: var(--ink);
  cursor: pointer;
  margin-bottom: 8px;
  transition: color 0.2s;
}

.article-title:hover {
  color: var(--leaf);
}

.article-preview {
  font-size: 15px;
  color: var(--muted);
  line-height: 1.7;
  margin-bottom: 16px;
}

.article-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.tag-dark {
  background: rgba(63, 107, 79, 0.1) !important;
  border-color: transparent !important;
  color: var(--leaf-deep) !important;
}

.empty-state {
  text-align: center;
  padding: 48px 20px;
  color: var(--muted);
  background: rgba(255, 250, 243, 0.7);
  border-radius: 18px;
  border: 1px dashed var(--paper-edge);
}

@media (max-width: 560px) {
  .article-card {
    padding: 20px;
  }

  .search-heading {
    font-size: 23px;
  }
}
</style>
