<template>
  <div class="friends">
    <section class="friends-hero">
      <p class="friends-kicker">people i follow</p>
      <h2 class="friends-title">GitHub 友邻</h2>
      <p class="friends-desc">我在 GitHub 上关注的开发者，都是一些有趣、值得学习的人。</p>
    </section>

    <div v-if="loading" class="friends-grid">
      <div v-for="i in 6" :key="i" class="friend-card skeleton"></div>
    </div>

    <div v-else-if="error" class="friends-empty">
      <p>GitHub 数据暂时获取失败，请稍后再试。</p>
      <button class="retry-btn" type="button" @click="load">重新加载</button>
    </div>

    <div v-else-if="friends.length" class="friends-grid">
      <a
        v-for="friend in friends"
        :key="friend.login"
        class="friend-card"
        :href="friend.htmlUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        <div class="friend-head">
          <img class="friend-avatar" :src="friend.avatarUrl" :alt="friend.login" loading="lazy" />
          <div class="friend-identity">
            <h3 class="friend-name">{{ friend.name || friend.login }}</h3>
            <p class="friend-login">@{{ friend.login }}</p>
          </div>
        </div>
        <p class="friend-bio">{{ friend.bio || '这个人很神秘，什么也没写。' }}</p>
        <span class="friend-more">访问主页 →</span>
      </a>
    </div>

    <div v-else class="friends-empty">
      <p>还没有关注任何人。</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { listGithubFriends } from '../api/github'

const friends = ref([])
const loading = ref(true)
const error = ref(false)

const load = async () => {
  loading.value = true
  error.value = false
  try {
    const res = await listGithubFriends()
    friends.value = res.data.data || []
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.friends-hero {
  margin-bottom: 26px;
}

.friends-kicker {
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--faint);
  margin-bottom: 8px;
}

.friends-title {
  font-family: var(--font-serif);
  font-size: 32px;
  line-height: 1.2;
  color: var(--leaf-deep);
  margin-bottom: 10px;
}

.friends-desc {
  font-size: 15px;
  line-height: 1.7;
  color: var(--muted);
}

.friends-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.friend-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: 176px;
  padding: 22px;
  background: rgba(255, 250, 243, 0.82);
  border: 1px solid var(--paper-edge);
  border-radius: 18px;
  box-shadow: var(--shadow-soft);
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.friend-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lift);
  border-color: rgba(63, 107, 79, 0.24);
}

.friend-head {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.friend-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  flex: none;
  border: 2px solid rgba(63, 107, 79, 0.14);
  background: #fff;
}

.friend-identity {
  min-width: 0;
}

.friend-name {
  font-family: var(--font-serif);
  font-size: 20px;
  line-height: 1.3;
  color: var(--ink);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.friend-login {
  margin-top: 2px;
  font-size: 12px;
  color: var(--faint);
}

.friend-bio {
  font-size: 13px;
  line-height: 1.65;
  color: var(--muted);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.friend-more {
  margin-top: auto;
  font-size: 12px;
  font-weight: 500;
  color: var(--leaf);
}

.skeleton {
  background: linear-gradient(100deg, rgba(255, 250, 243, 0.7) 30%, rgba(244, 236, 224, 0.9) 50%, rgba(255, 250, 243, 0.7) 70%);
  background-size: 220% 100%;
  animation: shimmer 1.4s infinite;
}

@keyframes shimmer {
  from {
    background-position: 180% 0;
  }
  to {
    background-position: -40% 0;
  }
}

.friends-empty {
  padding: 48px 20px;
  text-align: center;
  color: var(--muted);
  background: rgba(255, 250, 243, 0.7);
  border: 1px dashed var(--paper-edge);
  border-radius: 18px;
}

.retry-btn {
  margin-top: 14px;
  height: 34px;
  padding: 0 18px;
  border: 1px solid rgba(63, 107, 79, 0.3);
  border-radius: 999px;
  background: transparent;
  color: var(--leaf);
  font-size: 13px;
  cursor: pointer;
  transition: background 0.2s ease, color 0.2s ease;
}

.retry-btn:hover {
  background: var(--leaf);
  color: #fff;
}

@media (max-width: 640px) {
  .friends-grid {
    grid-template-columns: 1fr;
  }

  .friends-title {
    font-size: 26px;
  }

  .friend-card {
    min-height: 0;
    padding: 18px;
  }
}
</style>
