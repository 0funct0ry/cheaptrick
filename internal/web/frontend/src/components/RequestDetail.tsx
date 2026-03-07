import type { RequestDetail as RequestDetailType } from '../lib/types';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { useState } from 'react';
import clsx from 'clsx';

function CollapsibleSection({ title, defaultOpen = true, children }: { title: string, defaultOpen?: boolean, children: React.ReactNode }) {
    const [open, setOpen] = useState(defaultOpen);
    return (
        <div className="flex flex-col border-b border-zinc-200/50 dark:border-zinc-800/50">
            <button
                onClick={() => setOpen(!open)}
                className="flex items-center gap-2 py-3 px-4 text-sm font-semibold hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors"
            >
                {open ? <ChevronDown className="h-4 w-4 text-zinc-500" /> : <ChevronRight className="h-4 w-4 text-zinc-500" />}
                {title}
            </button>
            {open && <div className="p-4 pt-0">{children}</div>}
        </div>
    );
}

function JsonViewer({ data }: { data: unknown }) {
    return (
        <pre className="text-[11px] leading-snug font-mono bg-zinc-100 dark:bg-zinc-900 p-4 rounded-lg overflow-x-auto text-zinc-800 dark:text-zinc-300 shadow-inner">
            {JSON.stringify(data, null, 2)}
        </pre>
    );
}

export function RequestDetail({ req }: { req: RequestDetailType }) {
    return (
        <div className="flex flex-col h-full w-full overflow-y-auto bg-white dark:bg-zinc-950 shadow-none ring-1 ring-zinc-200/50 dark:ring-zinc-800/50 z-10 border-x border-zinc-200/50 dark:border-zinc-800/50">
            <div className="p-4 border-b border-zinc-200/50 dark:border-zinc-800/50 bg-zinc-50/80 dark:bg-zinc-950/50 flex flex-col gap-2">
                <h2 className="text-lg font-bold font-mono text-zinc-900 dark:text-zinc-100">{req.id}</h2>
                <div className="text-sm text-zinc-600 dark:text-zinc-400">Model: <span className="font-semibold text-zinc-900 dark:text-zinc-200">{req.model}</span></div>
            </div>

            <div className="flex flex-col">
                {req.system_instruction && (
                    <CollapsibleSection title="System Instruction">
                        <div className="text-sm text-zinc-700 dark:text-zinc-300 italic whitespace-pre-wrap bg-indigo-50 dark:bg-indigo-950/30 p-4 rounded-lg border border-indigo-100 dark:border-indigo-900/50">
                            {req.system_instruction}
                        </div>
                    </CollapsibleSection>
                )}

                {req.contents && req.contents.length > 0 && (
                    <CollapsibleSection title="Contents">
                        <div className="flex flex-col gap-4">
                            {req.contents.map((msgRaw: unknown, i: number) => {
                                const msg = msgRaw as Record<string, unknown>;
                                const isUser = msg.role === 'user';
                                return (
                                    <div key={i} className={clsx("flex flex-col w-full max-w-[90%] gap-1", isUser ? "self-end items-end" : "self-start items-start")}>
                                        <span className="text-xs font-bold text-zinc-500 uppercase px-1">{msg.role as string}</span>
                                        <div className={clsx(
                                            "p-3 rounded-2xl text-sm whitespace-pre-wrap font-sans",
                                            isUser
                                                ? "bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100 rounded-tr-sm"
                                                : "bg-zinc-50 dark:bg-zinc-800/50 text-zinc-900 dark:text-zinc-100 rounded-tl-sm border border-zinc-200/50 dark:border-zinc-700"
                                        )}>
                                            {(msg.parts as Record<string, unknown>[])?.map((p: Record<string, unknown>, j: number): React.ReactNode => {
                                                if (p.text) return <span key={j}>{p.text as string}</span>;
                                                if (p.functionCall) {
                                                    const fc = p.functionCall as Record<string, unknown>;
                                                    return (
                                                        <div key={j} className="font-mono text-xs bg-black/10 dark:bg-black/30 p-2 rounded mt-1 overflow-x-auto break-all">
                                                            {fc.name as string}({JSON.stringify(fc.args)})
                                                        </div>
                                                    );
                                                }
                                                if (p.functionResponse) {
                                                    const fr = p.functionResponse as Record<string, unknown>;
                                                    return (
                                                        <div key={j} className="font-mono text-xs p-2 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 rounded mt-1 overflow-x-auto whitespace-pre-wrap break-all">
                                                            {JSON.stringify(fr.response, null, 2)}
                                                        </div>
                                                    );
                                                }
                                                return <span key={j}>[Unknown part]</span>;
                                            })}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </CollapsibleSection>
                )}

                {req.tools && req.tools.length > 0 && (
                    <CollapsibleSection title="Tools">
                        <JsonViewer data={req.tools} />
                    </CollapsibleSection>
                )}

                {!!req.generation_config && (
                    <CollapsibleSection title="Generation Config" defaultOpen={false}>
                        <JsonViewer data={req.generation_config} />
                    </CollapsibleSection>
                )}

                <CollapsibleSection title="Raw Body JSON" defaultOpen={false}>
                    <JsonViewer data={req.body} />
                </CollapsibleSection>

                {!!req.response && (
                    <CollapsibleSection title="Sent Response" defaultOpen={true}>
                        <JsonViewer data={req.response} />
                    </CollapsibleSection>
                )}
            </div>
        </div>
    );
}
