// Core User types
export interface User {
    id: string;
    email: string;
    full_name: string;
    phone?: string;
    avatar?: string;
    is_super_admin: boolean;
    is_active: boolean;
    last_login_at?: string;
    created_at: string;
    updated_at?: string;
}

// Role types
export interface Role {
    id: string;
    organization_id: string;
    name: string;
    display_name: string;
    type: string;
    description?: string;
    level: number;
    is_default: boolean;
    is_active: boolean;
    permissions?: Permission[];
    created_at: string;
    updated_at: string;
}

// Permission types
export interface Permission {
    id: string;
    name: string;
    resource: string;
    action: string;
    description?: string;
    created_at: string;
}

// Perumahan types
export interface Perumahan {
    id: string;
    owner_id: string;
    name: string;
    code: string;
    type: string;
    address?: string;
    phone?: string;
    email?: string;
    logo?: string;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
}

// API Response types
export interface ApiResponse<T> {
    message: string;
    data: T;
}

export interface PaginatedResponse<T> {
    data: T[];
    pagination: {
        page: number;
        limit: number;
        total: number;
        total_pages: number;
    };
}
