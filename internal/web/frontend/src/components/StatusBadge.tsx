import { Check, Zap } from 'lucide-react';

export function StatusBadge({ status }: { status: string }) {
    if (status === 'pending') {
        return (
            <span className="inline-flex items-center gap-1.5 rounded-full bg-amber-500/10 px-2.5 py-0.5 text-xs font-medium text-amber-500 ring-1 ring-inset ring-amber-500/20">
                <span className="h-1.5 w-1.5 rounded-full bg-amber-500 animate-pulse"></span>
                PENDING
            </span>
        );
    }
    if (status === 'auto') {
        return (
            <span className="inline-flex items-center gap-1 rounded-full bg-blue-500/10 px-2.5 py-0.5 text-xs font-medium text-blue-400 ring-1 ring-inset ring-blue-500/20">
                <Zap className="h-3 w-3" />
                AUTO
            </span>
        );
    }
    return (
        <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2.5 py-0.5 text-xs font-medium text-emerald-400 ring-1 ring-inset ring-emerald-500/20">
            <Check className="h-3 w-3" />
            DONE
        </span>
    );
}
