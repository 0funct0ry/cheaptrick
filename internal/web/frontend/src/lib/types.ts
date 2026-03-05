export interface RequestItem {
    id: string;
    model: string;
    timestamp: string;
    status: 'pending' | 'responded' | 'auto';
    via?: string;
    fixture_hash?: string;
    preview: string;
    body: any;
    tools: string[];
}

export interface RequestDetail extends RequestItem {
    system_instruction: string;
    contents: any[];
    generation_config: any;
    response?: any;
}

export interface Template {
    id: string;
    label: string;
    shortcut: string;
    body: any;
}

export interface Fixture {
    hash: string;
    size: number;
}

export type WsStatus = 'connecting' | 'connected' | 'reconnecting' | 'disconnected';
