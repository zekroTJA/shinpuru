type Entity = { id: string };

export const listReducer = <T extends Entity>(
  state: T[],
  [type, payload]: ['set', T[]] | ['add', T | T[]] | ['remove', T],
) => {
  switch (type) {
    case 'set':
      return payload;
    case 'add':
      return [...state, ...(Array.isArray(payload) ? payload : [payload])];
    case 'remove':
      return state.filter((e) => e.id !== payload.id);
    default:
      return state;
  }
};
