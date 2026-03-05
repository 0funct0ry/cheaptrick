import { useEffect } from 'react';

type ShortcutMap = {
    [key: string]: (e: KeyboardEvent) => void;
};

export function useKeyboardShortcuts(shortcuts: ShortcutMap, active: boolean = true) {
    useEffect(() => {
        if (!active) return;

        const handleKeyDown = (e: KeyboardEvent) => {
            let combo = '';
            if (e.ctrlKey || e.metaKey) combo += 'ctrl+';
            if (e.shiftKey) combo += 'shift+';
            if (e.altKey) combo += 'alt+';
            combo += e.key.toLowerCase();

            for (const [shortcut, handler] of Object.entries(shortcuts)) {
                if (shortcut.toLowerCase() === combo) {
                    e.preventDefault();
                    handler(e);
                    return;
                }
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [shortcuts, active]);
}
