import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Trash2, Plus, Save } from 'lucide-react';
import type { Fixture } from '../lib/types';
import clsx from 'clsx';

export function FixtureManager() {
    const [fixtures, setFixtures] = useState<Fixture[]>([]);
    const [selectedHash, setSelectedHash] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const [editingHash, setEditingHash] = useState('');
    const [editingContent, setEditingContent] = useState('');
    const [isSaving, setIsSaving] = useState(false);

    const [deleteModalVisible, setDeleteModalVisible] = useState(false);
    const [fixtureToDelete, setFixtureToDelete] = useState<string | null>(null);

    const [createModalVisible, setCreateModalVisible] = useState(false);
    const [newRequest, setNewRequest] = useState('');
    const [newContent, setNewContent] = useState('');
    const [isCreating, setIsCreating] = useState(false);

    const loadFixtures = () => {
        setLoading(true);
        api.getFixtures()
            .then(setFixtures)
            .catch(console.error)
            .finally(() => setLoading(false));
    };

    useEffect(() => {
        loadFixtures();
    }, []);

    useEffect(() => {
        if (!selectedHash) {
            setEditingHash('');
            setEditingContent('');
            return;
        }

        api.getFixture(selectedHash)
            .then(data => {
                setEditingHash(selectedHash);
                setEditingContent(JSON.stringify(data.content, null, 2));
            })
            .catch(console.error);
    }, [selectedHash]);

    const handleCreateNew = () => {
        setNewRequest('{\n  "contents": [\n    {\n      "parts": [\n        {\n          "text": "<prompt-text>"\n        }\n      ],\n      "role": "user"\n    }\n  ]\n}');
        setNewContent('{\n  "candidates": [\n    {\n      "content": {\n        "role": "model",\n        "parts": [\n          {\n            "text": "Your textual response here."\n          }\n        ]\n      },\n      "finishReason": "STOP"\n    }\n  ],\n  "usageMetadata": {\n    "promptTokenCount": 0,\n    "candidatesTokenCount": 0,\n    "totalTokenCount": 0\n  }\n}');
        setCreateModalVisible(true);
    };

    const handleCreateSubmit = async () => {
        if (!newRequest.trim() || !newContent.trim()) return;

        try {
            let parsedRequest, parsedResponse;
            try {
                parsedRequest = JSON.parse(newRequest);
                parsedResponse = JSON.parse(newContent);
            } catch {
                alert('Invalid JSON! Please check your syntax for both fields.');
                return;
            }

            setIsCreating(true);
            const result = await api.createFixture(null, parsedResponse, parsedRequest);
            await loadFixtures();

            setCreateModalVisible(false);
            if (result.hash) {
                setSelectedHash(result.hash);
            }
        } catch (error) {
            console.error('Failed to create fixture:', error);
            alert('Failed to create fixture');
        } finally {
            setIsCreating(false);
        }
    };

    const handleSave = async () => {
        if (!editingHash.trim() || !editingContent.trim()) return;

        try {
            // Validate JSON
            let parsed;
            try {
                parsed = JSON.parse(editingContent);
            } catch {
                alert('Invalid JSON! Please check your syntax.');
                return;
            }

            setIsSaving(true);
            await api.createFixture(editingHash, parsed);
            await loadFixtures();
        } catch (error) {
            console.error('Failed to save fixture:', error);
            alert('Failed to save fixture');
        } finally {
            setIsSaving(false);
        }
    };

    const handleDeleteClick = (hash: string) => {
        setFixtureToDelete(hash);
        setDeleteModalVisible(true);
    };

    const confirmDelete = async () => {
        if (!fixtureToDelete) return;

        try {
            await api.deleteFixture(fixtureToDelete);
            if (selectedHash === fixtureToDelete) {
                setSelectedHash(null);
            }
            await loadFixtures();
        } catch (error) {
            console.error('Failed to delete fixture:', error);
        } finally {
            setDeleteModalVisible(false);
            setFixtureToDelete(null);
        }
    };

    return (
        <div className="flex w-full h-full relative">
            {/* Sidebar List */}
            <div className="w-1/3 min-w-[300px] max-w-[400px] border-r border-slate-200 dark:border-slate-800 flex flex-col bg-white dark:bg-slate-900 shadow-sm z-10">
                <div className="p-4 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center">
                    <h2 className="font-semibold text-slate-800 dark:text-slate-200">Fixtures</h2>
                    <button
                        onClick={handleCreateNew}
                        className="p-1.5 rounded-md bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-100 dark:hover:bg-indigo-500/20 transition-colors"
                        title="Create New Fixture"
                    >
                        <Plus className="w-4 h-4" />
                    </button>
                </div>

                <div className="flex-1 overflow-y-auto">
                    {loading && fixtures.length === 0 ? (
                        <div className="p-8 text-center text-slate-500 text-sm">Loading fixtures...</div>
                    ) : fixtures.length === 0 ? (
                        <div className="p-8 text-center text-slate-500 text-sm">No fixtures found.<br />Create one to get started!</div>
                    ) : (
                        <div className="flex flex-col">
                            {fixtures.map(f => (
                                <div
                                    key={f.hash}
                                    onClick={() => setSelectedHash(f.hash)}
                                    className={clsx(
                                        "flex items-center justify-between p-3 border-b border-slate-100 dark:border-slate-800/50 cursor-pointer transition-colors group",
                                        selectedHash === f.hash
                                            ? "bg-indigo-50 dark:bg-indigo-500/10 border-l-4 border-l-indigo-500"
                                            : "hover:bg-slate-50 dark:hover:bg-slate-800 border-l-4 border-l-transparent"
                                    )}
                                >
                                    <div className="flex flex-col min-w-0 pr-4">
                                        <div className={clsx(
                                            "font-medium truncate",
                                            selectedHash === f.hash ? "text-indigo-700 dark:text-indigo-300" : "text-slate-700 dark:text-slate-300"
                                        )}>
                                            {f.hash}
                                        </div>
                                        <div className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                                            {(f.size / 1024).toFixed(1)} KB
                                        </div>
                                    </div>
                                    <button
                                        onClick={(e) => { e.stopPropagation(); handleDeleteClick(f.hash); }}
                                        className="p-1.5 rounded-md text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/30 opacity-0 group-hover:opacity-100 transition-all focus:opacity-100"
                                        title="Delete Fixture"
                                    >
                                        <Trash2 className="w-4 h-4" />
                                    </button>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            {/* Editor Area */}
            <div className="flex-1 flex flex-col bg-slate-50 dark:bg-slate-950 overflow-hidden">
                {selectedHash ? (
                    <div className="flex flex-col h-full max-w-5xl w-full mx-auto shadow-sm bg-white dark:bg-slate-900/50 border-x border-slate-200 dark:border-slate-800">
                        <div className="p-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between bg-white dark:bg-slate-900">
                            <div className="flex-1 mr-4">
                                <label className="block text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
                                    Fixture Name / Hash
                                </label>
                                <input
                                    type="text"
                                    value={editingHash}
                                    readOnly
                                    className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-md shadow-sm outline-none text-slate-500 dark:text-slate-400 font-mono text-sm cursor-not-allowed"
                                />
                            </div>
                            <button
                                onClick={handleSave}
                                disabled={isSaving || !editingHash.trim() || !editingContent.trim()}
                                className="mt-5 flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-md shadow-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                <Save className="w-4 h-4" />
                                {isSaving ? 'Saving...' : 'Save'}
                            </button>
                        </div>

                        <div className="flex-1 flex flex-col p-4 overflow-hidden">
                            <label className="block text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
                                JSON Response Payload
                            </label>
                            <textarea
                                value={editingContent}
                                onChange={(e) => setEditingContent(e.target.value)}
                                className="flex-1 w-full p-4 bg-slate-50 dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 rounded-lg shadow-inner outline-none focus:ring-2 focus:ring-indigo-500/50 font-mono text-sm text-slate-800 dark:text-slate-300 resize-none"
                                spellCheck={false}
                                placeholder="Enter valid JSON here..."
                            />
                        </div>
                    </div>
                ) : (
                    <div className="flex-1 flex flex-col items-center justify-center text-slate-500">
                        <div className="bg-white dark:bg-slate-900 p-8 rounded-xl shadow-sm border border-slate-200 dark:border-slate-800 flex flex-col items-center">
                            <p className="text-lg font-medium text-slate-700 dark:text-slate-300">No fixture selected</p>
                            <p className="text-sm mt-2 mb-6 opacity-70">Select a fixture from the sidebar to edit, or create a new one.</p>
                            <button
                                onClick={handleCreateNew}
                                className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300 font-medium rounded-md transition-colors"
                            >
                                <Plus className="w-4 h-4" />
                                Create New
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Create Fixture Modal */}
            {createModalVisible && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 text-left">
                    <div className="bg-white dark:bg-slate-900 rounded-xl shadow-xl w-full max-w-2xl overflow-hidden animate-in fade-in zoom-in duration-200 flex flex-col max-h-[90vh]">
                        <div className="p-4 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center bg-white dark:bg-slate-900">
                            <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100">Create New Fixture</h3>
                        </div>

                        <div className="p-4 flex flex-col gap-4 overflow-y-auto">
                            <div className="flex flex-col flex-1 min-h-[200px]">
                                <label className="block text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
                                    Request JSON
                                </label>
                                <textarea
                                    value={newRequest}
                                    onChange={(e) => setNewRequest(e.target.value)}
                                    className="flex-1 w-full p-4 bg-slate-50 dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 rounded-lg shadow-inner outline-none focus:ring-2 focus:ring-indigo-500/50 font-mono text-sm text-slate-800 dark:text-slate-300 resize-none min-h-[150px]"
                                    spellCheck={false}
                                    placeholder="Enter valid Request JSON here..."
                                />
                            </div>

                            <div className="flex flex-col flex-1 min-h-[300px]">
                                <label className="block text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
                                    Response JSON
                                </label>
                                <textarea
                                    value={newContent}
                                    onChange={(e) => setNewContent(e.target.value)}
                                    className="flex-1 w-full p-4 bg-slate-50 dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 rounded-lg shadow-inner outline-none focus:ring-2 focus:ring-indigo-500/50 font-mono text-sm text-slate-800 dark:text-slate-300 resize-none min-h-[250px]"
                                    spellCheck={false}
                                    placeholder="Enter valid Response JSON here..."
                                />
                            </div>
                        </div>

                        <div className="bg-slate-50 dark:bg-slate-950/50 px-6 py-4 flex items-center justify-end gap-3 border-t border-slate-100 dark:border-slate-800/50">
                            <button
                                onClick={() => setCreateModalVisible(false)}
                                className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleCreateSubmit}
                                disabled={isCreating || !newRequest.trim() || !newContent.trim()}
                                className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-md shadow-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                            >
                                <Save className="w-4 h-4" />
                                {isCreating ? 'Creating...' : 'Create'}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Delete Confirmation Modal */}
            {deleteModalVisible && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
                    <div className="bg-white dark:bg-slate-900 rounded-xl shadow-xl w-full max-w-sm overflow-hidden animate-in fade-in zoom-in duration-200">
                        <div className="p-6">
                            <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-2">Delete Fixture</h3>
                            <p className="text-slate-600 dark:text-slate-400 text-sm">
                                Are you sure you want to delete the fixture <span className="font-mono bg-slate-100 dark:bg-slate-800 px-1 rounded text-red-600 dark:text-red-400 break-all">{fixtureToDelete}</span>?
                                This action cannot be undone.
                            </p>
                        </div>
                        <div className="bg-slate-50 dark:bg-slate-950/50 px-6 py-4 flex items-center justify-end gap-3 border-t border-slate-100 dark:border-slate-800/50">
                            <button
                                onClick={() => setDeleteModalVisible(false)}
                                className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={confirmDelete}
                                className="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md shadow-sm transition-colors"
                            >
                                Confirm Delete
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
