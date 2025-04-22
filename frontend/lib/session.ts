import {api} from "@/lib/api";

export async function validateSession() {
    try {
        const response = await api.getMe();

        if (!response.user_id || !response.username || !response.role) {
            throw new Error('Session validation failed: Invalid Session');
        }

        return response;
    }catch (error) {
        console.error('Session validation failed:', error);
        throw new Error('Session validation failed');
    }
}