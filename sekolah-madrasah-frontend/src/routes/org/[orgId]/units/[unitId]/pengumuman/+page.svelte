<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";

    // Get IDs from URL
    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    // State
    let posts: any[] = [];
    let isLoading = true;
    let error = "";
    let showCreateModal = false;
    let selectedPost: any = null;

    // Role-based access
    let canCreatePost = false;
    let userRole = "";
    let userJabatanName = "";

    // For scope selection
    let perumahanName = "";
    let userRtId: string | null = null;
    let userRtName: string | null = null;

    // Post scope selection: 'org' = organization wide, 'unit' = unit only
    let postScope: "org" | "unit" = "unit";

    // Allowed jabatan names for creating posts
    const allowedJabatans = [
        "ketua rt",
        "wakil ketua rt",
        "sekretaris",
        "bendahara",
    ];

    // Pagination
    let currentPage = 1;
    let totalPages = 1;
    let pageSize = 10;

    // Create post form
    let newPost = {
        title: "",
        content: "",
        post_type: "text",
        image_url: "",
        link_url: "",
        link_title: "",
        is_pinned: false,
        is_important: false,
        poll_options: ["", ""],
    };

    // Comment form
    let newComment = "";
    let replyingTo: any = null;
    let replyContent = "";

    async function checkUserRole() {
        try {
            // Check user membership
            const membershipRes: any = await api.get("/users/me/memberships");
            console.log("Pengumuman: Membership Response:", membershipRes);

            const memberships = membershipRes.data?.unit_memberships || [];
            console.log("Pengumuman: Perumahan Memberships:", memberships);

            const membership = memberships.find(
                (m: any) => m.unit_id === unitId,
            );
            if (membership) {
                userRole = membership.role;
                console.log("Pengumuman: User Role:", userRole);

                // Admin and staff can always create posts
                if (userRole === "admin" || userRole === "staff") {
                    canCreatePost = true;
                    return;
                }
            }

            // Get unit_member_id from membership
            let perumahanMemberId: string | null = null;
            if (membership && membership.unit_member_id) {
                perumahanMemberId = membership.unit_member_id;
            }

            if (!perumahanMemberId) {
                console.log("Pengumuman: No perumahan member found");
                return;
            }
            console.log("Pengumuman: Perumahan Member ID:", perumahanMemberId);

            // Get user's warga profile by unit_member_id
            const profileRes: any = await api.get(
                `/warga-profiles?unit_member_id=${perumahanMemberId}&limit=1`,
            );
            if (!profileRes.data || profileRes.data.length === 0) {
                console.log("Pengumuman: No warga profile found");
                return;
            }
            const wargaProfileId = profileRes.data[0].id;
            userRtId = profileRes.data[0].rt_id || null;
            // Get RT name from the profile or fetch separately
            if (profileRes.data[0].rt_name) {
                userRtName = profileRes.data[0].rt_name;
            } else if (userRtId) {
                // Try to get RT name from RTs endpoint
                try {
                    const rtRes: any = await api.get(`/rts/${userRtId}`);
                    if (rtRes.data) {
                        userRtName =
                            rtRes.data.name ||
                            rtRes.data.nama ||
                            `RT ${userRtId.slice(0, 4)}`;
                    }
                } catch (e) {
                    userRtName = "RT Saya";
                }
            }
            console.log(
                "Pengumuman: Warga Profile ID:",
                wargaProfileId,
                "RT:",
                userRtId,
                userRtName,
            );

            // Fetch struktur RT which includes member assignments
            const strukturRes: any = await api.get(`/units/${unitId}/struktur`);
            console.log("Pengumuman: Struktur Response:", strukturRes);

            if (strukturRes.data) {
                for (const jabatan of strukturRes.data) {
                    const jabatanName = jabatan.nama?.toLowerCase() || "";
                    console.log(
                        "Pengumuman: Checking jabatan:",
                        jabatanName,
                        "Members:",
                        jabatan.anggota,
                    );

                    // Check if this jabatan is in allowed list
                    const isAllowed = allowedJabatans.some((allowed) =>
                        jabatanName.includes(allowed),
                    );

                    if (
                        isAllowed &&
                        jabatan.anggota &&
                        jabatan.anggota.length > 0
                    ) {
                        // Check if current user is assigned to this jabatan
                        for (const anggota of jabatan.anggota) {
                            console.log(
                                "Pengumuman: Checking anggota:",
                                anggota.warga_profile_id,
                                "vs",
                                wargaProfileId,
                            );

                            if (anggota.warga_profile_id === wargaProfileId) {
                                userJabatanName = jabatan.nama || "";
                                canCreatePost = true;
                                console.log(
                                    "Pengumuman: Can create post! Jabatan:",
                                    userJabatanName,
                                );
                                return;
                            }
                        }
                    }
                }
            }
            console.log("Pengumuman: No matching jabatan found");
        } catch (err) {
            console.error("Failed to check user role:", err);
        }
    }

    async function loadPerumahanInfo() {
        try {
            // Get perumahan (which is actually RT in this schema)
            const perumahanRes: any = await api.get(`/units/${unitId}`);
            if (perumahanRes.data) {
                // The "perumahan" name is actually the RT name (e.g., "RT 01")
                // Only set userRtName if not already set from warga profile
                if (!userRtName) {
                    userRtName =
                        perumahanRes.data.name ||
                        perumahanRes.data.nama ||
                        "RT Saya";
                }
                // Don't overwrite userRtId if already set from warga profile
                // userRtId should come from warga profile's rt_id, not unitId
                console.log("RT Name (from perumahan):", userRtName);
            }

            // Get organization name (which is the actual perumahan name like "Perumahan Griya Sehati")
            const orgRes: any = await api.get(`/organizations/${orgId}`);
            if (orgRes.data) {
                perumahanName =
                    orgRes.data.name || orgRes.data.nama || "Seluruh Perumahan";
                console.log("Organization Name:", perumahanName);
            }
        } catch (err) {
            console.error("Failed to load perumahan info:", err);
            perumahanName = "Seluruh Perumahan";
            if (!userRtName) {
                userRtName = "RT Saya";
            }
        }
    }

    async function loadPosts() {
        try {
            isLoading = true;
            const response = await api.get(
                `/posts?unit_id=${unitId}&page=${currentPage}&page_size=${pageSize}`,
            );
            if (response.data) {
                posts = response.data;
            }
            if (response.paginate) {
                totalPages = Math.ceil(
                    response.paginate.total / response.paginate.page_size,
                );
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat pengumuman";
        } finally {
            isLoading = false;
        }
    }

    async function createPost() {
        try {
            const payload: any = {
                unit_id: unitId,
                is_org_wide: postScope === "org", // true for org level, false for unit level
                title: newPost.title,
                content: newPost.content,
                post_type: newPost.post_type,
                is_pinned: newPost.is_pinned,
                is_important: newPost.is_important,
            };

            if (newPost.post_type === "photo" && newPost.image_url) {
                payload.image_url = newPost.image_url;
            }
            if (newPost.post_type === "link" && newPost.link_url) {
                payload.link_url = newPost.link_url;
                payload.link_title = newPost.link_title;
            }
            if (newPost.post_type === "poll") {
                payload.poll_options = newPost.poll_options.filter(
                    (o) => o.trim() !== "",
                );
            }

            await api.post("/posts", payload);
            showCreateModal = false;
            postScope = "unit"; // Reset scope
            resetForm();
            loadPosts();
        } catch (err: any) {
            error = err.message || "Gagal membuat pengumuman";
        }
    }

    function resetForm() {
        newPost = {
            title: "",
            content: "",
            post_type: "text",
            image_url: "",
            link_url: "",
            link_title: "",
            is_pinned: false,
            is_important: false,
            poll_options: ["", ""],
        };
    }

    async function deletePost(postId: string) {
        if (!confirm("Apakah Anda yakin ingin menghapus pengumuman ini?"))
            return;
        try {
            await api.delete(`/posts/${postId}`);
            loadPosts();
            if (selectedPost?.id === postId) {
                selectedPost = null;
            }
        } catch (err: any) {
            error = err.message || "Gagal menghapus pengumuman";
        }
    }

    async function viewPost(post: any) {
        try {
            const response = await api.get(`/posts/${post.id}`);
            if (response.data) {
                selectedPost = response.data;
                // Load comments
                loadComments(post.id);
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat detail pengumuman";
        }
    }

    async function loadComments(postId: string) {
        try {
            const response = await api.get(`/posts/${postId}/comments`);
            if (response.data && selectedPost) {
                selectedPost.comments = response.data;
            }
        } catch (err: any) {
            console.error("Failed to load comments:", err);
        }
    }

    async function addComment() {
        if (!newComment.trim() || !selectedPost) return;
        try {
            await api.post(`/posts/${selectedPost.id}/comments`, {
                content: newComment,
            });
            newComment = "";
            loadComments(selectedPost.id);
            // Update comment count
            const postIndex = posts.findIndex((p) => p.id === selectedPost.id);
            if (postIndex >= 0) {
                posts[postIndex].comment_count++;
            }
        } catch (err: any) {
            error = err.message || "Gagal menambah komentar";
        }
    }

    async function addReply() {
        if (!replyContent.trim() || !replyingTo || !selectedPost) return;
        try {
            await api.post(`/posts/${selectedPost.id}/comments`, {
                content: replyContent,
                parent_id: replyingTo.id,
            });
            replyContent = "";
            replyingTo = null;
            loadComments(selectedPost.id);
        } catch (err: any) {
            error = err.message || "Gagal menambah balasan";
        }
    }

    async function votePoll(optionId: string) {
        if (!selectedPost) return;
        try {
            await api.post(`/posts/${selectedPost.id}/vote`, {
                option_id: optionId,
            });
            // Reload post to get updated vote counts
            viewPost(selectedPost);
        } catch (err: any) {
            error = err.message || "Gagal melakukan voting";
        }
    }

    function addPollOption() {
        newPost.poll_options = [...newPost.poll_options, ""];
    }

    function removePollOption(index: number) {
        if (newPost.poll_options.length > 2) {
            newPost.poll_options = newPost.poll_options.filter(
                (_, i) => i !== index,
            );
        }
    }

    function formatDate(dateString: string) {
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "short",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function getPostTypeIcon(type: string) {
        switch (type) {
            case "photo":
                return "🖼️";
            case "poll":
                return "📊";
            case "link":
                return "🔗";
            default:
                return "📝";
        }
    }

    onMount(async () => {
        // checkUserRole first to get userRtId from warga profile
        await checkUserRole();
        // Then loadPerumahanInfo (won't overwrite userRtId if already set)
        await loadPerumahanInfo();
        loadPosts();
    });
</script>

<svelte:head>
    <title>Pengumuman - Perumahan</title>
</svelte:head>

<div class="pengumuman-container">
    <div class="page-header">
        <div class="header-content">
            <h1>📢 Pengumuman</h1>
            <p class="subtitle">
                Informasi dan pengumuman untuk seluruh warga perumahan
            </p>
        </div>
        {#if canCreatePost}
            <button
                class="btn-create"
                on:click={() => (showCreateModal = true)}
            >
                <span class="icon">+</span>
                Buat Pengumuman
            </button>
        {/if}
    </div>

    {#if error}
        <div class="error-message">
            {error}
            <button on:click={() => (error = "")}>×</button>
        </div>
    {/if}

    <div class="content-wrapper">
        <!-- Posts List -->
        <div class="posts-list">
            {#if isLoading}
                <div class="loading-state">
                    <div class="loader"></div>
                    <p>Memuat pengumuman...</p>
                </div>
            {:else if posts.length === 0}
                <div class="empty-state">
                    <div class="icon">📭</div>
                    <h3>Belum Ada Pengumuman</h3>
                    <p>Buat pengumuman pertama untuk warga perumahan</p>
                </div>
            {:else}
                {#each posts as post}
                    <div
                        class="post-card"
                        class:pinned={post.is_pinned}
                        class:important={post.is_important}
                        class:selected={selectedPost?.id === post.id}
                        on:click={() => viewPost(post)}
                        on:keypress={() => viewPost(post)}
                        role="button"
                        tabindex="0"
                    >
                        <div class="post-header">
                            <div class="post-meta">
                                <span class="post-type-icon"
                                    >{getPostTypeIcon(post.post_type)}</span
                                >
                                {#if post.is_pinned}
                                    <span class="badge pinned"
                                        >📌 Disematkan</span
                                    >
                                {/if}
                                {#if post.is_important}
                                    <span class="badge important"
                                        >⚠️ Penting</span
                                    >
                                {/if}
                                {#if post.is_org_wide}
                                    <span class="badge unit">Organisasi</span>
                                {:else}
                                    <span class="badge rt">Unit</span>
                                {/if}
                            </div>
                            <button
                                class="btn-delete"
                                on:click|stopPropagation={() =>
                                    deletePost(post.id)}
                                title="Hapus"
                            >
                                🗑️
                            </button>
                        </div>

                        {#if post.title}
                            <h3 class="post-title">{post.title}</h3>
                        {/if}

                        <p class="post-content">
                            {post.content.length > 200
                                ? post.content.substring(0, 200) + "..."
                                : post.content}
                        </p>

                        {#if post.post_type === "photo" && post.image_url}
                            <div class="post-image-preview">
                                <img src={post.image_url} alt="Preview" />
                            </div>
                        {/if}

                        <div class="post-footer">
                            <span class="author">👤 {post.author_name}</span>
                            <span class="date"
                                >{formatDate(post.created_at)}</span
                            >
                            <span class="comments">💬 {post.comment_count}</span
                            >
                        </div>
                    </div>
                {/each}

                <!-- Pagination -->
                {#if totalPages > 1}
                    <div class="pagination">
                        <button
                            disabled={currentPage === 1}
                            on:click={() => {
                                currentPage--;
                                loadPosts();
                            }}
                        >
                            ← Sebelumnya
                        </button>
                        <span>Halaman {currentPage} dari {totalPages}</span>
                        <button
                            disabled={currentPage === totalPages}
                            on:click={() => {
                                currentPage++;
                                loadPosts();
                            }}
                        >
                            Selanjutnya →
                        </button>
                    </div>
                {/if}
            {/if}
        </div>

        <!-- Post Detail -->
        {#if selectedPost}
            <div class="post-detail">
                <button
                    class="btn-close"
                    on:click={() => (selectedPost = null)}
                >
                    ×
                </button>

                <div class="detail-header">
                    <div class="badges">
                        {#if selectedPost.is_pinned}
                            <span class="badge pinned">📌 Disematkan</span>
                        {/if}
                        {#if selectedPost.is_important}
                            <span class="badge important">⚠️ Penting</span>
                        {/if}
                    </div>
                    {#if selectedPost.title}
                        <h2>{selectedPost.title}</h2>
                    {/if}
                    <div class="author-info">
                        <span class="author">👤 {selectedPost.author_name}</span
                        >
                        <span class="date"
                            >{formatDate(selectedPost.created_at)}</span
                        >
                    </div>
                </div>

                <div class="detail-content">
                    <p>{selectedPost.content}</p>

                    {#if selectedPost.post_type === "photo" && selectedPost.image_url}
                        <div class="post-image">
                            <img
                                src={selectedPost.image_url}
                                alt="Post image"
                            />
                        </div>
                    {/if}

                    {#if selectedPost.post_type === "link" && selectedPost.link_url}
                        <a
                            href={selectedPost.link_url}
                            target="_blank"
                            rel="noopener"
                            class="post-link"
                        >
                            🔗 {selectedPost.link_title ||
                                selectedPost.link_url}
                        </a>
                    {/if}

                    {#if selectedPost.post_type === "poll" && selectedPost.options}
                        <div class="poll-options">
                            <h4>📊 Polling</h4>
                            {#each selectedPost.options as option}
                                <button
                                    class="poll-option"
                                    class:voted={option.has_voted}
                                    on:click={() => votePoll(option.id)}
                                >
                                    <span class="option-text"
                                        >{option.text}</span
                                    >
                                    <span class="vote-count"
                                        >{option.vote_count} suara</span
                                    >
                                </button>
                            {/each}
                        </div>
                    {/if}
                </div>

                <!-- Comments Section -->
                <div class="comments-section">
                    <h4>💬 Komentar ({selectedPost.comments?.length || 0})</h4>

                    <div class="comment-form">
                        <textarea
                            bind:value={newComment}
                            placeholder="Tulis komentar..."
                            rows="2"
                        ></textarea>
                        <button
                            class="btn-submit"
                            on:click={addComment}
                            disabled={!newComment.trim()}
                        >
                            Kirim
                        </button>
                    </div>

                    <div class="comments-list">
                        {#if selectedPost.comments}
                            {#each selectedPost.comments as comment}
                                <div class="comment">
                                    <div class="comment-header">
                                        <span class="author"
                                            >👤 {comment.author_name}</span
                                        >
                                        <span class="date"
                                            >{formatDate(
                                                comment.created_at,
                                            )}</span
                                        >
                                    </div>
                                    <p class="comment-content">
                                        {comment.content}
                                    </p>

                                    <button
                                        class="btn-reply"
                                        on:click={() => (replyingTo = comment)}
                                    >
                                        Balas ({comment.reply_count || 0})
                                    </button>

                                    <!-- Reply form -->
                                    {#if replyingTo?.id === comment.id}
                                        <div class="reply-form">
                                            <textarea
                                                bind:value={replyContent}
                                                placeholder="Tulis balasan..."
                                                rows="2"
                                            ></textarea>
                                            <div class="reply-actions">
                                                <button
                                                    class="btn-cancel"
                                                    on:click={() => {
                                                        replyingTo = null;
                                                        replyContent = "";
                                                    }}
                                                >
                                                    Batal
                                                </button>
                                                <button
                                                    class="btn-submit"
                                                    on:click={addReply}
                                                    disabled={!replyContent.trim()}
                                                >
                                                    Kirim
                                                </button>
                                            </div>
                                        </div>
                                    {/if}

                                    <!-- Replies -->
                                    {#if comment.replies && comment.replies.length > 0}
                                        <div class="replies">
                                            {#each comment.replies as reply}
                                                <div class="reply">
                                                    <div class="comment-header">
                                                        <span class="author"
                                                            >👤 {reply.author_name}</span
                                                        >
                                                        <span class="date"
                                                            >{formatDate(
                                                                reply.created_at,
                                                            )}</span
                                                        >
                                                    </div>
                                                    <p class="comment-content">
                                                        {reply.content}
                                                    </p>
                                                </div>
                                            {/each}
                                        </div>
                                    {/if}
                                </div>
                            {/each}
                        {/if}
                    </div>
                </div>
            </div>
        {/if}
    </div>
</div>

<!-- Create Post Modal -->
{#if showCreateModal}
    <div
        class="modal-overlay"
        on:click={() => (showCreateModal = false)}
        on:keypress={() => (showCreateModal = false)}
        role="button"
        tabindex="0"
    >
        <div
            class="modal"
            on:click|stopPropagation
            role="dialog"
            on:keypress|stopPropagation
        >
            <div class="modal-header">
                <h2>Buat Pengumuman Baru</h2>
                <button
                    class="btn-close"
                    on:click={() => (showCreateModal = false)}>×</button
                >
            </div>

            <div class="modal-body">
                <div class="form-group">
                    <label>Tipe Postingan</label>
                    <div class="post-type-selector">
                        <button
                            class:active={newPost.post_type === "text"}
                            on:click={() => (newPost.post_type = "text")}
                        >
                            📝 Teks
                        </button>
                        <button
                            class:active={newPost.post_type === "photo"}
                            on:click={() => (newPost.post_type = "photo")}
                        >
                            🖼️ Foto
                        </button>
                        <button
                            class:active={newPost.post_type === "poll"}
                            on:click={() => (newPost.post_type = "poll")}
                        >
                            📊 Polling
                        </button>
                        <button
                            class:active={newPost.post_type === "link"}
                            on:click={() => (newPost.post_type = "link")}
                        >
                            🔗 Link
                        </button>
                    </div>
                </div>

                <!-- Scope Selector: Organization wide or Unit specific -->
                <div class="form-group">
                    <label>Target Pengumuman</label>
                    <div class="post-type-selector">
                        <button
                            class:active={postScope === "org"}
                            on:click={() => (postScope = "org")}
                        >
                            🏢 Seluruh Organisasi
                        </button>
                        <button
                            class:active={postScope === "unit"}
                            on:click={() => (postScope = "unit")}
                        >
                            🏠 Unit Ini Saja
                        </button>
                    </div>
                    <small class="form-hint">
                        {#if postScope === "org"}
                            Pengumuman akan terlihat oleh semua unit dalam
                            organisasi
                        {:else}
                            Pengumuman hanya terlihat oleh warga di unit ini
                        {/if}
                    </small>
                </div>

                <div class="form-group">
                    <label for="title">Judul (Opsional)</label>
                    <input
                        type="text"
                        id="title"
                        bind:value={newPost.title}
                        placeholder="Judul pengumuman"
                    />
                </div>

                <div class="form-group">
                    <label for="content">Isi Pengumuman</label>
                    <textarea
                        id="content"
                        bind:value={newPost.content}
                        placeholder="Tulis isi pengumuman..."
                        rows="4"
                    ></textarea>
                </div>

                {#if newPost.post_type === "photo"}
                    <div class="form-group">
                        <label for="imageUrl">URL Gambar</label>
                        <input
                            type="url"
                            id="imageUrl"
                            bind:value={newPost.image_url}
                            placeholder="https://..."
                        />
                    </div>
                {/if}

                {#if newPost.post_type === "link"}
                    <div class="form-group">
                        <label for="linkUrl">URL Link</label>
                        <input
                            type="url"
                            id="linkUrl"
                            bind:value={newPost.link_url}
                            placeholder="https://..."
                        />
                    </div>
                    <div class="form-group">
                        <label for="linkTitle">Judul Link (Opsional)</label>
                        <input
                            type="text"
                            id="linkTitle"
                            bind:value={newPost.link_title}
                            placeholder="Judul tampilan link"
                        />
                    </div>
                {/if}

                {#if newPost.post_type === "poll"}
                    <div class="form-group">
                        <label>Opsi Polling</label>
                        {#each newPost.poll_options as option, index}
                            <div class="poll-option-input">
                                <input
                                    type="text"
                                    bind:value={newPost.poll_options[index]}
                                    placeholder="Opsi {index + 1}"
                                />
                                {#if newPost.poll_options.length > 2}
                                    <button
                                        class="btn-remove"
                                        on:click={() => removePollOption(index)}
                                    >
                                        ×
                                    </button>
                                {/if}
                            </div>
                        {/each}
                        <button class="btn-add-option" on:click={addPollOption}>
                            + Tambah Opsi
                        </button>
                    </div>
                {/if}

                <div class="form-row">
                    <label class="checkbox-label">
                        <input
                            type="checkbox"
                            bind:checked={newPost.is_pinned}
                        />
                        📌 Sematkan di atas
                    </label>
                    <label class="checkbox-label">
                        <input
                            type="checkbox"
                            bind:checked={newPost.is_important}
                        />
                        ⚠️ Tandai Penting
                    </label>
                </div>
            </div>

            <div class="modal-footer">
                <button
                    class="btn-cancel"
                    on:click={() => (showCreateModal = false)}
                >
                    Batal
                </button>
                <button
                    class="btn-submit"
                    on:click={createPost}
                    disabled={!newPost.content.trim()}
                >
                    Posting
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .pengumuman-container {
        padding: 1.5rem;
        max-width: 1400px;
        margin: 0 auto;
    }

    .page-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.5rem;
        padding-bottom: 1rem;
        border-bottom: 1px solid #e5e7eb;
    }

    .header-content h1 {
        font-size: 1.75rem;
        font-weight: 700;
        color: #111827;
        margin: 0;
    }

    .subtitle {
        color: #6b7280;
        margin-top: 0.25rem;
    }

    .btn-create {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.25rem;
        background: linear-gradient(135deg, #3b82f6, #2563eb);
        color: white;
        border: none;
        border-radius: 0.5rem;
        font-weight: 600;
        cursor: pointer;
        transition:
            transform 0.2s,
            box-shadow 0.2s;
    }

    .btn-create:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
    }

    .btn-create .icon {
        font-size: 1.25rem;
    }

    .error-message {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        background: #fef2f2;
        color: #dc2626;
        border-radius: 0.5rem;
        margin-bottom: 1rem;
    }

    .error-message button {
        background: none;
        border: none;
        font-size: 1.25rem;
        cursor: pointer;
        color: #dc2626;
    }

    .content-wrapper {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
    }

    @media (max-width: 1024px) {
        .content-wrapper {
            grid-template-columns: 1fr;
        }
    }

    .posts-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .loading-state,
    .empty-state {
        text-align: center;
        padding: 3rem;
        background: white;
        border-radius: 0.75rem;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .empty-state .icon {
        font-size: 3rem;
        margin-bottom: 1rem;
    }

    .empty-state h3 {
        margin: 0;
        color: #374151;
    }

    .empty-state p {
        color: #6b7280;
        margin: 0.5rem 0 0;
    }

    .loader {
        width: 40px;
        height: 40px;
        border: 3px solid #e5e7eb;
        border-top-color: #3b82f6;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin: 0 auto 1rem;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .post-card {
        background: white;
        border-radius: 0.75rem;
        padding: 1rem;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        cursor: pointer;
        transition:
            transform 0.2s,
            box-shadow 0.2s;
        border-left: 4px solid transparent;
    }

    .post-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    .post-card.selected {
        border-left-color: #3b82f6;
        background: #f0f7ff;
    }

    .post-card.pinned {
        border-left-color: #f59e0b;
    }

    .post-card.important {
        border-left-color: #ef4444;
    }

    .post-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 0.5rem;
    }

    .post-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        align-items: center;
    }

    .post-type-icon {
        font-size: 1.25rem;
    }

    .badge {
        font-size: 0.75rem;
        padding: 0.125rem 0.5rem;
        border-radius: 9999px;
        font-weight: 500;
    }

    .badge.pinned {
        background: #fef3c7;
        color: #92400e;
    }

    .badge.important {
        background: #fee2e2;
        color: #991b1b;
    }

    .badge.rt {
        background: #dbeafe;
        color: #1e40af;
    }

    .badge.unit {
        background: #d1fae5;
        color: #065f46;
    }

    .btn-delete {
        background: none;
        border: none;
        cursor: pointer;
        opacity: 0.5;
        transition: opacity 0.2s;
    }

    .btn-delete:hover {
        opacity: 1;
    }

    .post-title {
        font-size: 1.1rem;
        font-weight: 600;
        color: #111827;
        margin: 0 0 0.5rem;
    }

    .post-content {
        color: #4b5563;
        margin: 0;
        line-height: 1.5;
    }

    .post-image-preview {
        margin-top: 0.75rem;
        border-radius: 0.5rem;
        overflow: hidden;
        max-height: 150px;
    }

    .post-image-preview img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .post-footer {
        display: flex;
        gap: 1rem;
        margin-top: 0.75rem;
        padding-top: 0.75rem;
        border-top: 1px solid #e5e7eb;
        font-size: 0.85rem;
        color: #6b7280;
    }

    .post-detail {
        position: sticky;
        top: 1rem;
        background: white;
        border-radius: 0.75rem;
        padding: 1.5rem;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        max-height: calc(100vh - 8rem);
        overflow-y: auto;
    }

    .btn-close {
        position: absolute;
        top: 1rem;
        right: 1rem;
        background: #f3f4f6;
        border: none;
        width: 32px;
        height: 32px;
        border-radius: 50%;
        font-size: 1.25rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .detail-header {
        margin-bottom: 1.5rem;
    }

    .detail-header .badges {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }

    .detail-header h2 {
        font-size: 1.5rem;
        color: #111827;
        margin: 0 0 0.5rem;
    }

    .author-info {
        display: flex;
        gap: 1rem;
        color: #6b7280;
        font-size: 0.9rem;
    }

    .detail-content p {
        color: #374151;
        line-height: 1.7;
        white-space: pre-line;
    }

    .post-image {
        margin: 1rem 0;
        border-radius: 0.5rem;
        overflow: hidden;
    }

    .post-image img {
        width: 100%;
        max-height: 400px;
        object-fit: contain;
    }

    .post-link {
        display: inline-block;
        padding: 0.75rem 1rem;
        background: #f3f4f6;
        color: #1d4ed8;
        border-radius: 0.5rem;
        text-decoration: none;
        margin: 1rem 0;
    }

    .poll-options {
        margin: 1rem 0;
    }

    .poll-options h4 {
        margin: 0 0 0.75rem;
    }

    .poll-option {
        display: flex;
        justify-content: space-between;
        width: 100%;
        padding: 0.75rem 1rem;
        background: #f9fafb;
        border: 2px solid #e5e7eb;
        border-radius: 0.5rem;
        cursor: pointer;
        margin-bottom: 0.5rem;
        transition: all 0.2s;
    }

    .poll-option:hover {
        background: #eff6ff;
        border-color: #3b82f6;
    }

    .poll-option.voted {
        background: #dbeafe;
        border-color: #3b82f6;
    }

    .vote-count {
        font-weight: 500;
        color: #3b82f6;
    }

    .comments-section {
        margin-top: 1.5rem;
        padding-top: 1.5rem;
        border-top: 1px solid #e5e7eb;
    }

    .comments-section h4 {
        margin: 0 0 1rem;
    }

    .comment-form {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    .comment-form textarea {
        flex: 1;
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        resize: none;
    }

    .btn-submit {
        padding: 0.5rem 1rem;
        background: #3b82f6;
        color: white;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        font-weight: 500;
    }

    .btn-submit:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .comment {
        padding: 0.75rem;
        background: #f9fafb;
        border-radius: 0.5rem;
        margin-bottom: 0.75rem;
    }

    .comment-header {
        display: flex;
        justify-content: space-between;
        font-size: 0.85rem;
        margin-bottom: 0.5rem;
    }

    .comment-header .author {
        font-weight: 500;
        color: #374151;
    }

    .comment-header .date {
        color: #9ca3af;
    }

    .comment-content {
        color: #4b5563;
        margin: 0;
    }

    .btn-reply {
        background: none;
        border: none;
        color: #3b82f6;
        font-size: 0.85rem;
        cursor: pointer;
        margin-top: 0.5rem;
    }

    .reply-form {
        margin-top: 0.75rem;
        padding: 0.75rem;
        background: white;
        border-radius: 0.375rem;
    }

    .reply-form textarea {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        resize: none;
        margin-bottom: 0.5rem;
    }

    .reply-actions {
        display: flex;
        gap: 0.5rem;
        justify-content: flex-end;
    }

    .btn-cancel {
        padding: 0.375rem 0.75rem;
        background: #f3f4f6;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
    }

    .replies {
        margin-top: 0.75rem;
        padding-left: 1rem;
        border-left: 2px solid #e5e7eb;
    }

    .reply {
        padding: 0.5rem;
        background: white;
        border-radius: 0.375rem;
        margin-bottom: 0.5rem;
    }

    /* Modal */
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 1rem;
    }

    .modal {
        background: white;
        border-radius: 0.75rem;
        width: 100%;
        max-width: 600px;
        max-height: 90vh;
        overflow-y: auto;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem;
        border-bottom: 1px solid #e5e7eb;
    }

    .modal-header h2 {
        margin: 0;
        font-size: 1.25rem;
    }

    .modal-body {
        padding: 1.25rem;
    }

    .form-group {
        margin-bottom: 1rem;
    }

    .form-group label {
        display: block;
        margin-bottom: 0.375rem;
        font-weight: 500;
        color: #374151;
    }

    .form-group input,
    .form-group textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        font-size: 1rem;
    }

    .form-group textarea {
        resize: vertical;
    }

    .form-select {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        font-size: 1rem;
        background: white;
        cursor: pointer;
    }

    .form-select:focus {
        outline: none;
        border-color: #10b981;
        box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
    }

    .form-hint {
        display: block;
        margin-top: 0.5rem;
        font-size: 0.875rem;
        color: #6b7280;
    }

    .post-type-selector {
        display: flex;
        gap: 0.5rem;
    }

    .post-type-selector button {
        flex: 1;
        padding: 0.75rem;
        background: #f3f4f6;
        border: 2px solid transparent;
        border-radius: 0.5rem;
        cursor: pointer;
        font-weight: 500;
        transition: all 0.2s;
    }

    .post-type-selector button:hover {
        background: #e5e7eb;
    }

    .post-type-selector button.active {
        background: #dbeafe;
        border-color: #3b82f6;
        color: #1d4ed8;
    }

    .poll-option-input {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }

    .poll-option-input input {
        flex: 1;
    }

    .btn-remove {
        width: 36px;
        height: 36px;
        background: #fee2e2;
        color: #dc2626;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
        font-size: 1.25rem;
    }

    .btn-add-option {
        background: #e0f2fe;
        color: #0369a1;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 0.375rem;
        cursor: pointer;
        font-weight: 500;
    }

    .form-row {
        display: flex;
        gap: 1.5rem;
    }

    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
    }

    .checkbox-label input {
        width: 18px;
        height: 18px;
    }

    .modal-footer {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        padding: 1.25rem;
        border-top: 1px solid #e5e7eb;
    }

    .modal-footer .btn-cancel {
        padding: 0.75rem 1.5rem;
    }

    .modal-footer .btn-submit {
        padding: 0.75rem 1.5rem;
    }

    .pagination {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 1rem;
        padding: 1rem;
    }

    .pagination button {
        padding: 0.5rem 1rem;
        background: #f3f4f6;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
    }

    .pagination button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
</style>
