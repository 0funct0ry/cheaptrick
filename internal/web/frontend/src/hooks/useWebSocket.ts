import { useState, useEffect, useCallback, useRef } from 'react';
import type { WsStatus } from '../lib/types';

export function useWebSocket(onEvent: (event: Record<string, unknown>) => void) {
    const [status, setStatus] = useState<WsStatus>('connecting');
    const wsRef = useRef<WebSocket | null>(null);

    const connect = useCallback(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;

        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => setStatus('connected');

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                onEvent(data);
            } catch (e) {
                console.error('Failed to parse WS msg', e);
            }
        };

        ws.onclose = () => {
            setStatus('disconnected');
            // Remove the auto-reconnect from here to avoid the scope issue,
            // or pass it in as a param/handle it in a separate effect.
        };

        ws.onerror = () => {
            ws.close();
        };

        return ws;
    }, [onEvent]);

    useEffect(() => {
        connect();
        return () => {
            if (wsRef.current) {
                wsRef.current.onclose = null;
                wsRef.current.close();
            }
        };
    }, [connect]);

    useEffect(() => {
        if (status === 'disconnected') {
            const timer = setTimeout(() => {
                setStatus('reconnecting');
                connect();
            }, 2000);
            return () => clearTimeout(timer);
        }
    }, [status, connect]);

    return { status };
}
