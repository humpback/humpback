export interface ResponseSuccess {
  statusCode: number
  msg: string
}

export interface ResponseQuery<T> {
  total: number
  list: T[]
}
