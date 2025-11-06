export type userRole = "superadmin"

export interface user {
  uuid: string
  username: string
  role: userRole
}
