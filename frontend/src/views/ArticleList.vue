<template>
  <div>
    <div v-for="a in articles" :key="a.id" class="article-card">
      <h2 class="article-title" @click="$router.push(`/article/${a.id}`)">{{ a.title }}</h2>
      <p class="article-meta">
        <span>{{ formatTime(a.createdAt) }}</span>
        <span class="meta-dot">·</span>
        <span>Free Templates by My Little Yard</span>
      </p>
      <div class="article-preview">
        {{ (a.content || '').slice(0, 200) }}{{ (a.content || '').length > 200 ? '...' : '' }}
      </div>
      <div class="article-footer">
        <div class="article-tags">
          <el-tag v-for="t in a.tags" :key="t" size="small" class="tag-dark">{{ t }}</el-tag>
        </div>
        <div class="article-actions">
          <a class="action-link" @click="$router.push(`/edit/${a.id}`)">Continue reading</a>
          <span class="meta-dot">|</span>
          <a class="action-link" @click="onDelete(a)">Delete</a>
        </div>
      </div>
    </div>
    <div v-if="!articles.length" class="empty-state">
      <p>暂无文章，点击右上角 Write 开始创作</p>
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
  background-color: #141c33;
  border: 1px solid #1e293b;
  border-radius: 4px;
  padding: 20px;
  margin-bottom: 20px;
  transition: border-color 0.2s;
}

.article-card:hover {
  border-color: #2a3654;
}

.article-title {
  font-size: 20px;
  color: #3a82ff;
  cursor: pointer;
  margin-bottom: 8px;
  transition: color 0.2s;
}

.article-title:hover {
  color: #6aa1ff;
}

.article-meta {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 12px;
}

.meta-dot {
  margin: 0 6px;
  color: #334155;
}

.article-preview {
  font-size: 14px;
  color: #cbd5e1;
  line-height: 1.6;
  margin-bottom: 16px;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.article-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.tag-dark {
  background-color: #1e293b !important;
  border-color: #2a3654 !important;
  color: #94a3b8 !important;
}

.article-actions {
  font-size: 13px;
}

.action-link {
  color: #3a82ff;
  cursor: pointer;
  text-decoration: underline;
}

.action-link:hover {
  color: #6aa1ff;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #64748b;
  font-size: 16px;
}
</style>
