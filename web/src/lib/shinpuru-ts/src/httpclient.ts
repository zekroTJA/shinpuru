import { APIError } from './errors';
import { AccessTokenModel, ErrorReponse } from './models';

export type HttpMethod =
  | 'GET'
  | 'PUT'
  | 'POST'
  | 'DELETE'
  | 'PATCH'
  | 'OPTIONS';

export type HttpHeadersMap = { [key: string]: string };

export interface IHttpClient {
  req<T>(
    method: HttpMethod,
    path: string,
    body?: object,
    appendHeaders?: HttpHeadersMap
  ): Promise<T>;
}

export interface HttpClientOptions {
  authorization?: string;
  headers?: HttpHeadersMap;
}

interface AccessToken extends AccessTokenModel {
  expiresDate: Date;
}

export class HttpClient implements IHttpClient {
  private accessToken: AccessToken | null = null;

  constructor(
    protected endpoint: string,
    private options = {} as HttpClientOptions
  ) {}

  async req<T>(
    method: HttpMethod,
    path: string,
    body?: object,
    appendHeaders?: HttpHeadersMap
  ): Promise<T> {
    const headers = new Headers();
    headers.set('Content-Type', 'application/json');
    headers.set('Accept', 'application/json');
    kv<string, string>(this.options?.headers).forEach(([k, v]) =>
      headers.set(k, v)
    );
    kv<string, string>(headers).forEach(([k, v]) => headers.set(k, v));

    if (this.options.authorization)
      headers.set('Authorization', this.options.authorization);
    else if (this.accessToken) {
      if (Date.now() - this.accessToken.expiresDate.getTime() > 0)
        return await this.getAndSetAccessToken(() =>
          this.req(method, path, body, appendHeaders)
        );
      headers.set('Authorization', `accessToken ${this.accessToken.token}`);
    }
    const fullPath = `${this.endpoint}/${path}`.replace(/(?<=[^:])\/\//g, '/');
    console.log('fullpath', fullPath);
    const res = await window.fetch(fullPath, {
      method,
      headers,
      body: body ? JSON.stringify(body) : null,
      credentials: 'include',
    });

    if (res.status === 204) {
      return {} as T;
    }

    let data = {};
    try {
      data = await res.json();
    } catch {}

    if (
      res.status === 401 &&
      (data as ErrorReponse).error === 'invalid access token'
    ) {
      return await this.getAndSetAccessToken(() =>
        this.req(method, path, body, appendHeaders)
      );
    }

    if (res.status >= 400) throw new APIError(res, data as ErrorReponse);

    return data as T;
  }

  private async getAccessToken(): Promise<AccessTokenModel> {
    return this.req('POST', 'auth/accesstoken');
  }

  private async getAndSetAccessToken<T>(replay: () => Promise<T>): Promise<T> {
    const token = await this.getAccessToken();
    this.accessToken = token as AccessToken;
    this.accessToken.expiresDate = new Date(token.expires);
    return await replay();
  }
}

function kv<TKey, TVal>(m?: any) {
  if (!m) m = {};
  return Object.keys(m).map((k) => [k, m[k]]) as any as [TKey, TVal][];
}
