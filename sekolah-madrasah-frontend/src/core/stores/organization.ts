import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import type { Organization } from '../../business/types';

const STORAGE_KEY = 'selected_organization';

// Create writable store for selected organization
function createOrganizationStore() {
    // Initialize from localStorage if in browser
    const storedOrg = browser ? localStorage.getItem(STORAGE_KEY) : null;
    const initial: Organization | null = storedOrg ? JSON.parse(storedOrg) : null;

    const { subscribe, set, update } = writable<Organization | null>(initial);

    return {
        subscribe,
        set: (org: Organization | null) => {
            if (browser && org) {
                localStorage.setItem(STORAGE_KEY, JSON.stringify(org));
            } else if (browser && !org) {
                localStorage.removeItem(STORAGE_KEY);
            }
            set(org);
        },
        clear: () => {
            if (browser) {
                localStorage.removeItem(STORAGE_KEY);
            }
            set(null);
        },
    };
}

// Selected organization store
export const selectedOrganization = createOrganizationStore();

// Derived store for organization ID (for convenience)
export const selectedOrgId = derived(
    selectedOrganization,
    ($org) => $org?.id || null
);
