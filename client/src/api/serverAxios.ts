import { useAppStore } from '@/stores/app';
import axios from 'axios';
import { logout } from './auth';
import type { TokenDto } from '@/models/tokenDto';

const HOST_URL = "/api";
const REFRESH_TOKEN_URL = '/auth/refresh';
const TEN_MINUTES_IN_MS = 10 * 60 * 1000;

const serverApi = axios.create({
  baseURL: HOST_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});


// --- Request Interceptor ---
// Adds the access token to outgoing requests if available
serverApi.interceptors.request.use(
  (config) => {
    const authStore = useAppStore();
    console.log(authStore.authToken)
    if (authStore.authToken !== "") {
      config.headers = config.headers || {};
      config.headers['Authorization'] = `Bearer ${authStore.authToken}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// --- Response Error Interceptor ---
serverApi.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const authStore = useAppStore();
    if (error.response?.status === 401 && error.config.url !== REFRESH_TOKEN_URL) {
      console.error('Received 401 error. Proactive refresh might have failed or token expired. Logging out.');
      const router = useRouter()
      authStore.authToken = ""
      // BUG: router cant be used there
      await router.replace("/login")
      stopPeriodicRefresh()
    }
    return Promise.reject(error);
  }
);


// --- Token Refresh Function ---
const refreshToken = async () => {
  const authStore = useAppStore();
  console.log('Attempting scheduled token refresh...');

  try {
    const res = await serverApi.post<TokenDto>(REFRESH_TOKEN_URL);
    if (res.status !== 200) {
      throw Error
    }

    authStore.authToken = res.data.accessToken;
    console.log('Token refreshed successfully via schedule.');
  } catch (error) {
    console.error('Unable to refresh token via schedule:', error);
    const rez = await logout()
    if (!rez) return

    const router = useRouter()
    authStore.authToken = ""
    await router.replace("/login")
    stopPeriodicRefresh()
  }
};

// --- Setup Periodic Token Refresh ---
let refreshIntervalId: number | undefined;

export const startPeriodicRefresh = () => {
  if (refreshIntervalId) {
    clearInterval(refreshIntervalId);
  }
  // Immediately refresh token on startup if authenticated, then set interval
  const authStore = useAppStore();
  if (authStore.user) {
    refreshToken();
  }
  refreshIntervalId = setInterval(refreshToken, TEN_MINUTES_IN_MS);
  console.log(`Token refresh scheduled every ${TEN_MINUTES_IN_MS / 60000} minutes.`);
};

export const stopPeriodicRefresh = () => {
  if (refreshIntervalId) {
    clearInterval(refreshIntervalId);
    refreshIntervalId = undefined;
    console.log('Token refresh schedule stopped.');
  }
};


export default serverApi
