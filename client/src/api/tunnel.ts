import type { tunnel } from "@/models/tunnel";
import serverApi from "./serverAxios";

/**
 * Get a list of all tunnels.
 * @returns A promise that resolves to an array of tunnels.
 */
export async function getTunnels(): Promise<tunnel[] | undefined> {
  try {
    const rez = await serverApi.get<tunnel[]>("/tunnel");
    return rez.data;
  } catch (error: any) {
    console.error("Error fetching tunnels:", error);
  }
}

/**
 * Get detailed information for a single tunnel.
 * @param id The UUID of the tunnel.
 * @returns A promise that resolves to a single tunnel DTO.
 */
export async function getTunnelInfo(id: string): Promise<tunnel | undefined> {
  try {
    const rez = await serverApi.get<tunnel>(`/tunnel/${id}`);
    return rez.data;
  } catch (error: any) {
    console.error(`Error fetching tunnel info for ${id}:`, error);
  }
}

/**
 * Creates a new tunnel with the given name.
 * @param name The name for the new tunnel.
 * @returns A promise that resolves to the newly created tunnel.
 */
export async function createTunnel(name: string): Promise<tunnel | undefined> {
  try {
    const rez = await serverApi.post<tunnel>("/tunnel", { name });
    return rez.data;
  } catch (error: any) {
    console.error("Error creating tunnel:", error);
  }
}

/**
 * Deletes a tunnel by its ID.
 * @param id The UUID of the tunnel to delete.
 * @returns A promise that resolves to true if deletion was successful (HTTP 204).
 */
export async function deleteTunnel(id: string): Promise<boolean | undefined> {
  try {
    const rez = await serverApi.delete(`/tunnel/${id}`);
    return rez.status === 204; // 204 No Content
  } catch (error: any) {
    console.error(`Error deleting tunnel ${id}:`, error);
  }
}

/**
 * Creates a new DNS record for a tunnel.
 * @param id The UUID of the tunnel.
 * @param domain The domain name to route to the tunnel.
 * @returns A promise that resolves to the updated tunnel model.
 */
export async function createDnsRecord(id: string, domain: string): Promise<tunnel | undefined> {
  try {
    const rez = await serverApi.post<tunnel>(`/tunnel/dns/${id}`, { domain });
    return rez.data;
  } catch (error: any) {
    console.error(`Error creating DNS record for tunnel ${id}:`, error);
  }
}

/**
 * Starts a tunnel service.
 * @param id The UUID of the tunnel to start.
 * @returns A promise that resolves to true if the start command was successful (HTTP 204).
 */
export async function startTunnel(id: string): Promise<boolean | undefined> {
  try {
    const rez = await serverApi.put(`/tunnel/${id}/start`);
    return rez.status === 204;
  } catch (error: any) {
    console.error(`Error starting tunnel ${id}:`, error);
  }
}

/**
 * Stops a tunnel service.
 * @param id The UUID of the tunnel to stop.
 * @returns A promise that resolves to true if the stop command was successful (HTTP 204).
 */
export async function stopTunnel(id: string): Promise<boolean | undefined> {
  try {
    const rez = await serverApi.put(`/tunnel/${id}/stop`);
    return rez.status === 204;
  } catch (error: any) {
    console.error(`Error stopping tunnel ${id}:`, error);
  }
}

/**
 * Restarts a tunnel service.
 * @param id The UUID of the tunnel to restart.
 * @returns A promise that resolves to true if the restart command was successful (HTTP 204).
 */
export async function restartTunnel(id: string): Promise<boolean | undefined> {
  try {
    const rez = await serverApi.put(`/tunnel/${id}/restart`);
    return rez.status === 204;
  } catch (error: any) {
    console.error(`Error restarting tunnel ${id}:`, error);
  }
}
