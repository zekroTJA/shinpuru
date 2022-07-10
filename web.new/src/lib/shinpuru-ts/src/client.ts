import {
  AuthClient,
  ChannelsClient,
  EtcClient,
  GlobalSettingsClient,
  GuildsClient,
  ReportsClient,
  SearchClient,
  TokensClient,
  UnbanRequestsClient,
  UsersClient,
  UtilClient,
  VerificationsClient,
} from './bindings';
import { HttpClient } from './httpclient';

export class Client extends HttpClient {
  etc = new EtcClient(this);
  util = new UtilClient(this);
  auth = new AuthClient(this);
  search = new SearchClient(this);
  tokens = new TokensClient(this);
  settings = new GlobalSettingsClient(this);
  reports = new ReportsClient(this);
  guilds = new GuildsClient(this);
  unbanrequests = new UnbanRequestsClient(this);
  channels = new ChannelsClient(this);
  verification = new VerificationsClient(this);
  users = new UsersClient(this);

  constructor(endpoint: string = '/api') {
    super(endpoint);
  }

  public get clientEndpoint(): string {
    return this.endpoint;
  }
}
