export function getCryptoRandomString(ln: number): string {
  const array = new Uint8Array(ln);
  crypto.getRandomValues(array);
  return encode(array);
}

const DEFAULT_CHARSET = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';

function encode(v: Uint8Array, charset = DEFAULT_CHARSET): string {
  return Array.from(v)
    .map((c) => charset.charAt(Math.floor((c / 256) * charset.length)))
    .join('');
}
