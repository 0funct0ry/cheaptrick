export interface RequestItem {
    id: string;
    model: string;
    timestamp: string;
    status: 'pending' | 'responded' | 'auto';
    via?: string;
    fixture_hash?: string;
    preview: string;
    body: unknown;
    tools: string[];
}

export interface RequestDetail extends RequestItem {
    system_instruction: string;
    contents: unknown[];
    generation_config: unknown;
    response?: unknown;
}

export interface Template {
    id: string;
    label: string;
    shortcut: string;
    body: unknown;
}

export interface Fixture {
    hash: string;
    size: number;
}

export type WsStatus = 'connecting' | 'connected' | 'reconnecting' | 'disconnected';
