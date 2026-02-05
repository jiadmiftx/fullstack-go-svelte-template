<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import {
        isAuthenticated,
        api,
        selectedOrganization,
        selectedOrgId,
        type User,
        type Organization,
        type Role,
    } from "$lib";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    let users: User[] = [];
    let organizations: Organization[] = [];
    let roles: Role[] = [];
    let userOrganizations: Map<string, Organization[]> = new Map();
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingUser: User | null = null;

    // Form fields
    let formEmail = "";
    let formPassword = "";
    let formFullName = "";
    let formPhone = "";
    let formAvatar = "";
    let formOrganizationId = "";
    let formRoleId = "";

    $: if (!$isAuthenticated && typeof window !== "undefined") {
        goto("/auth/login");
    }

    $: if (formOrganizationId) {
        loadRolesByOrg(formOrganizationId);
    }

    onMount(async () => {
        if ($isAuthenticated) {
            // Load organizations first, then users (to get memberships)
            await loadOrganizations();
            await loadUsers();
        }
    });

    async function loadUsers() {
        isLoading = true;
        error = "";
        try {
            // Only fetch platform-level users (super admins and unassigned users)
            const response = await api.getUsers({
                limit: 100,
                platform_only: true,
            });
            users = response.data || [];
            // After loading users, load their org memberships
            await loadUserMemberships();
        } catch (err) {
            error = err instanceof Error ? err.message : "Failed to load users";
        } finally {
            isLoading = false;
        }
    }

    async function loadUserMemberships() {
        // Clear existing mappings
        userOrganizations = new Map();

        // For each org, get members and map to users
        for (const org of organizations) {
            try {
                const response = await api.getOrganizationMembers(org.id, {
                    limit: 100,
                });
                const members = response.data || [];

                for (const member of members) {
                    const userId = member.user_id;
                    if (!userOrganizations.has(userId)) {
                        userOrganizations.set(userId, []);
                    }
                    userOrganizations.get(userId)!.push(org);
                }
            } catch (err) {
                console.error(`Failed to load members for org ${org.id}:`, err);
            }
        }
        // Trigger reactivity
        userOrganizations = userOrganizations;
    }

    function getUserOrgs(userId: string): string {
        const orgs = userOrganizations.get(userId);
        if (!orgs || orgs.length === 0) return "-";
        return orgs.map((o) => o.name).join(", ");
    }

    async function loadOrganizations() {
        try {
            const response = await api.getOrganizations({ limit: 100 });
            organizations = response.data || [];
        } catch (err) {
            console.error("Failed to load organizations:", err);
        }
    }

    async function loadRolesByOrg(orgId: string) {
        try {
            const response = await api.getRoles({
                organization_id: orgId,
                limit: 100,
            });
            roles = response.data || [];
        } catch (err) {
            console.error("Failed to load roles:", err);
            roles = [];
        }
    }

    function openCreateModal() {
        editingUser = null;
        formEmail = "";
        formPassword = "";
        formFullName = "";
        formPhone = "";
        formOrganizationId = $selectedOrgId || "";
        formRoleId = "";
        showModal = true;
    }

    function openEditModal(user: User) {
        editingUser = user;
        formEmail = user.email;
        formPassword = "";
        formFullName = user.full_name;
        formPhone = user.phone || "";
        formAvatar = user.avatar || "";
        formOrganizationId = "";
        formRoleId = "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingUser) {
                const updateData: any = {
                    email: formEmail,
                    full_name: formFullName,
                    phone: formPhone || undefined,
                    avatar: formAvatar || undefined,
                };
                // Password cannot be updated via this endpoint
                await api.updateUser(editingUser.id, updateData);

                // Update org/role assignment if provided
                if (formOrganizationId && formRoleId) {
                    try {
                        await api.addMemberToOrganization(formOrganizationId, {
                            user_id: editingUser.id,
                            role_id: formRoleId,
                        });
                        showToast(
                            "User updated and assigned to organization",
                            "success",
                        );
                    } catch (assignErr) {
                        // Might already be member, try update instead
                        try {
                            await api.updateMemberRole(
                                formOrganizationId,
                                editingUser.id,
                                {
                                    role_id: formRoleId,
                                },
                            );
                            showToast(
                                "User and role updated successfully",
                                "success",
                            );
                        } catch {
                            showToast(
                                "User updated but failed to update organization",
                                "error",
                            );
                        }
                    }
                } else {
                    showToast("User updated successfully", "success");
                }
            } else {
                if (!formPassword) {
                    error = "Password is required for new users";
                    showToast("Password is required", "error");
                    return;
                }

                // Step 1: Create the user
                const createResponse = await api.createUser({
                    email: formEmail,
                    password: formPassword,
                    full_name: formFullName,
                    phone: formPhone || undefined,
                });

                const newUser = createResponse.data;

                // Step 2: If organization and role selected, add to org
                if (formOrganizationId && formRoleId && newUser?.id) {
                    try {
                        await api.addMemberToOrganization(formOrganizationId, {
                            user_id: newUser.id,
                            role_id: formRoleId,
                        });
                        showToast(
                            "User created and assigned to organization",
                            "success",
                        );
                    } catch (assignErr) {
                        // User created but org assignment failed
                        showToast(
                            "User created but failed to assign to organization",
                            "error",
                        );
                    }
                } else {
                    showToast("User created successfully", "success");
                }
            }
            showModal = false;
            await loadUsers();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to save user";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(user: User) {
        if (!confirm(`Are you sure you want to delete "${user.full_name}"?`))
            return;

        error = "";
        try {
            await api.deleteUser(user.id);
            showToast("User deleted successfully", "success");
            await loadUsers();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to delete user";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    const columns = [
        { key: "full_name" as keyof User, label: "Name", sortable: true },
        { key: "email" as keyof User, label: "Email", sortable: true },
        {
            key: "id" as keyof User,
            label: "Organizations",
            render: (_: string, row: User) => {
                const orgs = userOrganizations.get(row.id);
                if (!orgs || orgs.length === 0)
                    return '<span class="text-slate-500">-</span>';
                return orgs
                    .map(
                        (o) =>
                            `<span class="bg-gray-500/20 text-gray-600 px-2 py-0.5 rounded text-xs mr-1">${o.name}</span>`,
                    )
                    .join("");
            },
        },
        {
            key: "phone" as keyof User,
            label: "Phone",
            render: (phone: string) => phone || "-",
        },
        {
            key: "is_active" as keyof User,
            label: "Status",
            render: (active: boolean) =>
                active
                    ? '<span class="bg-green-500/20 text-green-400 px-2 py-1 rounded text-xs">Active</span>'
                    : '<span class="bg-red-500/20 text-red-400 px-2 py-1 rounded text-xs">Inactive</span>',
        },
        {
            key: "created_at" as keyof User,
            label: "Created",
            sortable: true,
            render: (value: string) => new Date(value).toLocaleDateString(),
        },
    ];

    const actions = [
        {
            label: "Edit",
            variant: "primary" as const,
            onClick: openEditModal,
        },
        {
            label: "Delete",
            variant: "danger" as const,
            onClick: handleDelete,
            // Hide delete button for super admin users
            isHidden: (user: User) => user.is_super_admin,
        },
    ];
