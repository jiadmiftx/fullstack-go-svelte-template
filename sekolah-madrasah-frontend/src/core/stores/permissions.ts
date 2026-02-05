import { derived, writable } from 'svelte/store';
import { auth } from './auth';
import type { Permission } from '../types';

// User permissions store - will be populated after login
export const userPermissions = writable<Permission[]>([]);

// Helper function to check if user has a specific permission
export function hasPermission(permissions: Permission[], permissionName: string): boolean {
    return permissions.some(p => p.name === permissionName);
}

// Helper to check multiple permissions (OR logic)
export function hasAnyPermission(permissions: Permission[], permissionNames: string[]): boolean {
    return permissionNames.some(name => hasPermission(permissions, name));
}

// Helper to check multiple permissions (AND logic)
export function hasAllPermissions(permissions: Permission[], permissionNames: string[]): boolean {
    return permissionNames.every(name => hasPermission(permissions, name));
}

// Derived store to check if current user is super admin (has all permissions)
export const isSuperAdminUser = derived(auth, $auth => $auth.user?.is_super_admin ?? false);

// Can user access perumahans?
export const canViewPerumahans = derived(
    [userPermissions, isSuperAdminUser],
    ([$perms, $isSuperAdmin]) => $isSuperAdmin || hasAnyPermission($perms, ['perumahans.list', 'perumahans.read'])
);

export const canCreatePerumahan = derived(
    [userPermissions, isSuperAdminUser],
    ([$perms, $isSuperAdmin]) => $isSuperAdmin || hasPermission($perms, 'perumahans.create')
);

export const canUpdatePerumahan = derived(
    [userPermissions, isSuperAdminUser],
    ([$perms, $isSuperAdmin]) => $isSuperAdmin || hasPermission($perms, 'perumahans.update')
);

export const canDeletePerumahan = derived(
    [userPermissions, isSuperAdminUser],
    ([$perms, $isSuperAdmin]) => $isSuperAdmin || hasPermission($perms, 'perumahans.delete')
);
