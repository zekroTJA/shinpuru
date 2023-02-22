import { HttpHeadersMap, HttpMethod, IHttpClient } from './httpclient';

import { Client } from './client';

export class SubClient implements IHttpClient {
  constructor(private client: Client, protected sub: string) {}

  req<TResp>(
    method: HttpMethod,
    path: string,
    body?: object,
    appendHeaders?: HttpHeadersMap,
  ): Promise<TResp> {
    return this.client.req(method, `${this.sub}/${path}`, body, appendHeaders);
  }

  protected get endpoint(): string {
    return `${this.client.clientEndpoint}/${this.sub}`;
  }
}
