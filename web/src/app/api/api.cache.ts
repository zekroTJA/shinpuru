/** @format */

import { Observable } from 'rxjs';

/** @format */

export class CacheValue<T> {
  public d: number;

  constructor(public val: T) {
    this.d = new Date().getTime();
  }

  public isUpDoDate(maxAge: number) {
    return new Date().getTime() - this.d <= maxAge;
  }
}

export class CacheBucket<TKey, TVal> {
  private cache: Map<TKey, CacheValue<TVal>>;

  constructor(private maxAge: number) {
    this.cache = new Map();
  }

  public put(key: TKey, val: TVal) {
    this.cache.set(key, new CacheValue(val));
  }

  public get(key: TKey): TVal | null {
    const v = this.cache.get(key);
    if (!v) {
      return null;
    }
    return v.isUpDoDate(this.maxAge) ? v.val : null;
  }

  public clear() {
    this.cache = new Map();
  }

  public putFromPipe = (key: TKey) => (source: Observable<TVal>) =>
    new Observable<TVal>((observer) => {
      return source.subscribe({
        next: (x: TVal) => {
          this.put(key, x);
          observer.next(x);
        },
        error: (err) => observer.error(err),
        complete: () => observer.complete(),
      });
    });
}
