import { useState, useMemo, useEffect } from 'react';
import { Layout } from './components/Layout';
import { RequestList } from './components/RequestList';
import { RequestDetail } from './components/RequestDetail';
import { ResponseComposer } from './components/ResponseComposer';
import { useWebSocket } from './hooks/useWebSocket';
import { useRequests } from './hooks/useRequests';
import { api } from './lib/api';
import type { RequestDetail as RequestDetailType } from './lib/types';

export default function App() {
  const { requests, loading, handleWsEvent } = useRequests();
  const { status: wsStatus } = useWebSocket(handleWsEvent);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detailCache, setDetailCache] = useState<Record<string, RequestDetailType>>({});
  const [loadingDetail, setLoadingDetail] = useState(false);

  const pendingCount = useMemo(() =>
    requests.filter(r => r.status === 'pending').length
    , [requests]);

  const selectedItem = useMemo(() =>
    requests.find(r => r.id === selectedId) || null
    , [requests, selectedId]);

  const selectedDetail = selectedId ? detailCache[selectedId] : null;

  useEffect(() => {
    if (!selectedId || detailCache[selectedId]) return;

    setLoadingDetail(true);
    api.getRequest(selectedId)
      .then(detail => {
        setDetailCache(prev => ({ ...prev, [selectedId]: detail }));
      })
      .catch(console.error)
      .finally(() => setLoadingDetail(false));
  }, [selectedId, detailCache]);

  const handleSent = (_: string) => {
    // Optionally auto-select the next pending request here
  };

  const handleSaved = (_: string) => {
    // Notification logic if needed, ResponseComposer handles saving
  };

  // Keyboard shortcut for navigating list (up/down)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        setSelectedId(null);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  return (
    <Layout wsStatus={wsStatus} pendingCount={pendingCount}>
      <div className="w-1/3 min-w-[300px] border-r border-slate-200 dark:border-slate-800 flex flex-col bg-white dark:bg-slate-900 overflow-hidden relative z-20">
        {loading && requests.length === 0 ? (
          <div className="flex-1 flex items-center justify-center p-8 text-slate-500">Loading requests...</div>
        ) : (
          <RequestList
            requests={requests}
            selectedId={selectedId}
            onSelect={setSelectedId}
          />
        )}
      </div>

      {selectedItem ? (
        <>
          <div className="w-1/3 min-w-[350px] shadow-[20px_0_40px_-15px_rgba(0,0,0,0.1)] dark:shadow-[20px_0_40px_-15px_rgba(0,0,0,0.5)] relative z-10 flex flex-col overflow-hidden">
            {loadingDetail && !selectedDetail ? (
              <div className="flex-1 flex items-center justify-center text-slate-500 bg-white dark:bg-slate-900">Loading details...</div>
            ) : selectedDetail ? (
              <RequestDetail req={selectedDetail} />
            ) : null}
          </div>

          <div className="flex-1 min-w-[400px] z-0 flex flex-col overflow-hidden bg-slate-100 dark:bg-slate-950">
            <ResponseComposer
              reqId={selectedItem.id}
              isAnswered={selectedItem.status !== 'pending'}
              onSent={handleSent}
              onSaved={handleSaved}
            />
          </div>
        </>
      ) : (
        <div className="flex-1 flex flex-col items-center justify-center text-slate-500 bg-slate-50 dark:bg-slate-950">
          <p className="text-lg font-medium text-slate-700 dark:text-slate-300">No request selected</p>
          <p className="text-sm mt-2 opacity-70">Select a request from the sidebar to inspect and respond.</p>
        </div>
      )}
    </Layout>
  );
}
