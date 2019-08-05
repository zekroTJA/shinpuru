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
