import type { TokenDto } from "@/models/tokenDto";
import serverApi from "./serverAxios";
import type { LoginDto } from "@/models/loginDto";

/**
 * Exchanges a Discord authorization code for an access token.
 * @param {string} code - The authorization code from the Discord SDK.
 * @returns {Promise<string | undefined>} A promise resolving to the access token.
 */
export async function authorizeServer(code: string): Promise<string | undefined> {
  try {
    const rez = await serverApi.post<TokenDto>("/token", {
      code: code,
    });
    console.log(rez.data);
    return rez.data.accessToken;
  } catch (error: any) {
    console.error(error);
  }
}

/**
 * Authenticates a user and returns an access token.
 * @param credentials The user's login credentials.
 * @returns A promise resolving to the token object.
 */
export async function login(credentials: LoginDto): Promise<TokenDto | undefined> {
  try {
    const rez = await serverApi.post<TokenDto>("/auth/login", credentials);
    return rez.data;
  } catch (error: any) {
    console.error("Login failed:", error);
  }
}


/**
 * Generates a new access token using a valid refresh token (sent via auth header).
 * @returns A promise resolving to a new token object.
 */
export async function refreshToken(): Promise<TokenDto | undefined> {
  try {
    // The refresh token is expected to be in the Authorization header,
    // which the serverApi interceptor should already be adding.
    const rez = await serverApi.post<TokenDto>("/auth/refresh");
    return rez.data;
  } catch (error: any) {
    console.error("Token refresh failed:", error);
  }
}

/**
 * Logs out the current user by invalidating their token.
 * @returns A promise resolving to true if logout was successful (HTTP 200).
 */
export async function logout(): Promise<boolean | undefined> {
  try {
    // Auth header is added by the interceptor
    const rez = await serverApi.post("/auth/logout");
    return rez.status === 200;
  } catch (error: any) {
    console.error("Logout failed:", error);
  }
}
