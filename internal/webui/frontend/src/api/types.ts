export interface Status {
  running: boolean
  start_time?: string
  scan_count: number
  found_count: number
  bought_count: number
  failed_count: number
  last_scan_time?: string
  last_error?: string
}

export interface PriceRange {
  min: number
  max: number
}

export interface MatchConfig {
  family_id: string
  type_id: string
  types: string[]
  price_range: PriceRange
  max_age: number
  max_cycles: number
  condition_min: number
  offer_types: string[]
  financing: string[]
  sort_by: string
}

export interface ActionConfig {
  auto_buy: boolean
  snatch: boolean
  max_bid_increment: number
  max_count?: number
}

export interface Rule {
  name: string
  enabled: boolean
  priority: number
  server_id: number
  auth_id: number
  company_name?: string
  match: MatchConfig
  action: ActionConfig
}

export interface AircraftFamily {
  id: string
  name: string
  types?: AircraftType[]
}

export interface AircraftType {
  id: string
  name: string
}

export interface AircraftData {
  families: AircraftFamily[]
  types: AircraftType[]
}