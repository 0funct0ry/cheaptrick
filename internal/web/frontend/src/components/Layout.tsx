import { useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { Moon, Sun } from 'lucide-react';
import type { WsStatus } from '../lib/types';
import clsx from 'clsx';

export function Layout({
    children,
    wsStatus,
    pendingCount
}: {
    children: ReactNode;
    wsStatus: WsStatus;
    pendingCount: number;
}) {
    const [theme, setTheme] = useState<'dark' | 'light'>('dark');

    useEffect(() => {
        if (theme === 'dark') {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }, [theme]);

    const wsColor =
        wsStatus === 'connected' ? 'bg-emerald-500' :
            wsStatus === 'connecting' ? 'bg-amber-500 animate-pulse' :
                wsStatus === 'reconnecting' ? 'bg-amber-500 animate-bounce' : 'bg-red-500';

    return (
        <div className="flex h-screen flex-col bg-zinc-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-200 overflow-hidden">
            {/* Header */}
            <header className="flex h-12 shrink-0 items-center justify-between border-b border-zinc-200/50 dark:border-zinc-800/50 bg-white dark:bg-zinc-950 px-4">
                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-2 font-semibold text-lg">
                        🎭 <span className="hidden sm:inline">Cheaptrick</span>
                    </div>
                    {pendingCount > 0 && (
                        <span className="inline-flex items-center rounded-full bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100 px-2.5 py-0.5 text-xs font-bold ring-1 ring-inset ring-zinc-300 dark:ring-zinc-700 animate-pulse">
                            {pendingCount} pending
                        </span>
                    )}
                </div>

                <div className="flex items-center gap-4 text-sm font-medium">
                    <div className="flex items-center gap-2">
                        <span className={clsx("h-2 w-2 rounded-full", wsColor)} />
                        <span className="hidden sm:inline capitalize">{wsStatus}</span>
                    </div>

                    <button
                        onClick={() => setTheme(t => t === 'dark' ? 'light' : 'dark')}
                        className="p-1.5 rounded-md text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
                    >
                        {theme === 'dark' ? <Moon className="h-4 w-4" /> : <Sun className="h-4 w-4" />}
                    </button>
                </div>
            </header>

            {/* Main Content Area */}
            <main className="flex flex-1 overflow-hidden relative">
                {children}
            </main>
        </div>
    );
}
