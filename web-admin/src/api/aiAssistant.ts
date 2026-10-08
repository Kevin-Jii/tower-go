import { http, unwrap } from './http'

export interface AIConversation { id: number; user_id: number; store_id: number; title: string; updated_at: string }
export interface AIMessage { id: number; conversation_id: number; role: 'user' | 'assistant'; content: string; analysis_context?: string; provider?: string; model?: string; created_at: string }
export interface AIConfig { provider: string; model: string; base_url: string; enabled: boolean; system_prompt: string; api_key_configured: boolean; api_key_masked: string }
export interface AIConfigInput { provider: string; model: string; base_url: string; enabled: boolean; system_prompt: string; api_key: string }
const storeParams = (store_id: number) => ({ params: { store_id } })
export async function getAIConfig(): Promise<AIConfig> { return unwrap(await http.get('/ai-assistant/config')) }
export async function saveAIConfig(data: AIConfigInput): Promise<AIConfig> { return unwrap(await http.put('/ai-assistant/config', data)) }
export async function testAIConfig(data: AIConfigInput): Promise<{ message: string }> { return unwrap(await http.post('/ai-assistant/config/test', data)) }
export async function listAIConversations(store_id: number): Promise<AIConversation[]> { return unwrap(await http.get('/ai-assistant/conversations', storeParams(store_id))) }
export async function createAIConversation(store_id: number): Promise<AIConversation> { return unwrap(await http.post('/ai-assistant/conversations', {}, storeParams(store_id))) }
export async function getAIMessages(id: number, store_id: number): Promise<AIMessage[]> { return unwrap(await http.get(`/ai-assistant/conversations/${id}/messages`, storeParams(store_id))) }
export async function renameAIConversation(id: number, title: string, store_id: number): Promise<void> { unwrap(await http.put(`/ai-assistant/conversations/${id}`, { title }, storeParams(store_id))) }
export async function deleteAIConversation(id: number, store_id: number): Promise<void> { unwrap(await http.delete(`/ai-assistant/conversations/${id}`, storeParams(store_id))) }
export async function sendAIMessage(data: { conversation_id: number; store_id: number; message: string; start_date: string; end_date: string }): Promise<AIMessage> { return unwrap(await http.post('/ai-assistant/chat', data)) }
