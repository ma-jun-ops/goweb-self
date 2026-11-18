<template>
  <div>
    <div v-for="a in articles" :key="a.id" class="article-card">
      <h2 class="article-title" @click="$router.push(`/article/${a.id}`)">{{ a.title }}</h2>
      <p class="article-meta">
        <span>{{ formatTime(a.createdAt) }}</span>
        <span class="meta-dot">·</span>
        <span>My Little Yard</span>
      </p>
      <div class="article-preview">
        {{ (a.content || '').slice(0, 200) }}{{ (a.content || '').length > 200 ? '...' : '' }}
      </div>
      <div class="article-footer">
        <div class="article-tags">
          <span v-for="t in a.tags" :key="t" class="tag">{{ t }}</span>
        </div>
        <div class="article-actions">
          <button class="pill-btn read" type="button" @click.stop="$router.push(`/article/${a.id}`)">阅读</button>
          <button class="pill-btn" type="button" @click.stop="$router.push(`/edit/${a.id}`)">编辑</button>
          <button class="pill-btn danger" type="button" @click.stop="onDelete(a)">删除</button>
        </div>
      </div>
    </div>
    <div v-if="!articles.length" class="empty-state">
      <p>暂无文章</p>
      <button class="pill-btn primary" type="button" @click="$router.push('/edit')">开始创作</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listArticles, deleteArticle } from '../api/article'

const articles = ref([])

const load = async () => {
  const res = await listArticles()
  articles.value = res.data.data || []
}

const onDelete = async (row) => {
  await ElMessageBox.confirm('确定删除这篇文章吗？', '提示', { type: 'warning' })
  await deleteArticle(row.id)
  ElMessage.success('删除成功')
  load()
}

const formatTime = (s) => new Date(s).toLocaleString()

onMounted(load)
</script>

<style scoped>
.article-card {
  background: rgba(255, 250, 243, 0.82);
  border: 1px solid var(--paper-edge);
  border-radius: 18px;
  padding: 24px 26px;
  margin-bottom: 16px;
  box-shadow: var(--shadow-soft);
  transition: box-shadow 0.2s, transform 0.2s, border-color 0.2s;
}

.article-card:hover {
  box-shadow: var(--shadow-lift);
  transform: translateY(-2px);
  border-color: rgba(63, 107, 79, 0.2);
}

.article-title {
  font-family: var(--font-serif);
  font-size: 27px;
  line-height: 1.35;
  color: var(--ink);
  cursor: pointer;
  margin-bottom: 8px;
  transition: color 0.2s;
}

.article-title:hover {
  color: var(--leaf);
}

.article-meta {
  font-size: 13px;
  color: var(--faint);
  margin-bottom: 12px;
}

.meta-dot {
  margin: 0 6px;
  color: var(--faint);
}

.article-preview {
  font-size: 15px;
  color: var(--muted);
  line-height: 1.7;
  margin-bottom: 16px;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  padding-top: 14px;
  border-top: 1px solid rgba(70, 58, 42, 0.08);
}

.article-tags {
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

.article-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.pill-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 32px;
  padding: 0 15px;
  border: 1px solid var(--paper-edge);
  background: transparent;
  color: var(--muted);
  border-radius: 999px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  white-space: nowrap;
  transition: color 0.18s, background 0.18s, border-color 0.18s;
}

.pill-btn:hover {
  color: var(--leaf-deep);
  border-color: rgba(63, 107, 79, 0.35);
  background: rgba(63, 107, 79, 0.06);
}

.pill-btn.read {
  background: rgba(63, 107, 79, 0.1);
  border-color: transparent;
  color: var(--leaf-deep);
}

.pill-btn.read:hover {
  background: rgba(63, 107, 79, 0.16);
}

.pill-btn.primary {
  background: var(--leaf);
  border-color: var(--leaf);
  color: #fffaf3;
  font-weight: 600;
  padding: 0 20px;
}

.pill-btn.primary:hover {
  background: var(--leaf-deep);
  border-color: var(--leaf-deep);
  color: #fffaf3;
}

.pill-btn.danger {
  color: var(--clay);
}

.pill-btn.danger:hover {
  color: #a94f33;
  border-color: rgba(196, 106, 74, 0.4);
  background: rgba(196, 106, 74, 0.07);
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

@media (max-width: 560px) {
  .article-card {
    padding: 20px;
  }

  .article-title {
    font-size: 24px;
  }

  .article-footer {
    flex-direction: column;
    align-items: stretch;
  }

  .article-actions {
    justify-content: flex-end;
  }
}
</style>
