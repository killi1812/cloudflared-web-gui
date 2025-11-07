export interface tunnel {
  id: string
  name: string
  dnsRecords: dnsRecord[]
  created_at: string
  deleted_at: string
}

export interface dnsRecord {
  id: string
  name: string
  type: string
  content: string
  proxiable: boolean
  proxied: boolean
  ttl: number
  settings: any
  meta: any
  commnet: string
  tags: any[]
  created_at: string
  modified_on: string
}
