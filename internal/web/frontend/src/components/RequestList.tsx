import type { RequestItem } from '../lib/types';
import { StatusBadge } from './StatusBadge';
import clsx from 'clsx';
import { formatDistanceToNow } from 'date-fns';

export function RequestList({
    requests,
    selectedId,
    onSelect
}: {
    requests: RequestItem[];
    selectedId: string | null;
    onSelect: (id: string) => void;
}) {
    if (requests.length === 0) {
        return (
            <div className="flex h-full w-full items-center justify-center text-sm text-slate-500">
                No requests yet. Send a request to the mock server.
            </div>
        );
    }

    return (
        <div className="flex flex-col h-full w-full divide-y divide-slate-200 dark:divide-slate-800/50 overflow-y-auto">
            {requests.map(req => {
                const isSelected = req.id === selectedId;
                const targetDate = new Date(req.timestamp);
                const timeAgo = !isNaN(targetDate.getTime()) ? formatDistanceToNow(targetDate, { addSuffix: true }) : '';

                return (
                    <button
                        key={req.id}
                        onClick={() => onSelect(req.id)}
                        className={clsx(
                            "flex flex-col items-start gap-2 p-4 text-left transition-colors focus:outline-none focus:bg-slate-100 dark:focus:bg-slate-800/80",
                            isSelected ? "bg-slate-100 dark:bg-slate-800/80" : "hover:bg-slate-50 dark:hover:bg-slate-800/40"
                        )}
                    >
                        <div className="flex w-full items-center justify-between gap-2">
                            <StatusBadge status={req.status === 'responded' && req.via === 'fixture' ? 'auto' : req.status} />
                            <div className="text-xs text-slate-500 whitespace-nowrap">
                                {timeAgo}
                            </div>
                        </div>
                        <div className="flex w-full flex-col gap-1">
                            <div className="flex w-full justify-between items-baseline">
                                <span className="font-mono text-xs font-semibold text-slate-900 dark:text-slate-100">{req.id}</span>
                                <span className="text-xs font-medium text-slate-600 dark:text-slate-400 line-clamp-1">{req.model}</span>
                            </div>
                            <div className="text-sm text-slate-600 dark:text-slate-300 line-clamp-2 pr-4 leading-snug">
                                {req.preview || <span className="italic opacity-50">No text preview available</span>}
                            </div>
                        </div>
                    </button>
                );
            })}
        </div>
    );
}
