import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';
import type { RequestItem } from '../lib/types';

export function useRequests() {
    const [requests, setRequests] = useState<RequestItem[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchRequests = useCallback(async () => {
        try {
            setLoading(true);
            const data = await api.getRequests();
            setRequests(data);
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchRequests();
    }, [fetchRequests]);

    const handleWsEvent = useCallback((event: Record<string, unknown>) => {
        if (event.type === 'new_request') {
            setRequests(prev => [event.request as RequestItem, ...prev]);
        } else if (event.type === 'request_responded') {
            setRequests(prev => prev.map(req =>
                req.id === event.id ? { ...req, status: 'responded', via: event.via as string } : req
            ));
        } else if (event.type === 'fixture_saved') {
            setRequests(prev => prev.map(req =>
                req.id === event.request_id ? { ...req, fixture_hash: event.hash as string } : req
            ));
        } else if (event.type === 'request_deleted') {
            setRequests(prev => prev.filter(req => req.id !== event.id));
        } else if (event.type === 'requests_cleared') {
            setRequests(prev => prev.filter(req => req.status === 'pending'));
        }
    }, []);

    return { requests, loading, handleWsEvent, fetchRequests };
}
