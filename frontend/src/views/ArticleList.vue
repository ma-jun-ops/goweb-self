<template>
  <div>
    <el-table :data="articles" style="width: 100%">
      <el-table-column label="标题">
        <template #default="{ row }">
          <a class="title-link" @click="$router.push(`/article/${row.id}`)">{{ row.title }}</a>
        </template>
      </el-table-column>
      <el-table-column label="标签">
        <template #default="{ row }">
          <el-tag v-for="t in row.tags" :key="t" style="margin-right: 4px">{{ t }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="时间" width="180">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push(`/edit/${row.id}`)">编辑</el-button>
          <el-button size="small" type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
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
.title-link {
  color: #409eff;
  cursor: pointer;
}
</style>
