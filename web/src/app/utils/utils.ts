/** @format */

export function toHexClr(n: number, op = 1): string {
  if (!n) {
    return '#00000000';
  }

  const clr = n.toString(16);
  const opacity = Math.floor(op * 255).toString(16);

  return `#${clr}${opacity}`;
}

export function permLvlColor(lvl: number): string {
  if (lvl < 1) {
    return '#424242';
  } else if (lvl < 3) {
    return '#0288D1';
  } else if (lvl < 5) {
    return '#689F38';
  } else if (lvl < 7) {
    return '#FFA000';
  } else if (lvl < 9) {
    return '#E64A19';
  } else if (lvl < 11) {
    return '#d32f2f';
  }

  return '#F50057';
}
