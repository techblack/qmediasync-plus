import { message } from 'ant-design-vue'
import { post } from './http'

export async function handleOAuthCallback() {
  const hashQuery = window.location.hash.includes('?') ? window.location.hash.split('?')[1] : ''
  const params = new URLSearchParams(window.location.search || hashQuery)
  const accountId = params.get('account_id'), tokenData = params.get('token_data'), source = params.get('source') || '115'
  if (!accountId || !tokenData) return false
  await post(`/${source === 'baidupan' ? 'baidupan' : '115'}/oauth-confirm`, { account_id: Number(accountId), data: tokenData })
  window.history.replaceState(null, '', `${window.location.pathname}#/accounts`)
  message.success('网盘授权成功')
  return true
}
