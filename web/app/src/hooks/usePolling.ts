import { useCallback, useEffect, useRef, useState } from "react";

export function usePolling() {
  const [isPolling, setIsPolling] = useState(false);
  const ref = useRef<ReturnType<typeof setInterval> | null>(null);

  const stop = useCallback(() => {
    if (ref.current !== null) {
      clearInterval(ref.current);
      ref.current = null;
    }
    setIsPolling(false);
  }, []);

  const start = useCallback(
    (callback: () => void, intervalMs: number) => {
      stop();
      setIsPolling(true);
      callback();
      ref.current = setInterval(callback, intervalMs);
    },
    [stop],
  );

  useEffect(() => {
    return () => {
      if (ref.current !== null) {
        clearInterval(ref.current);
      }
    };
  }, []);

  return { isPolling, start, stop };
}
