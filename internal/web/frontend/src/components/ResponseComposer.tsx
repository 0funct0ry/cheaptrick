import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import type { Template } from '../lib/types';
import { useKeyboardShortcuts } from '../hooks/useKeyboardShortcuts';
import { Save, Send, AlertCircle, FileJson, Check } from 'lucide-react';
import clsx from 'clsx';

export function ResponseComposer({
    reqId,
    isAnswered,
    onSent,
    onSaved
}: {
    reqId: string;
    isAnswered: boolean;
    onSent: (id: string) => void;
    onSaved: (hash: string) => void;
}) {
    const [value, setValue] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [templates, setTemplates] = useState<Template[]>([]);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        api.getTemplates().then(setTemplates).catch(console.error);
    }, []);

    useEffect(() => {
        if (reqId && !value && templates.length > 0) {
            const textTpl = templates.find(t => t.id === 'text');
            if (textTpl) {
                setValue(JSON.stringify(textTpl.body, null, 2));
            }
        }
    }, [reqId, templates]);

    useEffect(() => {
        if (!value.trim()) {
            setError(null);
            return;
        }
        try {
            JSON.parse(value);
            setError(null);
        } catch (e: any) {
            setError(e.message);
        }
    }, [value]);

    const handleSend = async () => {
        if (error || !value.trim() || submitting || isAnswered) return;
        try {
            setSubmitting(true);
            const parsed = JSON.parse(value);
            await api.respondToRequest(reqId, parsed);
            onSent(reqId);
            setValue('');
        } catch (e: any) {
            console.error(e);
            alert(e.message);
        } finally {
            setSubmitting(false);
        }
    };

    const handleSave = async () => {
        if (error || !value.trim() || isAnswered) return;
        try {
            const parsed = JSON.parse(value);
            const res = await api.saveFixture(reqId, parsed);
            onSaved(res.hash);
        } catch (e: any) {
            console.error(e);
            alert(e.message);
        }
    };

    useKeyboardShortcuts({
        'ctrl+s': handleSend,
        'ctrl+f': handleSave,
        'f1': () => applyTemplate('text'),
        'f2': () => applyTemplate('function_call'),
        'f3': () => applyTemplate('error_429'),
        'f4': () => applyTemplate('error_500'),
    }, !isAnswered);

    const applyTemplate = (id: string) => {
        const tpl = templates.find(t => t.id === id);
        if (tpl) {
            setValue(JSON.stringify(tpl.body, null, 2));
        }
    };

    if (isAnswered) {
        return (
            <div className="flex flex-col h-full w-full items-center justify-center p-8 text-center text-zinc-500 bg-zinc-50/50 dark:bg-zinc-900/50">
                <Check className="mx-auto h-12 w-12 text-emerald-500 mb-4 opacity-50" />
                <p className="font-medium text-lg text-zinc-700 dark:text-zinc-300">Request Responded</p>
                <p className="text-sm mt-2 opacity-70">This request has already been processed.</p>
            </div>
        );
    }

    return (
        <div className="flex flex-col h-full w-full bg-zinc-100 dark:bg-zinc-950">
            <div className="flex items-center justify-between px-4 py-3 bg-white dark:bg-zinc-900 border-b border-zinc-200/50 dark:border-zinc-800/50">
                <div className="flex items-center gap-2 font-medium text-sm">
                    <FileJson className="h-4 w-4 text-zinc-600 dark:text-zinc-400" />
                    Response Composer
                </div>
                <div className="flex items-center gap-2 text-xs">
                    {error ? (
                        <span className="flex items-center gap-1 text-red-500"><AlertCircle className="h-3 w-3" /> Invalid JSON</span>
                    ) : (
                        <span className="flex items-center gap-1 text-emerald-500"><div className="h-2 w-2 rounded-full bg-emerald-500" /> Valid JSON</span>
                    )}
                </div>
            </div>

            <div className="flex-1 p-4 relative">
                <textarea
                    value={value}
                    onChange={e => setValue(e.target.value)}
                    className={clsx(
                        "w-full h-full resize-none bg-white dark:bg-zinc-900 p-4 rounded-xl font-mono text-sm shadow-inner focus:outline-none focus:ring-2 focus:ring-zinc-400/50 dark:ring-zinc-600/50 border transition-colors",
                        error ? "border-red-500/50" : "border-zinc-200/50 dark:border-zinc-800/50"
                    )}
                    placeholder="Type JSON response here..."
                    spellCheck={false}
                />
            </div>

            <div className="p-4 bg-white dark:bg-zinc-900 border-t border-zinc-200/50 dark:border-zinc-800/50 flex flex-col gap-4">
                <div className="flex flex-wrap gap-2 justify-center pb-4 border-b border-zinc-100 dark:border-zinc-800/50">
                    {templates.map(t => (
                        <button
                            key={t.id}
                            onClick={() => applyTemplate(t.id)}
                            className="px-3 py-1.5 text-xs font-semibold bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-lg transition-colors text-zinc-700 dark:text-zinc-300 flex items-center gap-2"
                        >
                            <kbd className="font-mono text-[10px] text-zinc-500 border border-zinc-300 dark:border-zinc-600 rounded px-1.5">{t.shortcut.toUpperCase()}</kbd>
                            {t.label}
                        </button>
                    ))}
                </div>

                <div className="flex items-center justify-between gap-4">
                    <button
                        onClick={handleSave}
                        disabled={!!error || !value.trim()}
                        className="flex-1 flex items-center justify-center gap-2 p-3 text-sm font-bold rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 hover:bg-blue-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        <Save className="h-4 w-4" /> Save Fixture <kbd className="hidden sm:inline-block ml-2 text-xs opacity-60 font-mono">Ctrl+F</kbd>
                    </button>

                    <button
                        onClick={handleSend}
                        disabled={!!error || !value.trim() || submitting}
                        className="flex-1 flex items-center justify-center gap-2 p-3 text-sm font-bold rounded-xl bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-50 disabled:cursor-not-allowed transition-transform transform active:scale-95 shadow-lg "
                    >
                        {submitting ? 'Sending...' : <><Send className="h-4 w-4" /> Send Response <kbd className="hidden sm:inline-block ml-2 text-xs opacity-60 font-mono">Ctrl+S</kbd></>}
                    </button>
                </div>
            </div>
        </div>
    );
}
