import type { APIResponse } from '../types'

export class APIError extends Error { constructor(message:string, public status:number, public code?:number) { super(message) } }

export async function request<T>(path:string, options:RequestInit = {}):Promise<T> {
  const headers:HeadersInit = options.body instanceof FormData ? { ...(options.headers || {}) } : { 'Content-Type':'application/json', ...(options.headers || {}) }
  const response = await fetch(`/api${path}`, { credentials:'include', ...options, headers })
  const contentType = response.headers.get('content-type') || ''
  if (!contentType.includes('application/json')) throw new APIError('服务接口返回了非 JSON 数据', response.status)
  const result = await response.json() as APIResponse<T> | T
  if (!response.ok) {
    if (response.status === 401) { window.dispatchEvent(new Event('qms:unauthorized')); window.location.hash = '#/login' }
    throw new APIError((result as APIResponse<T>).message || '请求失败', response.status, (result as APIResponse<T>).code)
  }
  if (typeof result === 'object' && result !== null && 'code' in result) {
    const wrapped = result as APIResponse<T>
    if (wrapped.code === 401 || wrapped.code === 403) { window.dispatchEvent(new Event('qms:unauthorized')); window.location.hash = '#/login' }
    if (wrapped.code !== 0 && wrapped.code !== 200) throw new APIError(wrapped.message || '请求失败', response.status, wrapped.code)
    return wrapped.data
  }
  return result as T
}
export const get = <T>(path:string) => request<T>(path)
export const post = <T>(path:string, data:unknown) => request<T>(path, { method:'POST', body:JSON.stringify(data) })
export const put = <T>(path:string, data:unknown) => request<T>(path, { method:'PUT', body:JSON.stringify(data) })
export const del = <T>(path:string) => request<T>(path, { method:'DELETE' })
