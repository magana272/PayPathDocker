const store = new Map();

export const cache = {
  get(key) {
    return store.get(key) ?? null;
  },

  set(key, data) {
    store.set(key, data);
  },

  invalidate(...keys) {
    if (keys.length === 0) {
      store.clear();
      return;
    }
    for (const key of keys) {
      if (key.endsWith("*")) {
        const prefix = key.slice(0, -1);
        for (const k of store.keys()) {
          if (k.startsWith(prefix)) store.delete(k);
        }
      } else {
        store.delete(key);
      }
    }
  },
};

export function emitRefresh() {
  if (typeof window !== "undefined") {
    window.dispatchEvent(new Event("paypath:refresh"));
  }
}