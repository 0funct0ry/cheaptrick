import type { RequestItem } from '../lib/types';
import { StatusBadge } from './StatusBadge';
import clsx from 'clsx';
import { formatDistanceToNow } from 'date-fns';
import { Trash2, Trash } from 'lucide-react';
import { api } from '../lib/api';

export function RequestList({
    requests,
    selectedId,
    onSelect
}: {
    requests: RequestItem[];
    selectedId: string | null;
    onSelect: (id: string) => void;
}) {
    const handleDelete = async (e: React.MouseEvent, id: string) => {
        e.stopPropagation();
        try {
            await api.deleteRequest(id);
        } catch (err) {
            console.error(err);
        }
    };

    const handleClearAll = async () => {
        try {
            await api.clearRequests();
        } catch (err) {
            console.error(err);
        }
    };

    const hasResponded = requests.some(r => r.status !== 'pending');

    if (requests.length === 0) {
        return (
            <div className="flex h-full w-full items-center justify-center text-sm text-zinc-500">
                No requests yet. Send a request to the mock server.
            </div>
        );
    }

    return (
        <div className="flex flex-col h-full w-full">
            {hasResponded && (
                <div className="flex justify-end p-2 border-b border-zinc-200/50 dark:border-zinc-800/50 shrink-0">
                    <button
                        onClick={handleClearAll}
                        title="Clear all responded requests"
                        className="flex items-center gap-1.5 px-2 py-1 text-xs font-medium text-zinc-500 hover:text-red-500 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded transition-colors"
                    >
                        <Trash className="w-4 h-4" />
                        <span>Clear All</span>
                    </button>
                </div>
            )}

            <div className="flex flex-col flex-1 divide-y divide-zinc-200 dark:divide-zinc-800/50 overflow-y-auto">
                {requests.map(req => {
                    const isSelected = req.id === selectedId;
                    const targetDate = new Date(req.timestamp);
                    const timeAgo = !isNaN(targetDate.getTime()) ? formatDistanceToNow(targetDate, { addSuffix: true }) : '';

                    return (
                        <button
                            key={req.id}
                            onClick={() => onSelect(req.id)}
                            className={clsx(
                                "flex flex-col items-start gap-2 p-4 text-left transition-colors focus:outline-none focus:bg-zinc-100 dark:focus:bg-zinc-800/80",
                                isSelected ? "bg-zinc-100 dark:bg-zinc-800/80" : "hover:bg-zinc-50 dark:hover:bg-zinc-800/40"
                            )}
                        >
                            <div className="flex w-full items-center justify-between gap-2">
                                <StatusBadge status={req.status === 'responded' && req.via === 'fixture' ? 'auto' : req.status} />
                                <div className="flex items-center gap-2">
                                    <div className="text-xs text-zinc-500 whitespace-nowrap">
                                        {timeAgo}
                                    </div>
                                    {req.status !== 'pending' && (
                                        <button
                                            onClick={(e) => handleDelete(e, req.id)}
                                            title="Delete request"
                                            className="p-1 text-zinc-400 hover:text-red-500 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded transition-colors"
                                        >
                                            <Trash2 className="w-3.5 h-3.5" />
                                        </button>
                                    )}
                                </div>
                            </div>
                            <div className="flex w-full flex-col gap-1">
                                <div className="flex w-full justify-between items-baseline">
                                    <span className="font-mono text-xs font-semibold text-zinc-900 dark:text-zinc-100">{req.id}</span>
                                    <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400 line-clamp-1">{req.model}</span>
                                </div>
                                <div className="text-sm text-zinc-600 dark:text-zinc-300 line-clamp-2 pr-4 leading-snug">
                                    {req.preview || <span className="italic opacity-50">No text preview available</span>}
                                </div>
                            </div>
                        </button>
                    );
                })}
            </div>
        </div>
    );
}
