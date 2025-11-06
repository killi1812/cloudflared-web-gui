<template>
  <v-card>
    <v-card-title class="d-flex justify-center">
    </v-card-title>
    <v-card-text class="text-center">
      <h2 class="text-h5">{{ user.username }}</h2>
    </v-card-text>
  </v-card>

</template>
<script lang="ts">

export interface UserPreviewProps {
  user: user
}

</script>
<script setup lang="ts">
import type { user } from '@/models/user';

defineProps<UserPreviewProps>()

</script>