</script>

<svelte:head>
    <title>Users - Multi Tenant RBAC</title>
</svelte:head>

{#if $isAuthenticated}
    <div class="py-8 px-4">
        <div class="max-w-7xl mx-auto">
            <!-- Header -->
            <div class="flex justify-between items-center mb-8">
                <div>
                    <h1 class="text-3xl font-bold mb-2">Users</h1>
                    <p class="text-slate-400">Manage system users</p>
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

            <!-- Users Table -->
            <div class="card">
                {#if isLoading}
                    <div class="text-center py-12 text-slate-400">
                        Loading users...
                    </div>
                {:else}
                    <DataTable
                        data={users}
                        {columns}
                        {actions}
                        emptyMessage="No users found. Create your first one!"
                    />
                {/if}
            </div>
        </div>
    </div>

    <!-- Create/Edit Modal -->
    <Modal
        bind:isOpen={showModal}
        title={editingUser ? "Edit User" : "Create User"}
        size="lg"
    >
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
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
                        placeholder="e.g. John Doe"
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
                        placeholder="e.g. john@example.com"
                        required
                    />
                </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        for="password"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Password {editingUser ? "(Cannot be changed)" : "*"}
                    </label>
                    <input
                        id="password"
                        type="password"
                        bind:value={formPassword}
                        class="input-field"
                        placeholder={editingUser
                            ? "Password update not supported"
                            : "Enter password"}
                        required={!editingUser}
                        disabled={!!editingUser}
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
                        placeholder="e.g. +628123456789"
                    />
                </div>
            </div>

            {#if editingUser}
                <div>
                    <label
                        for="avatar"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Avatar URL
                    </label>
                    <input
                        id="avatar"
                        type="text"
                        bind:value={formAvatar}
                        class="input-field"
                        placeholder="https://example.com/avatar.jpg"
                    />
                </div>
            {/if}

            <!-- Organization & Role Assignment -->
            <div class="border-t border-slate-700 pt-4 mt-2">
                <h3 class="text-sm font-medium text-slate-300 mb-3">
                    Organization & Role Assignment
                </h3>

                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label
                            for="organization"
                            class="block text-sm font-medium text-slate-300 mb-1"
                        >
                            Organization
                        </label>
                        <select
                            id="organization"
                            bind:value={formOrganizationId}
                            class="input-field"
                        >
                            <option value="">Select organization...</option>
                            {#each organizations as org}
                                <option value={org.id}
                                    >{org.name} ({org.code})</option
                                >
                            {/each}
                        </select>
                        <p class="text-xs text-slate-500 mt-1">
                            {editingUser
                                ? "Change organization assignment"
                                : "Optional: Assign user to an organization"}
                        </p>
                    </div>

                    <div>
                        <label
                            for="role"
                            class="block text-sm font-medium text-slate-300 mb-1"
                        >
                            Role
                        </label>
                        <select
                            id="role"
                            bind:value={formRoleId}
                            class="input-field"
                            disabled={!formOrganizationId}
                        >
                            <option value="">Select role...</option>
                            {#each roles as role}
                                <option value={role.id}
                                    >{role.display_name}</option
                                >
                            {/each}
                        </select>
                        <p class="text-xs text-slate-500 mt-1">
                            {formOrganizationId
                                ? "Assign a role within the selected org"
                                : "Select organization first"}
                        </p>
                    </div>
                </div>
            </div>
        </form>

        <div slot="footer">
            <button on:click={() => (showModal = false)} class="btn-secondary">
                Cancel
            </button>
            <button on:click={handleSubmit} class="btn-primary">
                {editingUser ? "Update" : "Create"}
            </button>
        </div>
    </Modal>
{/if}
