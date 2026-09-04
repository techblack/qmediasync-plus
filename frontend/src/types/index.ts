export interface APIResponse<T> { code: number; message: string; data: T }
export interface User { id: string; username: string; role?: string }
export interface Account { id:number; source_type:string; name:string; username:string; user_id:string; token?:string; base_url?:string; token_failed_reason?:string; cookie_bound?:boolean; app_id_name?:string }
export interface SyncPath { id:number; source_type:string; account_id:number; account_name:string; base_cid:string; remote_path:string; local_path:string; enable_cron:boolean; cron:string; is_running:number; custom_config:boolean; [key:string]:unknown }
export interface SyncRecord { id:number; remote_path:string; local_path:string; status:number; is_full_sync:boolean; new_strm:number; new_meta:number; new_upload:number; finish_at:number; fail_reason:string }
export interface ScrapePath { id:number; source_type:string; media_type:string; source_path:string; dest_path:string; scrape_type:string; enable_cron:boolean; cron_expression:string; is_running:number }
export interface QueueTask { id:number; file_name:string; local_full_path?:string; remote_path?:string; source_type?:string; source?:string; status:number; file_size?:number; size?:number; error?:string }
export interface BackupRecord { id:number; created_at:number; backup_type:string; status:string; file_size:number; backup_duration:number }
export interface NotificationChannel { id:number; channel_type:string; channel_name:string; description:string; is_enabled:boolean }
export interface APIKeyItem { id:number; name:string; key_prefix:string; is_active:boolean; last_used_at:number; created_at:number }
export interface FileItem { id:string; name:string; is_directory:boolean; size:number; modified_time:number|string }
export type AnyMap = Record<string, any>
