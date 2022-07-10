export function getQueryParam(key: string): string | null {
  return new URLSearchParams(window.location.search).get(key);
}

export function setQueryParam(key: string, v: string) {
  const sp = new URLSearchParams(window.location.search);
  if (!v) sp.delete(key);
  else sp.set(key, v);
  const query = sp.toString();
  if (query.length === 0) window.history.pushState('', '', '/');
  else window.history.pushState('', '', `?${query}`);
}
