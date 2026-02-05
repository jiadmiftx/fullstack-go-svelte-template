import { writable, derived, get } from 'svelte/store';
import { api } from '../api/client';
import { browser } from '$app/environment';
import type { User } from '../types';

// Re-export User type for convenience
export type { User } from '../types';

interface AuthState {
    user: User | null;
    accessToken: string | null;
    refreshToken: string | null;
    isLoading: boolean;
}

// Initial state
const initialState: AuthState = {
    user: null,
    accessToken: null,
    refreshToken: null,
    isLoading: true,
};

// Create store
function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>(initialState);

    // Load from localStorage on init
    if (browser) {
        const stored = localStorage.getItem('auth');
        if (stored) {
            try {
                const parsed = JSON.parse(stored);
                set({ ...parsed, isLoading: false });
                if (parsed.accessToken) {
                    api.setToken(parsed.accessToken);
                }
            } catch {
                set({ ...initialState, isLoading: false });
            }
        } else {
            update(s => ({ ...s, isLoading: false }));
        }
    }

    return {
        subscribe,

        async login(email: string, password: string) {
            const response = await api.post<{
                message: string;
                data: {
                    access_token: string;
                    refresh_token: string;
                    expires_at: number;
                    user: User;
                };
            }>('/auth/login', { email, password });

            const authData = {
                user: response.data.user,
                accessToken: response.data.access_token,
                refreshToken: response.data.refresh_token,
                isLoading: false,
            };

            api.setToken(authData.accessToken);

            if (browser) {
                localStorage.setItem('auth', JSON.stringify(authData));
            }

            set(authData);
            return response;
        },

        async register(email: string, password: string, fullName: string) {
            const response = await api.post<{
                message: string;
                data: { user: User };
            }>('/auth/register', { email, password, full_name: fullName });

            return response;
        },

        logout() {
            api.setToken(null);
            if (browser) {
                localStorage.removeItem('auth');
            }
            set({ ...initialState, isLoading: false });
        },
    };
}

export const auth = createAuthStore();

// Derived stores
export const isAuthenticated = derived(auth, $auth => !!$auth.accessToken);
export const currentUser = derived(auth, $auth => $auth.user);
export const isSuperAdmin = derived(auth, $auth => $auth.user?.is_super_admin ?? false);
