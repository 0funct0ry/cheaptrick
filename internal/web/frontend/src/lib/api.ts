import type { RequestItem, RequestDetail, Template, Fixture } from './types';

const API_BASE = '/api';

export const api = {
    getRequests: async (): Promise<RequestItem[]> => {
        const res = await fetch(`${API_BASE}/requests`);
        const data = await res.json();
        return data.requests || [];
    },

    getRequest: async (id: string): Promise<RequestDetail> => {
        const res = await fetch(`${API_BASE}/requests/${id}`);
        if (!res.ok) throw new Error('Failed to fetch request detail');
        return res.json();
    },

    respondToRequest: async (id: string, response: unknown): Promise<{ ok: boolean }> => {
        const res = await fetch(`${API_BASE}/requests/${id}/respond`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ response })
        });
        if (!res.ok) throw new Error('Failed to respond');
        return res.json();
    },

    deleteRequest: async (id: string): Promise<void> => {
        const res = await fetch(`${API_BASE}/requests/${id}`, {
            method: 'DELETE'
        });
        if (!res.ok) throw new Error('Failed to delete request');
    },

    clearRequests: async (): Promise<void> => {
        const res = await fetch(`${API_BASE}/requests`, {
            method: 'DELETE'
        });
        if (!res.ok) throw new Error('Failed to clear requests');
    },

    saveFixture: async (id: string, response: unknown): Promise<{ hash: string, path: string }> => {
        const res = await fetch(`${API_BASE}/requests/${id}/fixture`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ response })
        });
        if (!res.ok) throw new Error('Failed to save fixture');
        return res.json();
    },

    getTemplates: async (): Promise<Template[]> => {
        const res = await fetch(`${API_BASE}/templates`);
        const data = await res.json();
        return data.templates || [];
    },

    getFixtures: async (): Promise<Fixture[]> => {
        const res = await fetch(`${API_BASE}/fixtures`);
        const data = await res.json();
        return data.fixtures || [];
    },

    getFixture: async (hash: string): Promise<{ content: unknown }> => {
        const res = await fetch(`${API_BASE}/fixtures/${hash}`);
        if (!res.ok) throw new Error('Failed to fetch fixture');
        return res.json();
    },

    createFixture: async (hash: string | null, response: unknown, request?: unknown): Promise<{ ok: boolean, hash: string, path: string }> => {
        const body: Record<string, unknown> = { response };
        if (hash) body.hash = hash;
        if (request) body.request = request;

        const res = await fetch(`${API_BASE}/fixtures`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        if (!res.ok) throw new Error('Failed to create/update fixture');
        return res.json();
    },

    deleteFixture: async (hash: string): Promise<void> => {
        const res = await fetch(`${API_BASE}/fixtures/${hash}`, {
            method: 'DELETE'
        });
        if (!res.ok) throw new Error('Failed to delete fixture');
    }
};
