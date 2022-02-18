export function getCryptoRandomString(ln: number): string {
  const array = new Uint8Array(ln);
  crypto.getRandomValues(array);
  return encode(array);
}

const ASCII_START = 33;
const ASCII_END = 126;

function encode(v: Uint8Array): string {
  v = v.map((c) =>
    Math.floor((c / 256) * (ASCII_END - ASCII_START + 1) + ASCII_START)
  );
  return String.fromCharCode(...v);
}
