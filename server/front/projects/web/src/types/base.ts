export interface BaseInfo {
  createdAt: number
  updatedAt: number
}

export function NewBaseEmptyInfo(): BaseInfo {
  return {
    createdAt: 0,
    updatedAt: 0
  }
}
