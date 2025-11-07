<template>
  <v-container>
    <!-- Page Header & Main Actions -->
    <v-toolbar flat color="transparent" class="mb-4">
      <v-toolbar-title class="text-h5 font-weight-bold">Tunnels</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn color="primary" @click="openCreateDialog()" prepend-icon="mdi-plus" :disabled="isLoading">
        New Tunnel
      </v-btn>
    </v-toolbar>

    <!-- Main Loading Bar -->
    <v-progress-linear indeterminate :active="isLoading" color="primary" class="mb-4"></v-progress-linear>

    <!-- Empty State -->
    <v-card v-if="!isLoading && tunnels && tunnels.length === 0" variant="tonal" class="pa-4 text-center">
      <v-card-title class="justify-center">No Tunnels Found</v-card-title>
      <v-card-text>Get started by creating your first tunnel.</v-card-text>
      <v-card-actions class="justify-center">
        <v-btn color="primary" @click="openCreateDialog()">
          Create Tunnel
        </v-btn>
      </v-card-actions>
    </v-card>

    <!-- Tunnels List -->
    <v-expansion-panels v-if="tunnels && tunnels.length > 0" variant="inset" class="tunnel-list">
      <v-expansion-panel v-for="tunnel in tunnels" :key="tunnel.id">
        <v-expansion-panel-title>
          <template #default>
            <v-icon icon="mdi-tunnel" start class="mr-3"></v-icon>
            <span class="font-weight-medium">{{ tunnel.name }}</span>
            <v-chip size="small" variant="tonal" class="ml-4" prepend-icon="mdi-clock-outline">
              Created: {{ tunnel.created_at }}
            </v-chip>
            <v-spacer></v-spacer>

            <!-- Control Buttons -->
            <div class="mr-2" @click.stop>
              <v-tooltip location="top" text="Start Tunnel">
                <template #activator="{ props }">
                  <v-btn v-bind="props" icon="mdi-play-circle" color="success" variant="text" size="small"
                    :disabled="isLoading" @click="handleStart(tunnel)"></v-btn>
                </template>
              </v-tooltip>
              <v-tooltip location="top" text="Stop Tunnel">
                <template #activator="{ props }">
                  <v-btn v-bind="props" icon="mdi-stop-circle" color="warning" variant="text" size="small"
                    :disabled="isLoading" @click="handleStop(tunnel)"></v-btn>
                </template>
              </v-tooltip>
              <v-tooltip location="top" text="Restart Tunnel">
                <template #activator="{ props }">
                  <v-btn v-bind="props" icon="mdi-refresh" color="info" variant="text" size="small"
                    :disabled="isLoading" @click="handleRestart(tunnel)"></v-btn>
                </template>
              </v-tooltip>
              <v-tooltip location="top" text="Delete Tunnel">
                <template #activator="{ props }">
                  <v-btn v-bind="props" icon="mdi-delete" color="error" variant="text" size="small"
                    :disabled="isLoading" @click="openDeleteDialog(tunnel)"></v-btn>
                </template>
              </v-tooltip>
            </div>
          </template>
        </v-expansion-panel-title>

        <!-- Expansion Panel Content (DNS Records) -->
        <v-expansion-panel-text>
          <v-toolbar dense flat color="transparent">
            <v-toolbar-title class="text-subtitle-1">DNS Records</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn color="primary" variant="text" size="small" prepend-icon="mdi-plus"
              @click="openCreateDnsDialog(tunnel)" :disabled="isLoading">
              Add Domain
            </v-btn>
          </v-toolbar>
          <v-divider></v-divider>

          <v-list v-if="tunnel.dnsRecords && tunnel.dnsRecords.length > 0" lines="two">
            <v-list-item v-for="dns in tunnel.dnsRecords" :key="dns.id" :title="dns.name"
              :subtitle="`Type: ${dns.type}  |  Content: ${dns.content}`" prepend-icon="mdi-dns">
              <template #append>
                <v-chip :color="dns.proxied ? 'orange' : 'grey'" size="small">
                  {{ dns.proxied ? 'Proxied' : 'DNS Only' }}
                </v-chip>
                <!-- Add a delete button for DNS records if your API supports it -->
              </template>
            </v-list-item>
          </v-list>
          <v-alert v-else type="info" variant="tonal" class="mt-4" dense>
            This tunnel has no DNS records routed to it yet.
          </v-alert>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>

    <!-- --- Dialogs --- -->

    <!-- Create Tunnel Dialog -->
    <v-dialog v-model="createDialog" max-width="500px">
      <v-card>
        <v-card-title>Create New Tunnel</v-card-title>
        <v-card-text>
          <v-text-field v-model="newTunnelName" label="Tunnel Name" placeholder="my-cool-tunnel" variant="outlined"
            autofocus @keydown.enter="handleCreateConfirm"></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="createDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="isLoading" :disabled="isLoading || !newTunnelName"
            @click="handleCreateConfirm">
            Create
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Tunnel Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500px">
      <v-card>
        <v-card-title>Delete Tunnel</v-card-title>
        <v-card-text>
          Are you sure you want to delete
          <strong>{{ selectedTunnel?.name }}</strong>? This will also remove any
          associated DNS records. This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="isLoading" :disabled="isLoading" @click="handleDeleteConfirm">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Add DNS Record Dialog -->
    <v-dialog v-model="createDnsDialog" max-width="500px">
      <v-card>
        <v-card-title>Add DNS Record</v-card-title>
        <v-card-subtitle>
          Route a new domain to <strong>{{ selectedTunnel?.name }}</strong>
        </v-card-subtitle>
        <v-card-text>
          <v-text-field v-model="newDnsDomain" label="Domain Name" placeholder="subdomain.example.com"
            variant="outlined" autofocus @keydown.enter="handleCreateDnsConfirm"></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="createDnsDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="isLoading" :disabled="isLoading || !newDnsDomain"
            @click="handleCreateDnsConfirm">
            Add Route
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
// Import all functions from your tunnel API file
import * as tunnelApi from '@/api/tunnel';
// Import the type definition
import type { tunnel } from '@/models/tunnel';

