export type userRole = "superadmin"

/**
 * Represents a user's data.
 * Based on dto.UserDto.
 */
export interface user {
  uuid: string
  username: string
  role: userRole
}

/**
 * Data required for creating a new user.
 * Based on dto.NewUserDto.
 */
export interface NewUserDto {
  username: string;
  password: string;
  first_name: string;
  last_name: string;
  oib: string;
  role: string;
  // Add any other fields
}
