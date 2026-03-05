import type { RequestDetail as RequestDetailType } from '../lib/types';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { useState } from 'react';
import clsx from 'clsx';

function CollapsibleSection({ title, defaultOpen = true, children }: { title: string, defaultOpen?: boolean, children: React.ReactNode }) {
    const [open, setOpen] = useState(defaultOpen);
    return (
        <div className="flex flex-col border-b border-slate-200 dark:border-slate-800">
            <button
                onClick={() => setOpen(!open)}
                className="flex items-center gap-2 py-3 px-4 text-sm font-semibold hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors"
            >
                {open ? <ChevronDown className="h-4 w-4 text-slate-500" /> : <ChevronRight className="h-4 w-4 text-slate-500" />}
                {title}
            </button>
            {open && <div className="p-4 pt-0">{children}</div>}
        </div>
    );
}

function JsonViewer({ data }: { data: any }) {
    return (
        <pre className="text-[11px] leading-snug font-mono bg-slate-100 dark:bg-slate-900 p-4 rounded-lg overflow-x-auto text-slate-800 dark:text-slate-300 shadow-inner">
            {JSON.stringify(data, null, 2)}
        </pre>
    );
}

export function RequestDetail({ req }: { req: RequestDetailType }) {
    return (
        <div className="flex flex-col h-full w-full overflow-y-auto bg-white dark:bg-slate-900 shadow-xl z-10 border-x border-slate-200 dark:border-slate-800">
            <div className="p-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50/80 dark:bg-slate-950/50 flex flex-col gap-2">
                <h2 className="text-lg font-bold font-mono text-slate-900 dark:text-slate-100">{req.id}</h2>
                <div className="text-sm text-slate-600 dark:text-slate-400">Model: <span className="font-semibold text-slate-900 dark:text-slate-200">{req.model}</span></div>
            </div>

            <div className="flex flex-col">
                {req.system_instruction && (
                    <CollapsibleSection title="System Instruction">
                        <div className="text-sm text-slate-700 dark:text-slate-300 italic whitespace-pre-wrap bg-indigo-50 dark:bg-indigo-950/30 p-4 rounded-lg border border-indigo-100 dark:border-indigo-900/50">
                            {req.system_instruction}
                        </div>
                    </CollapsibleSection>
                )}

                {req.contents && req.contents.length > 0 && (
                    <CollapsibleSection title="Contents">
                        <div className="flex flex-col gap-4">
                            {req.contents.map((msg: any, i: number) => {
                                const isUser = msg.role === 'user';
                                return (
                                    <div key={i} className={clsx("flex flex-col w-full max-w-[90%] gap-1", isUser ? "self-end items-end" : "self-start items-start")}>
                                        <span className="text-xs font-bold text-slate-500 uppercase px-1">{msg.role}</span>
                                        <div className={clsx(
                                            "p-3 rounded-2xl text-sm whitespace-pre-wrap font-sans",
                                            isUser
                                                ? "bg-amber-500 text-slate-900 rounded-tr-sm"
                                                : "bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-slate-100 rounded-tl-sm border border-slate-200 dark:border-slate-700"
                                        )}>
                                            {msg.parts?.map((p: any, j: number) => {
                                                if (p.text) return <span key={j}>{p.text}</span>;
                                                if (p.functionCall) return (
                                                    <div key={j} className="font-mono text-xs bg-black/10 dark:bg-black/30 p-2 rounded mt-1">
                                                        {p.functionCall.name}({JSON.stringify(p.functionCall.args)})
                                                    </div>
                                                );
                                                if (p.functionResponse) return (
                                                    <div key={j} className="font-mono text-xs p-2 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 rounded mt-1">
                                                        {JSON.stringify(p.functionResponse.response)}
                                                    </div>
                                                );
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

                {req.generation_config && (
                    <CollapsibleSection title="Generation Config" defaultOpen={false}>
                        <JsonViewer data={req.generation_config} />
                    </CollapsibleSection>
                )}

                <CollapsibleSection title="Raw Body JSON" defaultOpen={false}>
                    <JsonViewer data={req.body} />
                </CollapsibleSection>

                {req.response && (
                    <CollapsibleSection title="Sent Response" defaultOpen={true}>
                        <JsonViewer data={req.response} />
                    </CollapsibleSection>
                )}
            </div>
        </div>
    );
}
