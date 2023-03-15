export const isAllowed = (p: string) => p.startsWith('+');
export const isDisallowed = (p: string) => !isAllowed(p);
