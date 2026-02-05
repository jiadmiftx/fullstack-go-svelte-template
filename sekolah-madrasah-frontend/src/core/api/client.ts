import type {
    ApiResponse,
    PaginatedResponse,
    Role,
    Permission,
    User,
} from '../types';

import type { Organization } from '../../business/types';

// API Configuration
export const API_BASE_URL = 'http://localhost:8080/api/v1';

// API Client wrapper
class ApiClient {
    private baseUrl: string;
    private token: string | null = null;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    setToken(token: string | null) {
        this.token = token;
    }

    private getHeaders(): HeadersInit {
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        return headers;
    }

    async get<T>(endpoint: string): Promise<T> {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            method: 'GET',
            headers: this.getHeaders(),
        });
        if (!response.ok) {
            throw await this.handleError(response);
        }
        return response.json();
    }

    async post<T>(endpoint: string, data: unknown): Promise<T> {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw await this.handleError(response);
        }
        return response.json();
    }

    async put<T>(endpoint: string, data: unknown): Promise<T> {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            method: 'PUT',
            headers: this.getHeaders(),
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw await this.handleError(response);
        }
        return response.json();
    }

    async delete<T>(endpoint: string): Promise<T> {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            method: 'DELETE',
            headers: this.getHeaders(),
        });
        if (!response.ok) {
            throw await this.handleError(response);
        }
        return response.json();
    }

    private async handleError(response: Response): Promise<Error> {
        try {
            const data = await response.json();
            return new Error(data.error || data.message || 'An error occurred');
        } catch {
            return new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
    }

    // Organizations API
    async getOrganizations(params?: { page?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Organization[]>>(`/organizations?${query}`);
    }

    async getOrganization(id: string) {
        return this.get<ApiResponse<Organization>>(`/organizations/${id}`);
    }

    async createOrganization(data: { name: string; code: string; type: string; description?: string; address?: string; logo?: string }) {
        return this.post<ApiResponse<Organization>>('/organizations', data);
    }

    async updateOrganization(id: string, data: { name?: string; description?: string; address?: string; logo?: string; settings?: string }) {
        return this.put<ApiResponse<Organization>>(`/organizations/${id}`, data);
    }

    async deleteOrganization(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/organizations/${id}`);
    }

    // Organization Members API
    async addMemberToOrganization(orgId: string, data: { user_id: string; role_id: string }) {
        return this.post<ApiResponse<any>>(`/organizations/${orgId}/members`, data);
    }

    async getOrganizationMembers(orgId: string, params?: { page?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/organizations/${orgId}/members?${query}`);
    }

    async updateMemberRole(orgId: string, userId: string, data: { role_id: string }) {
        return this.put<ApiResponse<any>>(`/organizations/${orgId}/members/${userId}`, data);
    }

    async removeMemberFromOrganization(orgId: string, userId: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/organizations/${orgId}/members/${userId}`);
    }

    // Roles API
    async getRoles(params?: { page?: number; limit?: number; organization_id?: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Role[]>>(`/roles?${query}`);
    }

    async getRole(id: string) {
        return this.get<ApiResponse<Role>>(`/roles/${id}`);
    }

    async createRole(data: {
        organization_id: string;
        name: string;
        display_name: string;
        type: string;
        level: number;
        description: string;
        is_default?: boolean;
        permission_ids?: string[]
    }) {
        return this.post<ApiResponse<Role>>('/roles', data);
    }

    async updateRole(id: string, data: {
        name?: string;
        display_name?: string;
        type?: string;
        level?: number;
        description?: string;
        permission_ids?: string[]
    }) {
        return this.put<ApiResponse<Role>>(`/roles/${id}`, data);
    }

    async deleteRole(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/roles/${id}`);
    }

    // Permissions API
    async getPermissions(params?: { page?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Permission[]>>(`/permissions?${query}`);
    }

    async getPermission(id: string) {
        return this.get<ApiResponse<Permission>>(`/permissions/${id}`);
    }

    async createPermission(data: { resource: string; action: string; description?: string }) {
        return this.post<ApiResponse<Permission>>('/permissions', data);
    }

    async deletePermission(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/permissions/${id}`);
    }

    // Auth API
    async register(data: { email: string; password: string; full_name: string; phone?: string }) {
        return this.post<ApiResponse<{ user: User; token: string }>>('/auth/register', data);
    }

    // Users API
    async getUsers(params?: { page?: number; limit?: number; platform_only?: boolean }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<User[]>>(`/users?${query}`);
    }

    async getUser(id: string) {
        return this.get<ApiResponse<User>>(`/users/${id}`);
    }

    async createUser(data: { email: string; password: string; full_name: string; phone?: string }) {
        return this.post<ApiResponse<User>>('/users', data);
    }

    async updateUser(id: string, data: { email?: string; password?: string; full_name?: string; phone?: string }) {
        return this.put<ApiResponse<User>>(`/users/${id}`, data);
    }

    async deleteUser(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/users/${id}`);
    }

    async getCurrentUser() {
        return this.get<ApiResponse<{
            id: string;
            email: string;
            full_name: string;
            phone: string;
            avatar: string;
            is_super_admin: boolean;
            is_active: boolean;
        }>>('/users/me');
    }

    async getMyMemberships() {
        return this.get<ApiResponse<{
            user_id: string;
            is_super_admin: boolean;
            organization_memberships: Array<{
                org_id: string;
                org_name: string;
                role_id: string;
                role_name: string;
            }>;
            unit_memberships: Array<{
                unit_member_id: string;
                unit_id: string;
                perumahan_name: string;
                org_id: string;
                org_name: string;
                role: 'pengurus' | 'warga' | 'admin' | 'staff' | 'parent';
                is_active: boolean;
            }>;
        }>>('/users/me/memberships');
    }

    // Perumahans API
    async getPerumahans(params?: { page?: number; limit?: number; organization_id?: string; type?: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/units?${query}`);
    }

    async getPerumahan(id: string) {
        return this.get<ApiResponse<any>>(`/units/${id}`);
    }

    async createPerumahan(data: {
        organization_id: string;
        name: string;
        code: string;
        type?: string;
        address?: string;
        phone?: string;
        email?: string;
        logo?: string;
    }) {
        return this.post<ApiResponse<any>>('/units', data);
    }

    async updatePerumahan(id: string, data: {
        name?: string;
        address?: string;
        phone?: string;
        email?: string;
        logo?: string;
        settings?: string;
    }) {
        return this.put<ApiResponse<any>>(`/units/${id}`, data);
    }

    async deletePerumahan(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/units/${id}`);
    }

    // Pengurus Profiles API
    async getPenguruss(params?: { page?: number; limit?: number; unit_id?: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/pengurus-profiles?${query}`);
    }

    async getPengurus(id: string) {
        return this.get<ApiResponse<any>>(`/pengurus-profiles/${id}`);
    }

    async createPengurus(data: {
        unit_member_id: string;
        nip?: string;
        specialization?: string;
        gender?: string;
        address?: string;
    }) {
        return this.post<ApiResponse<any>>('/pengurus-profiles', data);
    }

    async updatePengurus(id: string, data: {
        nip?: string;
        specialization?: string;
        gender?: string;
        address?: string;
        status?: string;
    }) {
        return this.put<ApiResponse<any>>(`/pengurus-profiles/${id}`, data);
    }

    async deletePengurus(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/pengurus-profiles/${id}`);
    }

    // Pengurus-Kegiatan Assignments API
    async getPengurusKegiatans(pengurusId: string, params?: { page?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/pengurus-profiles/${pengurusId}/kegiatans?${query}`);
    }

    async assignKegiatanToPengurus(pengurusId: string, data: {
        kegiatan_id: string;
        academic_year?: string;
        is_primary?: boolean;
    }) {
        return this.post<ApiResponse<any>>(`/pengurus-profiles/${pengurusId}/kegiatans`, data);
    }

    async removeKegiatanFromPengurus(pengurusId: string, kegiatanId: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/pengurus-profiles/${pengurusId}/kegiatans/${kegiatanId}`);
    }

    async getPengurussByKegiatan(kegiatanId: string, params?: { page?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/kegiatans/${kegiatanId}/penguruss?${query}`);
    }

    // Warga Profiles API
    async getWargas(params?: { page?: number; limit?: number; unit_id?: string; rt_id?: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/warga-profiles?${query}`);
    }

    async getWarga(id: string) {
        return this.get<ApiResponse<any>>(`/warga-profiles/${id}`);
    }

    async createWarga(data: {
        unit_member_id: string;
        rt_id?: string;
        nis: string;
        nisn?: string;
        gender?: string;
        birth_place?: string;
        address?: string;
        parent_name?: string;
        parent_phone?: string;
    }) {
        return this.post<ApiResponse<any>>('/warga-profiles', data);
    }

    async updateWarga(id: string, data: {
        // Kategori 1: Identitas
        nik?: string;
        gender?: string;
        agama?: string;
        pekerjaan?: string;
        no_whatsapp?: string;
        // Kategori 2: Domisili
        blok_rumah?: string;
        nomor_rumah?: string;
        rt_rw?: string;
        status_kepemilikan?: string;
        status_hunian?: string;
        // Kategori 3: Keluarga
        jumlah_anggota_keluarga?: number;
        nama_kontak_darurat?: string;
        no_kontak_darurat?: string;
        no_plat_mobil?: string;
        no_plat_motor?: string;
        memiliki_art?: boolean;
        // Kategori 4: Keuangan
        status_iuran?: string;
        metode_pembayaran?: string;
        keterangan_khusus?: string;
        // System
        status?: string;
    }) {
        return this.put<ApiResponse<any>>(`/warga-profiles/${id}`, data);
    }

    async deleteWarga(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/warga-profiles/${id}`);
    }

    async getWargaRTHistory(wargaId: string) {
        return this.get<ApiResponse<any[]>>(`/warga-profiles/${wargaId}/rt-history`);
    }

    // RTs API
    async getRTs(params?: { page?: number; limit?: number; unit_id?: string; iuran?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/rts?${query}`);
    }

    async getRT(id: string) {
        return this.get<ApiResponse<any>>(`/rts/${id}`);
    }

    async createRT(data: {
        unit_id: string;
        name: string;
        iuran: number;
        homeroom_id?: string;
        academic_year: string;
        type?: string;
        capacity?: number;
    }) {
        return this.post<ApiResponse<any>>('/rts', data);
    }

    async updateRT(id: string, data: {
        name?: string;
        iuran?: number;
        homeroom_id?: string;
        academic_year?: string;
        type?: string;
        capacity?: number;
        is_active?: boolean;
    }) {
        return this.put<ApiResponse<any>>(`/rts/${id}`, data);
    }

    async deleteRT(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/rts/${id}`);
    }

    // RT-Warga API
    async getWargasByRT(rtId: string) {
        return this.get<ApiResponse<any[]>>(`/rts/${rtId}/wargas`);
    }

    async addWargaToRT(rtId: string, data: { warga_id: string; rt_type: string }) {
        return this.post<ApiResponse<any>>(`/rts/${rtId}/wargas`, data);
    }

    async removeWargaFromRT(rtId: string, wargaId: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/rts/${rtId}/wargas/${wargaId}`);
    }

    async getRTsByWarga(wargaId: string) {
        return this.get<ApiResponse<any[]>>(`/warga-profiles/${wargaId}/rts`);
    }
    // Rooms API
    async getRooms(params?: { page?: number; limit?: number; unit_id?: string; type?: string; building?: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/rooms?${query}`);
    }

    async getRoom(id: string) {
        return this.get<ApiResponse<any>>(`/rooms/${id}`);
    }

    async createRoom(data: {
        unit_id: string;
        code: string;
        name: string;
        type?: string;
        building?: string;
        floor?: number;
        capacity?: number;
        facilities?: string;
    }) {
        return this.post<ApiResponse<any>>('/rooms', data);
    }

    async updateRoom(id: string, data: {
        code?: string;
        name?: string;
        type?: string;
        building?: string;
        floor?: number;
        capacity?: number;
        facilities?: string;
    }) {
        return this.put<ApiResponse<any>>(`/rooms/${id}`, data);
    }

    async deleteRoom(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/rooms/${id}`);
    }

    // Perumahan Members API
    async getPerumahanMembers(unitId: string, params?: { page?: number; limit?: number; role?: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/units/${unitId}/members?${query}`);
    }

    async getPerumahanMember(unitId: string, memberId: string) {
        return this.get<ApiResponse<any>>(`/units/${unitId}/members/${memberId}`);
    }

    async addPerumahanMember(unitId: string, data: {
        user_id: string;
        role: 'admin' | 'pengurus' | 'warga' | 'parent' | 'staff';
    }) {
        return this.post<ApiResponse<any>>(`/units/${unitId}/members`, data);
    }

    async updatePerumahanMember(unitId: string, memberId: string, data: {
        role?: 'admin' | 'pengurus' | 'warga' | 'parent' | 'staff';
        is_active?: boolean;
    }) {
        return this.put<ApiResponse<any>>(`/units/${unitId}/members/${memberId}`, data);
    }

    async removePerumahanMember(unitId: string, memberId: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/units/${unitId}/members/${memberId}`);
    }

    // Quick Registration API
    async registerPengurus(unitId: string, data: {
        name: string;
        email: string;
        nip: string;
        phone?: string;
        address?: string;
    }) {
        return this.post<ApiResponse<{
            user_id: string;
            pengurus_id: string;
            unit_member_id: string;
            generated_password: string;
        }>>(`/units/${unitId}/penguruss/register`, data);
    }

    async registerWarga(unitId: string, data: {
        name: string;
        email: string;
        phone?: string;
        // Kategori 1: Identitas
        nik?: string;
        gender?: string;
        agama?: string;
        pekerjaan?: string;
        // Kategori 2: Domisili
        blok_rumah?: string;
        nomor_rumah?: string;
        status_kepemilikan?: string;
        status_hunian?: string;
    }) {
        return this.post<ApiResponse<{
            user_id: string;
            warga_id: string;
            unit_member_id: string;
            generated_password: string;
        }>>(`/units/${unitId}/wargas/register`, data);
    }

    // Schedule API
    async getSchedules(params?: { unit_id?: string; rt_id?: string; pengurus_id?: string; day_of_week?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Array<{
            id: string;
            unit_id: string;
            rt_id: string;
            kegiatan_id: string;
            pengurus_id: string;
            day_of_week: number;
            start_time: string;
            end_time: string;
            room?: string;
            is_active: boolean;
            rt?: { id: string; name: string };
            kegiatan?: { id: string; name: string };
            pengurus?: { id: string; name: string };
        }>>>(`/schedules?${query}`);
    }

    async getSchedule(id: string) {
        return this.get<ApiResponse<{
            id: string;
            unit_id: string;
            rt_id: string;
            kegiatan_id: string;
            pengurus_id: string;
            day_of_week: number;
            start_time: string;
            end_time: string;
            room?: string;
            is_active: boolean;
        }>>(`/schedules/${id}`);
    }

    async createSchedule(data: {
        unit_id: string;
        rt_id: string;
        kegiatan_id: string;
        pengurus_id: string;
        day_of_week: number;
        start_time: string;
        end_time: string;
        room?: string;
        start_date?: string;
        end_date?: string;
        recurrence_type?: string;
    }) {
        return this.post<ApiResponse<{ id: string }>>('/schedules', data);
    }

    async updateSchedule(id: string, data: {
        rt_id?: string;
        kegiatan_id?: string;
        pengurus_id?: string;
        day_of_week?: number;
        start_time?: string;
        end_time?: string;
        room?: string;
        is_active?: boolean;
        start_date?: string;
        end_date?: string;
        recurrence_type?: string;
    }) {
        return this.put<ApiResponse<{ id: string }>>(`/schedules/${id}`, data);
    }

    async deleteSchedule(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/schedules/${id}`);
    }

    async checkScheduleConflicts(data: {
        unit_id: string;
        rt_id: string;
        pengurus_id: string;
        room: string;
        day_of_week: number;
        start_time: string;
        end_time: string;
        exclude_id?: string;
    }) {
        return this.post<ApiResponse<{
            has_conflict: boolean;
            conflicts: Array<{
                type: string;
                conflicts_with: any;
                message: string;
            }>;
        }>>('/schedules/check-conflicts', data);
    }

    async copyScheduleFromRT(data: {
        source_rt_id: string;
        target_rt_id: string;
        unit_id: string;
    }) {
        return this.post<ApiResponse<any[]>>('/schedules/copy-from-rt', data);
    }

    // Schedule Template API
    async getScheduleTemplates(params: { unit_id: string }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Array<{
            id: string;
            unit_id: string;
            name: string;
            description?: string;
            academic_year?: string;
            is_active: boolean;
            items?: any[];
        }>>>(`/schedules/templates?${query}`);
    }

    async saveScheduleTemplate(data: {
        unit_id: string;
        rt_id: string;
        name: string;
        description?: string;
        academic_year?: string;
    }) {
        return this.post<ApiResponse<any>>('/schedules/templates', data);
    }

    async loadScheduleTemplate(data: {
        template_id: string;
        rt_id: string;
        unit_id: string;
    }) {
        return this.post<ApiResponse<any[]>>('/schedules/templates/load', data);
    }

    async deleteScheduleTemplate(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/schedules/templates/${id}`);
    }

    // Kegiatans API
    async getKegiatans(params?: { unit_id?: string; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Array<{
            id: string;
            unit_id: string;
            code: string;
            name: string;
            description?: string;
            credits: number;
            is_active: boolean;
        }>>>(`/kegiatans?${query}`);
    }

    async getKegiatan(id: string) {
        return this.get<ApiResponse<{
            id: string;
            unit_id: string;
            code: string;
            name: string;
            description?: string;
            credits: number;
            is_active: boolean;
        }>>(`/kegiatans/${id}`);
    }

    async createKegiatan(data: {
        unit_id: string;
        code: string;
        name: string;
        description?: string;
        credits?: number;
    }) {
        return this.post<ApiResponse<{ id: string }>>('/kegiatans', data);
    }

    async updateKegiatan(id: string, data: {
        code?: string;
        name?: string;
        description?: string;
        credits?: number;
        is_active?: boolean;
    }) {
        return this.put<ApiResponse<{ id: string }>>(`/kegiatans/${id}`, data);
    }

    async deleteKegiatan(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/kegiatans/${id}`);
    }

    // Iurans API
    async getIurans(params?: {
        unit_id?: string;
        warga_id?: string;
        rt_id?: string;
        kegiatan_id?: string;
        academic_year?: string;
        semester?: number;
        type?: string;
        limit?: number
    }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Array<{
            id: string;
            unit_id: string;
            warga_id: string;
            kegiatan_id: string;
            rt_id: string;
            pengurus_id: string;
            academic_year: string;
            semester: number;
            type: string;
            score: number;
            max_score: number;
            notes?: string;
            warga?: { id: string; name: string; nis: string };
            kegiatan?: { id: string; name: string };
            rt?: { id: string; name: string };
        }>>>(`/iurans?${query}`);
    }

    async createIuran(data: {
        unit_id: string;
        warga_id: string;
        kegiatan_id: string;
        rt_id: string;
        pengurus_id: string;
        academic_year: string;
        semester: number;
        type: string;
        score: number;
        max_score?: number;
        notes?: string;
    }) {
        return this.post<ApiResponse<{ id: string }>>('/iurans', data);
    }

    async updateIuran(id: string, data: {
        score?: number;
        max_score?: number;
        notes?: string;
    }) {
        return this.put<ApiResponse<{ id: string }>>(`/iurans/${id}`, data);
    }

    async deleteIuran(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/iurans/${id}`);
    }

    // Attendance API
    async getAttendances(params?: {
        unit_id?: string;
        warga_id?: string;
        rt_id?: string;
        date?: string;
        status?: string;
        limit?: number
    }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Array<{
            id: string;
            unit_id: string;
            warga_id: string;
            rt_id: string;
            date: string;
            status: string;
            notes?: string;
            warga?: { id: string; name: string; nis: string };
            rt?: { id: string; name: string };
        }>>>(`/attendances?${query}`);
    }

    async createAttendance(data: {
        unit_id: string;
        warga_id: string;
        rt_id: string;
        date: string;
        status: string;
        notes?: string;
    }) {
        return this.post<ApiResponse<{ id: string }>>('/attendances', data);
    }

    async updateAttendance(id: string, data: {
        status?: string;
        notes?: string;
    }) {
        return this.put<ApiResponse<{ id: string }>>(`/attendances/${id}`, data);
    }

    async deleteAttendance(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/attendances/${id}`);
    }

    // Perumahan Settings API
    async getPerumahanSettings(unitId: string) {
        return this.get<ApiResponse<{
            id: string;
            unit_id: string;
            period_duration: number;
            start_time: string;
            total_periods: number;
            break_after_period: number;
            break_duration: number;
            academic_year: string;
            current_semester: number;
            semester_1_start: string | null;
            semester_1_end: string | null;
            semester_2_start: string | null;
            semester_2_end: string | null;
        }>>(`/units/${unitId}/settings`);
    }

    async updatePerumahanSettings(unitId: string, data: {
        period_duration?: number;
        start_time?: string;
        total_periods?: number;
        break_after_period?: number;
        break_duration?: number;
        academic_year?: string;
        current_semester?: number;
        semester_1_start?: string;
        semester_1_end?: string;
        semester_2_start?: string;
        semester_2_end?: string;
    }) {
        return this.put<ApiResponse<{ id: string }>>(`/units/${unitId}/settings`, data);
    }

    // Period Definitions API
    async getPeriodDefinitions(unitId: string) {
        return this.get<ApiResponse<Array<{
            id: string;
            period: number;
            start_time: string;
            end_time: string;
            is_break: boolean;
        }>>>(`/units/${unitId}/periods`);
    }

    async generatePeriodDefinitions(unitId: string) {
        return this.post<ApiResponse<{ message: string }>>(`/units/${unitId}/periods/generate`, {});
    }

    async updatePeriodDefinition(unitId: string, periodId: string, data: {
        start_time?: string;
        end_time?: string;
    }) {
        return this.put<ApiResponse<{ id: string }>>(`/units/${unitId}/periods/${periodId}`, data);
    }

    // Events API
    async getEvents(params?: { unit_id?: string; rt_id?: string; event_type?: string; start_date?: string; end_date?: string; page?: number; limit?: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/events?${query}`);
    }

    async getEvent(id: string) {
        return this.get<ApiResponse<any>>(`/events/${id}`);
    }

    async getCalendarEvents(params: { unit_id: string; year: number; month: number }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<any[]>>(`/events/calendar?${query}`);
    }

    async createEvent(data: {
        unit_id: string;
        rt_id?: string;
        title: string;
        description?: string;
        event_type: string;
        start_date: string;
        end_date?: string;
        start_time?: string;
        end_time?: string;
        recurrence_type: string;
        day_of_week?: number;
        is_all_day?: boolean;
        location?: string;
        color?: string;
    }) {
        return this.post<ApiResponse<any>>('/events', data);
    }

    async updateEvent(id: string, data: {
        rt_id?: string;
        title?: string;
        description?: string;
        event_type?: string;
        start_date?: string;
        end_date?: string;
        start_time?: string;
        end_time?: string;
        recurrence_type?: string;
        day_of_week?: number;
        is_all_day?: boolean;
        location?: string;
        color?: string;
        is_active?: boolean;
    }) {
        return this.put<ApiResponse<any>>(`/events/${id}`, data);
    }

    async deleteEvent(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/events/${id}`);
    }

    // Calendar Entries API (Unified schedules + events)
    async getCalendarEntries(params: {
        unit_id: string;
        entry_type?: 'schedule' | 'event';
        rt_id?: string;
        pengurus_id?: string;
        kegiatan_id?: string;
        year?: number;
        month?: number;
        page?: number;
        limit?: number;
    }) {
        const query = new URLSearchParams(params as any).toString();
        return this.get<ApiResponse<Array<{
            id: string;
            unit_id: string;
            entry_type: 'schedule' | 'event';
            title: string;
            description?: string;
            day_of_week?: number;
            date?: string;
            start_date?: string;
            end_date?: string;
            start_time: string;
            end_time: string;
            is_all_day: boolean;
            recurrence_type: string;
            rt_id?: string;
            kegiatan_id?: string;
            pengurus_id?: string;
            room?: string;
            rt_name?: string;
            kegiatan_name?: string;
            pengurus_name?: string;
            event_type?: string;
            location?: string;
            color?: string;
            is_active: boolean;
            created_at: string;
            updated_at: string;
        }>>>(`/calendar-entries?${query}`);
    }

    async getCalendarEntry(id: string) {
        return this.get<ApiResponse<any>>(`/calendar-entries/${id}`);
    }

    async createCalendarEntry(data: {
        entry_type: 'schedule' | 'event';
        title?: string;
        description?: string;
        day_of_week?: number;
        start_date?: string;
        end_date?: string;
        start_time: string;
        end_time: string;
        is_all_day?: boolean;
        recurrence_type?: string;
        rt_id?: string;
        kegiatan_id?: string;
        pengurus_id?: string;
        room?: string;
        event_type?: string;
        location?: string;
        color?: string;
    }, unitId: string) {
        return this.post<ApiResponse<any>>(`/calendar-entries?unit_id=${unitId}`, data);
    }

    async updateCalendarEntry(id: string, data: {
        title?: string;
        description?: string;
        day_of_week?: number;
        start_date?: string;
        end_date?: string;
        start_time?: string;
        end_time?: string;
        is_all_day?: boolean;
        recurrence_type?: string;
        rt_id?: string;
        kegiatan_id?: string;
        pengurus_id?: string;
        room?: string;
        event_type?: string;
        location?: string;
        color?: string;
    }) {
        return this.put<ApiResponse<any>>(`/calendar-entries/${id}`, data);
    }

    async deleteCalendarEntry(id: string) {
        return this.delete<ApiResponse<{ message: string }>>(`/calendar-entries/${id}`);
    }
}

export const api = new ApiClient(API_BASE_URL);
