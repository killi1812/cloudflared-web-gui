import type { NewUserDto, user } from "@/models/user";
import serverApi from "./serverAxios";




// --- API Calls ---

/**
 * Fetches the currently logged-in user's data.
 * @returns A promise resolving to the user's data.
 */
export async function getLoggedInUserData(): Promise<user | undefined> {
  try {
    const rez = await serverApi.get<user>("/user/my-data");
    return rez.data;
  } catch (error: any) {
    console.error("Failed to get logged-in user data:", error);
  }
}

/**
 * Fetches a specific user by their UUID.
 * @param uuid The UUID of the user.
 * @returns A promise resolving to the user's data.
 */
export async function getUser(uuid: string): Promise<user | undefined> {
  try {
    const rez = await serverApi.get<user>(`/user/${uuid}`);
    return rez.data;
  } catch (error: any) {
    console.error(`Failed to get user ${uuid}:`, error);
  }
}

/**
 * Creates a new user.
 * @param userData The data for the new user.
 * @returns A promise resolving to the newly created user's data.
 */
export async function createUser(userData: NewUserDto): Promise<user | undefined> {
  try {
    const rez = await serverApi.post<user>("/user", userData);
    return rez.data;
  } catch (error: any) {
    console.error("Failed to create user:", error);
  }
}

/**
 * Updates an existing user.
 * @param uuid The UUID of the user to update.
 * @param userData The user data to update.
 * @returns A promise resolving to the updated user's data.
 */
export async function updateUser(uuid: string, userData: user): Promise<user | undefined> {
  try {
    const rez = await serverApi.put<user>(`/user/${uuid}`, userData);
    return rez.data;
  } catch (error: any) {
    console.error(`Failed to update user ${uuid}:`, error);
  }
}

/**
 * Deletes a user by their UUID.
 * @param uuid The UUID of the user to delete.
 * @returns A promise resolving to true if deletion was successful.
 */
export async function deleteUser(uuid: string): Promise<boolean | undefined> {
  try {
    const rez = await serverApi.delete(`/user/${uuid}`);
    return rez.status === 204; // 204 No Content
  } catch (error: any) {
    console.error(`Failed to delete user ${uuid}:`, error);
  }
}

/**
 * Fetches all users (Super Admin only).
 * @returns A promise resolving to an array of all users.
 */
export async function getAllUsers(): Promise<user[] | undefined> {
  try {
    const rez = await serverApi.get<user[]>("/user/all-users");
    return rez.data;
  } catch (error: any) {
    console.error("Failed to get all users:", error);
  }
}

/**
 * Searches for users by name.
 * @param query The search query string.
 * @returns A promise resolving to an array of matching users.
 */
export async function searchUsersByName(query: string): Promise<user[] | undefined> {
  try {
    const rez = await serverApi.get<user[]>("/user/search", {
      params: { query }
    });
    return rez.data;
  } catch (error: any) {
    console.error(`Failed to search users with query "${query}":`, error);
  }
}

/**
 * Fetches a specific user by their OIB.
 * @param oib The OIB of the user.
 * @returns A promise resolving to the user's data.
 */
export async function getUserByOib(oib: string): Promise<user | undefined> {
  try {
    const rez = await serverApi.get<user>(`/user/oib/${oib}`);
    return rez.data;
  } catch (error: any) {
    console.error(`Failed to get user with OIB ${oib}:`, error);
  }
}
