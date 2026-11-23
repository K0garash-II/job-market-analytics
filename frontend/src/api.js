import { useEffect, useState } from "react";

export function useMarket(path) {
  const [revision, refresh] = useState(0);
  const [result, setResult] = useState({});
  useEffect(() => {
    if (!path) return;
    const controller = new AbortController();
    let timer;
    async function read() {
      try {
        const response = await fetch("/api/v1/" + path, {
          signal: controller.signal,
        });
        if (!response.ok) throw new Error(String(response.status));
        const data = await response.json();
        if (!controller.signal.aborted) setResult({ path, revision, data });
      } catch {
        if (!controller.signal.aborted)
          setResult({ path, revision, error: true });
      } finally {
        if (!controller.signal.aborted) timer = setTimeout(read, 5 * 60 * 1000);
      }
    }
    read();
    return () => {
      controller.abort();
      clearTimeout(timer);
    };
  }, [path, revision]);
  const current =
    result.path === path && result.revision === revision ? result : {};
  return {
    ...current,
    loading: !current.data && !current.error,
    retry: () => refresh((n) => n + 1),
  };
}
export const number = (value) => value.toLocaleString("ru-RU");
export const sourceName = (name) =>
  ({ habr: "Хабр Карьера", remoteok: "Remote OK", hh: "hh.ru" })[name] || name;
export const dayLabel = (day) =>
  new Date(day + "T12:00:00").toLocaleDateString("ru-RU", {
    day: "numeric",
    month: "short",
  });
