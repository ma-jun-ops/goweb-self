<template>
  <div class="edit">
    <el-input v-model="form.title" placeholder="标题" />
    <el-input
      v-model="tagsText"
      placeholder="标签，用逗号分隔"
      style="margin: 12px 0"
    />
    <el-input
      v-model="form.content"
      type="textarea"
      :rows="15"
      placeholder="正文内容"
    />
    <div style="margin-top: 12px">
      <input
        ref="fileInput"
        type="file"
        accept=".docx,.doc"
        style="display: none"
        @change="onFileChange"
      />
      <el-button @click="$refs.fileInput.click()">上传 Word</el-button>
      <el-button type="primary" @click="onSave">保存</el-button>
      <el-button @click="$router.back()">返回</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createArticle, updateArticle, getArticle, uploadFile } from '../api/article'

const route = useRoute()
const router = useRouter()

const form = ref({ title: '', content: '', tags: [], attachments: [] })
const tagsText = ref('')
const id = route.params.id
const fileInput = ref(null)

onMounted(async () => {
  if (id) {
    const res = await getArticle(id)
    form.value = res.data.article
    tagsText.value = (form.value.tags || []).join(',')
  }
})

const onFileChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  const res = await uploadFile(file)
  form.value.content = res.data.text
  form.value.attachments = [{
    fileId: res.data.fileId,
    filename: res.data.filename,
    text: res.data.text,
  }]
  ElMessage.success('上传成功，已提取文字')
  e.target.value = ''
}

const onSave = async () => {
  const data = {
    title: form.value.title,
    content: form.value.content,
    tags: tagsText.value.split(',').map((s) => s.trim()).filter(Boolean),
    attachments: form.value.attachments || [],
  }
  if (id) {
    await updateArticle(id, data)
    ElMessage.success('更新成功')
  } else {
    await createArticle(data)
    ElMessage.success('发布成功')
  }
  router.push('/')
}
</script>

<style scoped>
.edit {
  background-color: #fff;
}
</style>
