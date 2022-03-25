// This is here because you COULD simply use a positive lookbehind
// regex to replace the double slashes, but because webkit wants to
// be the next Internet Explorer, they donÄt support it. :)
export function replaceDoublePath(url: string): string {
  const split = url.split('://');
  split[split.length - 1] = split[split.length - 1].replace(/\/\//g, '/');
  return split.join('://');
}