// State
const tunnels = ref<tunnel[] | null>(null);
const isLoading = ref(false);
const selectedTunnel = ref<tunnel | null>(null);

// Dialog State
const createDialog = ref(false);
const newTunnelName = ref("");

const deleteDialog = ref(false);

const createDnsDialog = ref(false);
const newDnsDomain = ref("");

// --- Data Fetching ---

async function loadTunnels() {
  isLoading.value = true;
  try {
    tunnels.value = await tunnelApi.getTunnels() ?? [];
  } catch (err) {
    console.error("Failed to load tunnels:", err);
    tunnels.value = []; // Set to empty array on error to stop loading
  } finally {
    isLoading.value = false;
  }
}

// Load tunnels when the component is mounted
onMounted(loadTunnels);

// --- Tunnel Actions ---

async function handleStart(tunnel: tunnel) {
  isLoading.value = true;
  console.log(`Starting tunnel ${tunnel.name}...`);
  await tunnelApi.startTunnel(tunnel.id);
  // TODO: Add snackbar for success/error
  isLoading.value = false;
}

async function handleStop(tunnel: tunnel) {
  isLoading.value = true;
  console.log(`Stopping tunnel ${tunnel.name}...`);
  await tunnelApi.stopTunnel(tunnel.id);
  isLoading.value = false;
}

async function handleRestart(tunnel: tunnel) {
  isLoading.value = true;
  console.log(`Restarting tunnel ${tunnel.name}...`);
  await tunnelApi.restartTunnel(tunnel.id);
  isLoading.value = false;
}

// --- Create Tunnel ---

function openCreateDialog() {
  newTunnelName.value = "";
  createDialog.value = true;
}

async function handleCreateConfirm() {
  if (!newTunnelName.value) return;
  isLoading.value = true;
  await tunnelApi.createTunnel(newTunnelName.value);
  createDialog.value = false;
  // No need to set isLoading to false, as loadTunnels will do it
  await loadTunnels();
}

// --- Delete Tunnel ---

function openDeleteDialog(tunnel: tunnel) {
  selectedTunnel.value = tunnel;
  deleteDialog.value = true;
}

async function handleDeleteConfirm() {
  if (!selectedTunnel.value) return;
  isLoading.value = true;
  await tunnelApi.deleteTunnel(selectedTunnel.value.id);
  deleteDialog.value = false;
  await loadTunnels();
}

// --- Create DNS Record ---

function openCreateDnsDialog(tunnel: tunnel) {
  selectedTunnel.value = tunnel;
  newDnsDomain.value = "";
  createDnsDialog.value = true;
}

async function handleCreateDnsConfirm() {
  if (!selectedTunnel.value || !newDnsDomain.value) return;
  isLoading.value = true;
  await tunnelApi.createDnsRecord(selectedTunnel.value.id, newDnsDomain.value);
  createDnsDialog.value = false;
  await loadTunnels(); // Refresh to show new DNS record
}
</script>

<style scoped>
.tunnel-list .v-expansion-panel-title__icon {
  /* This ensures the control buttons don't get shifted by the icon */
  margin-left: auto;
}
</style>
