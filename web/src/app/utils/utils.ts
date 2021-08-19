/** @format */

import { Role, Member } from '../api/api.models';

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

export function rolePosDiff(roles: Role[], m1: Member, m2: Member): number {
  const rm1 = roles
    .filter((r) => m1.roles.includes(r.id))
    .sort((a, b) => b.position - a.position)[0];
  const rm2 = roles
    .filter((r) => m2.roles.includes(r.id))
    .sort((a, b) => b.position - a.position)[0];

  if (!rm1 || !rm2) {
    return null;
  }

  return rm1.position - rm2.position;
}

export function padNumber(n: number, minLen: number, padChar = '0'): string {
  const neg = n < 0;
  const nStr = Math.abs(n).toString();
  const diff = minLen - nStr.length;
  if (diff <= 0) {
    return nStr;
  }
  return (neg ? '-' : '') + padChar.repeat(diff) + nStr;
}

export function prefixNumner(n: number): string {
  return (n < 0 ? '-' : '+') + n.toString();
}
