<template>
  <div>
    <AppHeader v-if="isLoggedIn" />
    <SnackbarProvider>
      <v-main>
        <router-view v-slot="{ Component, route }">
          <component :is="Component" :key="route" />
        </router-view>
      </v-main>
    </SnackbarProvider>
  </div>
</template>

<script lang="ts" setup>
const router = useRouter()
const isLoggedIn = computed(() => (router.currentRoute.value.fullPath != "/login") && (router.currentRoute.value.fullPath != "/"))
//
</script>
