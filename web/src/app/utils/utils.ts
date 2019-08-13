/** @format */

import { Role } from '../api/api.models';

/** @format */

export function toHexClr(n: number, op = 1): string {
  if (!n) {
    return '#00000000';
  }

  const clr = n.toString(16);
  const opacity = Math.floor(op * 255).toString(16);

  return `#${clr}${opacity}`;
}

export function clone<T>(v: T): T {
  return JSON.parse(JSON.stringify(v)) as T;
}

export function topRole(roles: Role[], roleIDs: string[]): Role {
  if (!roleIDs || !roleIDs.length) {
    return { position: -1 } as Role;
  }

  const uRoles = roleIDs.map((rID) => roles.find((r) => r.id === rID) || null);

  let top = uRoles[0];
  uRoles
    .slice(1)
    .forEach((r) => (top = r && r.position > top.position ? r : top));
  return top;
}
