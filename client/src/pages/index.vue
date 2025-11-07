<template>
</template>

<script lang="ts" setup>
import { useAppStore } from '@/stores/app'

const router = useRouter()
const app = useAppStore()

onMounted(async () => {
  if (app.authToken == "")
    await router.replace("/login")
  else
    await router.replace("/home/dashboard")
})
</script>
