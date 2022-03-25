/** @format */

export default class LocalStorageUtil {
  public static set<T>(key: string, value: T) {
    const valJSON = JSON.stringify(value);
    window.localStorage.setItem(key, valJSON);
  }

  public static get<T>(key: string, def?: T): T | null {
    const valJSON = window.localStorage.getItem(key);
    if (!valJSON) {
      return def;
    }
    return JSON.parse(valJSON) as T;
  }

  public static remove(key: string) {
    window.localStorage.removeItem(key);
  }
}
