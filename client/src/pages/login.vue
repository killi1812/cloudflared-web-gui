<template>
  <div style="margin-top: 20rem;">
    <v-card class="mx-auto pa-12 pb-8" elevation="8" max-width="448" rounded="lg">
      <div class="text-subtitle-1 text-medium-emphasis">Account</div>

      <v-text-field v-model="username" density="compact" placeholder="Username" type="text"
        prepend-inner-icon="mdi-account-outline" variant="outlined" :error-messages="emailError ? [emailError] : []"
        @focus="emailError = ''"></v-text-field>

      <div class="text-subtitle-1 text-medium-emphasis d-flex align-center justify-space-between">
        Password
      </div>

      <v-text-field v-model="password" :append-inner-icon="visible ? 'mdi-eye-off' : 'mdi-eye'"
        :type="visible ? 'text' : 'password'" density="compact" placeholder="Enter your password"
        prepend-inner-icon="mdi-lock-outline" variant="outlined" :error-messages="passwordError ? [passwordError] : []"
        @click:append-inner="visible = !visible" @focus="passwordError = ''"></v-text-field>

      <v-btn class="mb-8" color="blue" size="large" variant="tonal" block :loading="loading" @click="handleLogin">
        Log In
      </v-btn>

    </v-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSnackbar } from '@/components/generic/snackbarProvider.vue';
import { useAppStore } from '@/stores/app';
import { login } from '@/api/auth';
import { startPeriodicRefresh } from '@/api/serverAxios';
import { getLoggedInUserData } from '@/api/user';

const router = useRouter();
const authStore = useAppStore()
const snackbar = useSnackbar()

const username = ref('');
const password = ref('');
const emailError = ref('');
const passwordError = ref('');
const visible = ref(false);
const tosAccepted = ref(false);

const loading = ref(false)

const validateForm = () => {
  let isValid = true;

  if (!username.value.trim()) {
    emailError.value = 'E-mail is required';
    isValid = false;
  }

  if (!password.value) {
    passwordError.value = 'Password is required';
    isValid = false;
  } else if (password.value.length < 6) {
    passwordError.value = 'Password must be at least 6 characters';
    isValid = false;
  }

  return isValid;
};

async function handleLogin() {
  if (!validateForm()) {
    return;
  }
  loading.value = true
  try {
    const rez = await login({ username: username.value, password: password.value });
    if (!rez)
      return
    authStore.authToken = rez.accessToken

    const userRez = await getLoggedInUserData()
    if (userRez)
      authStore.user = userRez

  } catch (error) {
    console.error('Login failed:', error);
    snackbar.Error(`Failed to login`)
    return
  }
  finally {
    loading.value = false
  }
  startPeriodicRefresh()
  await router.push('/home/dashboard');
}

</script>
