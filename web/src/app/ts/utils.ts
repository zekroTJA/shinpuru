/** @format */

export function toHexClr(n: number, op = 1): string {
  if (n === 0) {
    return '#00000000';
  }

  const clr = n.toString(16);
  const opacity = Math.floor(op * 255).toString(16);

  return `#${clr}${opacity}`;
}
