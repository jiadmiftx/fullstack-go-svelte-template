<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { api, selectedOrganization, type User, type Role } from "$lib";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    interface Member {
        id: string;
        user_id: string;
        organization_id: string;
        role_id: string;
        is_active: boolean;
        joined_at: string;
        user?: User;
        role?: Role;
    }

    let orgId: string = "";
    let members: Member[] = [];
    let users: User[] = [];
    let roles: Role[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;

    // Form fields for creating new user
    let formEmail = "";
    let formPassword = "";
    let formFullName = "";
    let formPhone = "";
    let formRoleId = "";

    $: orgId = $page.params.orgId ?? "";
    $: organization = $selectedOrganization;

    onMount(async () => {
        if (orgId) {
            await Promise.all([loadMembers(), loadRoles()]);
        }
    });

    async function loadMembers() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getOrganizationMembers(orgId, {
                limit: 100,
            });
            members = response.data || [];

            // Load user details for each member
            const userIds = members.map((m) => m.user_id);
            if (userIds.length > 0) {
                const usersResponse = await api.getUsers({ limit: 100 });
                users = usersResponse.data || [];
            }
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Failed to load members";
        } finally {
            isLoading = false;
        }
    }

    async function loadRoles() {
        try {
            const response = await api.getRoles({
                limit: 100,
                organization_id: orgId,
            });
            roles = response.data || [];
        } catch (err) {
            console.error("Failed to load roles:", err);
        }
    }

    function openCreateModal() {
        formEmail = "";
        formPassword = "";
        formFullName = "";
        formPhone = "";
        formRoleId = "";
        showModal = true;
    }

    async function handleCreateUser() {
        error = "";
        if (!formEmail || !formPassword || !formFullName || !formRoleId) {
            showToast("Please fill all required fields", "error");
            return;
        }

        try {
            // 1. Register the user
            const registerResponse = await api.register({
                email: formEmail,
                password: formPassword,
                full_name: formFullName,
                phone: formPhone,
            });

            const newUser = registerResponse.data.user;

            // 2. Add user as member to organization
            await api.addMemberToOrganization(orgId, {
                user_id: newUser.id,
                role_id: formRoleId,
            });

            showToast(
                `User ${formFullName} created and added to organization`,
                "success",
            );
            showModal = false;
            await loadMembers();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to create user";
            showToast(errorMsg, "error");
        }
    }

    async function handleRemoveMember(member: Member) {
        if (
            !confirm(
                "Are you sure you want to remove this member from the organization?",
            )
        )
            return;

        try {
            await api.removeMemberFromOrganization(orgId, member.user_id);
            showToast("Member removed successfully", "success");
            await loadMembers();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to remove member";
            showToast(errorMsg, "error");
        }
    }

    const columns = [
        {
            key: "user_id" as keyof Member,
            label: "User",
            render: (userId: string) => {
                const user = users.find((u) => u.id === userId);
                return user
                    ? `<div><div class="font-medium">${user.full_name}</div><div class="text-xs text-slate-400">${user.email}</div></div>`
                    : userId;
            },
        },
        {
            key: "role_id" as keyof Member,
            label: "Role",
            render: (roleId: string) => {
                const role = roles.find((r) => r.id === roleId);
                return role
                    ? `<span class="bg-gray-500/20 text-gray-600 px-2 py-1 rounded text-xs">${role.display_name || role.name}</span>`
                    : roleId;
            },
        },
        {
            key: "is_active" as keyof Member,
            label: "Status",
            render: (active: boolean) =>
                active
                    ? '<span class="bg-green-500/20 text-green-400 px-2 py-1 rounded text-xs">Active</span>'
                    : '<span class="bg-red-500/20 text-red-400 px-2 py-1 rounded text-xs">Inactive</span>',
        },
        {
            key: "joined_at" as keyof Member,
            label: "Joined",
            render: (value: string) => new Date(value).toLocaleDateString(),
        },
    ];

    const actions = [
        {
            label: "Remove",
            variant: "danger" as const,
            onClick: handleRemoveMember,
        },
    ];
</script>

<svelte:head>
    <title
        >Users - {organization?.name || "Organization"} - Multi Tenant RBAC</title
    >
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-7xl mx-auto">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
            <div>
                <h1 class="text-3xl font-bold mb-2">Organization Users</h1>
                <p class="text-slate-400">
                    Manage users in {organization?.name || "this organization"}
                </p>
            </div>
            <button on:click={openCreateModal} class="btn-primary">
                + Create User
            </button>
        </div>

        <!-- Error Alert -->
        {#if error}
            <div
                class="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-6"
            >
                {error}
            </div>
        {/if}

        <!-- Members Table -->
        <div class="card">
            {#if isLoading}
                <div class="text-center py-12 text-slate-400">
                    Loading members...
                </div>
            {:else}
                <DataTable
                    data={members}
                    {columns}
                    {actions}
                    emptyMessage="No users in this organization yet. Create your first user!"
                />
            {/if}
        </div>
    </div>
</div>

<!-- Create User Modal -->
<Modal bind:isOpen={showModal} title="Create New User" size="md">
    <form on:submit|preventDefault={handleCreateUser} class="space-y-4">
        <div>
            <label
                for="fullName"
                class="block text-sm font-medium text-slate-300 mb-1"
            >
                Full Name *
            </label>
            <input
                id="fullName"
                type="text"
                bind:value={formFullName}
                class="input-field"
                placeholder="e.g. Budi Santoso"
                required
            />
        </div>

        <div>
            <label
                for="email"
                class="block text-sm font-medium text-slate-300 mb-1"
            >
                Email *
            </label>
            <input
                id="email"
                type="email"
                bind:value={formEmail}
                class="input-field"
                placeholder="e.g. budi@perumahan.id"
                required
            />
        </div>

        <div>
            <label
                for="phone"
                class="block text-sm font-medium text-slate-300 mb-1"
            >
                Phone
            </label>
            <input
                id="phone"
                type="tel"
                bind:value={formPhone}
                class="input-field"
                placeholder="e.g. +6281234567890"
            />
        </div>

        <div>
            <label
                for="password"
                class="block text-sm font-medium text-slate-300 mb-1"
            >
                Password *
            </label>
            <input
                id="password"
                type="password"
                bind:value={formPassword}
                class="input-field"
                placeholder="Minimum 8 characters"
                required
                minlength="8"
            />
            <p class="text-xs text-slate-500 mt-1">
                User can change this password after first login
            </p>
        </div>

        <div>
            <label
                for="role"
                class="block text-sm font-medium text-slate-300 mb-1"
            >
                Role *
            </label>
            <select
                id="role"
                bind:value={formRoleId}
                class="input-field"
                required
            >
                <option value="">Select Role</option>
                {#each roles as role}
                    <option value={role.id}
                        >{role.display_name || role.name}</option
                    >
                {/each}
            </select>
            {#if roles.length === 0}
                <p class="text-xs text-yellow-400 mt-1">
                    No roles available. <a
                        href="/org/{orgId}/roles"
                        class="underline">Create roles first</a
                    >.
                </p>
            {/if}
        </div>
    </form>

    <div slot="footer">
        <button on:click={() => (showModal = false)} class="btn-secondary">
            Cancel
        </button>
        <button on:click={handleCreateUser} class="btn-primary">
            Create User
        </button>
    </div>
</Modal>
