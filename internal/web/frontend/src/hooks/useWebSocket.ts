import { useState, useEffect, useCallback, useRef } from 'react';
import type { WsStatus } from '../lib/types';

export function useWebSocket(onEvent: (event: any) => void) {
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
            setTimeout(() => {
                setStatus('reconnecting');
                connect();
            }, 2000);
        };

        ws.onerror = () => {
            ws.close();
        };
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

    return { status };
}
