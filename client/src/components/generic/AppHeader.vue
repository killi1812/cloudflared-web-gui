<template>
  <v-app-bar color="primary" :elevation="8" height="60">

    <v-btn icon @click="router.back()">
      <v-icon>mdi-arrow-left</v-icon>
    </v-btn>

    <img id="logo" src="../../assets/logo.png" alt="logo" />

    <v-app-bar-title>Cloudflared web</v-app-bar-title>

    <v-spacer></v-spacer>

    <v-menu location="bottom">
      <template v-slot:activator="{ props }">
        <v-btn icon v-bind="props">
          <v-icon> mdi-account-circle </v-icon>
        </v-btn>
      </template>
      <v-list>
        <IconListItem to="/me/myProfile" icon="mdi-account">
          MY PROFILE
        </IconListItem>
        <v-list-item prepend-icon="mdi-logout" @click="handleLogout()">
          Logout
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>

  <v-navigation-drawer class="offset" :rail="!open" rail-width="60" location="left" permanent persistent elevation="2">
    <v-list nav>
      <!-- Information section -->
      <!-- Title -->
      <v-list-item v-show="open" readonly>INFORMATION</v-list-item>
      <v-divider></v-divider>
      <!-- Navigations -->
      <IconListItem to="/home/dashboard" icon="mdi-home" tooltip-text="START.DASHBOARD" :show-tooltip="!open">
        HOME
      </IconListItem>
      <!-- Master Data section -->
      <!-- Title -->
      <v-list-item v-show="open" readonly>MASTER_DATA</v-list-item>
      <v-divider></v-divider>
      <!-- Navigations  -->
      <IconListItem to="/help/help" icon="mdi-help-circle" tooltip-text="START.HELP" :show-tooltip="!open">
        HELP
      </IconListItem>
      <IconListItem to="/help/about" icon="mdi-information" tooltip-text="START.ABOUT" :show-tooltip="!open">
        ABOUT
      </IconListItem>
    </v-list>

    <template v-slot:append>
      <v-list>
        <v-list-item @click="open = !open">
          <v-icon v-if="open">
            mdi-arrow-left
          </v-icon>
          <v-icon v-else>
            mdi-arrow-right
          </v-icon>
        </v-list-item>
      </v-list>
    </template>

  </v-navigation-drawer>

</template>


<script lang="ts" setup>
import { logout } from '@/api/auth';
import { useSnackbar } from './snackbarProvider.vue';


const open = ref(false)
const router = useRouter()
const snackbar = useSnackbar()


async function handleLogout() {
  const rez = await logout()
  if (rez) {
    await router.push("/login")
    return
  }
  snackbar.Error("Failed to login")
}

</script>
<style lang="css" scoped>
#logo {
  width: 40px;
  height: 30px;
}

.offset {
  padding-bottom: 40px;
}

.uppercase {
  text-transform: uppercase;
}
</style>
