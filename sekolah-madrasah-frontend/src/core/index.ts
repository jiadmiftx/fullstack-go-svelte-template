// Core exports - Reusable infrastructure

// Components
export { default as DataTable } from './components/DataTable.svelte';
export { default as Toast } from './components/Toast.svelte';
export { default as StatsCard } from './components/StatsCard.svelte';
export { default as Modal } from './components/Modal.svelte';
export { default as RouteGuard } from './components/RouteGuard.svelte';

// API
export { api, API_BASE_URL } from './api/client';

// Stores
export {
    auth,
    isAuthenticated,
    currentUser,
    isSuperAdmin,
    type User
} from './stores/auth';
export { selectedOrganization, selectedOrgId } from './stores/organization';
export {
    userPermissions,
    hasPermission,
    hasAnyPermission,
    canViewPerumahans,
    canCreatePerumahan,
    canUpdatePerumahan,
    canDeletePerumahan
} from './stores/permissions';

// Types
export type {
    Role,
    Permission,
    Perumahan,
    ApiResponse,
    PaginatedResponse
} from './types';
