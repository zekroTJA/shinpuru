/** @format */

import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import { catchError, map, share, tap } from 'rxjs/operators';
import { HttpClient, HttpParams } from '@angular/common/http';
import {
  User,
  ListReponse,
  Guild,
  Member,
  Report,
  PermissionResponse,
  GuildSettings,
  PermissionsUpdate,
  ReportRequest,
  ReasonRequest,
  Presence,
  InviteSettingsResponse,
  InviteSettingsRequest,
  Count,
  SystemInfo,
  APIToken,
  GuildBackup,
  GuildScoreboardEntry,
  CommandInfo,
  KarmaSettings,
  AntiraidSettings,
  JoinlogEntry,
  UnbanRequest,
  UserSettingsOTA,
  GuildStarboardEntry,
  AccessTokenModel,
  KarmaRule,
  GuildLogEntry,
  State,
  SearchResult,
  GuildSettingsApi,
  MessageEmbed,
  AntiraidAction,
  ChannelWithPermissions,
  VerificationSiteKey,
  GuildSettingsVerification,
  CodeExecSettings,
} from './api.models';
import { environment } from 'src/environments/environment';
import { ToastService } from '../components/toast/toast.service';
import { CacheBucket } from './api.cache';
import { Router } from '@angular/router';
import { NO_LOGIN_ROUTES } from '../utils/consts';

/** @format */

@Injectable({
  providedIn: 'root',
})
export class APIService {
  private rootURL = '';

  private accessToken: AccessTokenModel;

  private readonly cacheMembers = new CacheBucket<string, Member>(
    10 * 60 * 1000
  );
  private readonly cacheUsers = new CacheBucket<string, User>(10 * 60 * 1000);
  private readonly cacheGuilds = new CacheBucket<string, Guild>(30 * 1000);
  private readonly cachePermissionsAllowed = new CacheBucket<string, string[]>(
    3600 * 1000
  );

  private readonly defopts = (obj?: object) => {
    const defopts = {
      withCredentials: true,
      headers: {},
    };

    if (obj) {
      Object.keys(obj).forEach((k) => {
        defopts[k] = obj[k];
      });
    }

    return defopts;
  };

  private readonly rcAPI = (rc: string = '') =>
    `${this.rootURL}/api${rc ? '/' + rc : ''}`;

  private readonly rcAuth = (rc: string = '') =>
    `${this.rcAPI('auth')}${rc ? '/' + rc : ''}`;

  private readonly rcSearch = (rc: string = '') =>
    `${this.rcAPI('search')}${rc ? '/' + rc : ''}`;

  private readonly rcGuilds = (guildID: string = '') =>
    `${this.rcAPI('guilds')}${guildID ? '/' + guildID : ''}`;

  private readonly rcGuildMembers = (guildID: string, memberID: string = '') =>
    `${this.rcGuilds(guildID)}/${memberID ? memberID : 'members'}`;

  private readonly rcGuildMembersPermissions = (
    guildID: string,
    memberID: string
  ) => `${this.rcGuildMembers(guildID, memberID)}/permissions`;

  private readonly rcGuildMembersPermissionsAllowed = (
    guildID: string,
    memberID: string
  ) => `${this.rcGuildMembersPermissions(guildID, memberID)}/allowed`;

  private readonly rcGuildReports = (guildID: string) =>
    `${this.rcGuilds(guildID)}/reports`;

  private readonly rcGuildMemberReports = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/reports`;

  private readonly rcGuildReportsCount = (guildID: string) =>
    `${this.rcGuildReports(guildID)}/count`;

  private readonly rcGuildMemberReportsCount = (
    guildID: string,
    memberID: string
  ) => `${this.rcGuildMemberReports(guildID, memberID)}/count`;

  private readonly rcReports = (reportID: string, rc: string = '') =>
    `${this.rcAPI('reports')}/${reportID}${rc ? '/' + rc : ''}`;

  private readonly rcGuildSettings = (guildID: string) =>
    `${this.rcGuilds(guildID)}/settings`;

  private readonly rcUserSettings = (rc: string) =>
    `${this.rcAPI('usersettings')}${rc ? '/' + rc : ''}`;

  private readonly rcGuildSettingsKarma = (guildID: string) =>
    `${this.rcGuildSettings(guildID)}/karma`;

  private readonly rcGuildSettingsKarmaBlocklist = (
    guildID: string,
    rc: string = ''
  ) => `${this.rcGuildSettingsKarma(guildID)}/blocklist${rc ? '/' + rc : ''}`;

  private readonly rcGuildSettingsKarmaRules = (
    guildID: string,
    rc: string = ''
  ) => `${this.rcGuildSettingsKarma(guildID)}/rules${rc ? '/' + rc : ''}`;

  private readonly rcGuildSettingsLogs = (guildID: string, rc: string = '') =>
    `${this.rcGuildSettings(guildID)}/logs${rc ? '/' + rc : ''}`;

  private readonly rcGuildSettingsAntiraid = (guildID: string) =>
    `${this.rcGuildSettings(guildID)}/antiraid`;

  private readonly rcGuildSettingsAntiraidAction = (guildID: string) =>
    `${this.rcGuildSettingsAntiraid(guildID)}/action`;

  private readonly rcGuildSettingsFlushData = (guildID: string) =>
    `${this.rcGuildSettings(guildID)}/flushguilddata`;

  private readonly rcGuildSettingsApi = (guildID: string) =>
    `${this.rcGuildSettings(guildID)}/api`;

  private readonly rcGuildSettingsVerification = (guildID: string) =>
    `${this.rcGuildSettings(guildID)}/verification`;

  private readonly rcGuildSettingsCodeExec = (guildID: string) =>
    `${this.rcGuildSettings(guildID)}/codeexec`;

  private readonly rcGuildPermissions = (guildID: string) =>
    `${this.rcGuilds(guildID)}/permissions`;

  private readonly rcGuildBackups = (guildID: string, rc: string = '') =>
    `${this.rcGuilds(guildID)}/backups${rc ? '/' + rc : ''}`;

  private readonly rcGuildInviteBlock = (guildID: string) =>
    `${this.rcGuilds(guildID)}/inviteblock`;

  private readonly rcGuildScoreboard = (guildID: string) =>
    `${this.rcGuilds(guildID)}/scoreboard`;

  private readonly rcGuildStarboard = (guildID: string) =>
    `${this.rcGuilds(guildID)}/starboard`;

  public readonly rcGuildAntiraidJoinlog = (guildID: string) =>
    `${this.rcGuilds(guildID)}/antiraid/joinlog`;

  private readonly rcGuildMemberKick = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/kick`;

  private readonly rcGuildMemberBan = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/ban`;

  private readonly rcGuildMemberMute = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/mute`;

  private readonly rcGuildMemberUnmute = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/unmute`;

  private readonly rcSetting = (rc: string = '') =>
    `${this.rcAPI('settings')}${rc ? '/' + rc : ''}`;

  private readonly rcUtil = (rc: string = '') =>
    `${this.rcAPI('util')}${rc ? '/' + rc : ''}`;

  private readonly rcGuildUnbanRequest = (guildId: string, id: string = '') =>
    `${this.rcGuilds(guildId)}/unbanrequests${id ? '/' + id : ''}`;

  private readonly rcGuildMemberUnbanRequest = (
    guildId: string,
    memberId: string,
    id: string = ''
  ) =>
    `${this.rcGuildMembers(guildId, memberId)}/unbanrequests${
      id ? '/' + id : ''
    }`;

  private readonly rcUnbanRequests = (rc: string = '') =>
    `${this.rcAPI('unbanrequests')}${rc ? '/' + rc : ''}`;

  private readonly rcChannels = (
    guildId: string,
    chanId: string = '',
    msgId: string = ''
  ) =>
    `${this.rcAPI('channels')}/${guildId}${chanId ? '/' + chanId : ''}${
      msgId ? '/' + msgId : ''
    }`;

  private readonly rcUsers = (rc: string = '') =>
    `${this.rcAPI('users')}${rc ? '/' + rc : ''}`;

  private readonly rcVerification = (rc: string = '') =>
    `${this.rcAPI('verification')}${rc ? '/' + rc : ''}`;

  private readonly errorCatcher = (err) => {
    if (err instanceof TypeError) {
      return of({});
    }
    if (err.status === 401) {
      let path = window.location.pathname;
      if (NO_LOGIN_ROUTES.find((r) => path.startsWith(r))) return;
      if (!(path?.length > 0)) path = null;
      console.log('api', path);
      this.router.navigate(['/login'], {
        queryParams: {
          redirect: path,
        },
      });
      return of(null);
    }

    const msg = err?.error?.error ?? err.message;
    this.toasts.push(msg, 'Request Error', 'error', 10000);
    return throwError(err);
  };

  constructor(
    private http: HttpClient,
    private toasts: ToastService,
    private router: Router
  ) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }

  public getStoredAccessToken(): AccessTokenModel {
    return this.accessToken;
  }

  public getRcGuildBackupDownload(
    guildID: string,
    backupID: string,
    otaToken: string
  ): string {
    return `${this.rcGuildBackups(
      guildID
    )}/${backupID}/download?ota_token=${otaToken}`;
  }

  public getAndSetAccessToken(): Observable<AccessTokenModel> {
    return this.http
      .post<any>(this.rcAuth('accesstoken'), null, this.defopts())
      .pipe(catchError(this.errorCatcher))
      .pipe(tap((res) => (this.accessToken = res)));
  }

  public logout(): Observable<any> {
    return this.http
      .post<any>(this.rcAuth('logout'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getSelfUser(): Observable<User> {
    const u = this.cacheUsers.get('me');
    if (u) {
      return of(u);
    }

    return this.http.get<User>(this.rcAPI('me'), this.defopts()).pipe(
      this.cacheUsers.putFromPipe('me'),
      catchError((err) => {
        if (err.status !== 401) {
          return this.errorCatcher(err);
        }
      })
    );
  }

  public getGuilds(): Observable<Guild[]> {
    return this.http
      .get<ListReponse<Guild>>(this.rcGuilds(), this.defopts())
      .pipe(
        map((lr) => {
          return lr.data;
        }),
        catchError(this.errorCatcher)
      );
  }

  public getGuild(id: string): Observable<Guild> {
    const g = this.cacheGuilds.get(id);
    if (g) {
      return of(g);
    }

    return this.http
      .get<Guild>(this.rcGuilds(id), this.defopts())
      .pipe(this.cacheGuilds.putFromPipe(id), catchError(this.errorCatcher));
  }

  public getGuildMembers(
    guildID: string,
    after: string = '',
    limit: number = 0
  ): Observable<Member[]> {
    const opts = this.defopts({
      params: new HttpParams()
        .set('after', after)
        .set('limit', limit.toString()),
    });
    return this.http
      .get<ListReponse<Member>>(this.rcGuildMembers(guildID), opts)
      .pipe(
        map((lr) => lr.data),
        catchError(this.errorCatcher)
      );
  }

  public getGuildMember(
    guildID: string,
    memberID: string,
    ignoreError: boolean = false
  ): Observable<Member> {
    const m = this.cacheMembers.get(`${guildID}#${memberID}`);
    if (m) {
      return of(m);
    }

    return this.http
      .get<Member>(this.rcGuildMembers(guildID, memberID), this.defopts())
      .pipe(
        this.cacheMembers.putFromPipe(`${guildID}#${memberID}`),
        catchError(ignoreError ? (err) => of(null) : this.errorCatcher)
      );
  }

  public getPermissions(
    guildID: string,
    memberID: string
  ): Observable<string[]> {
    return this.http
      .get<PermissionResponse>(
        this.rcGuildMembersPermissions(guildID, memberID),
        this.defopts()
      )
      .pipe(
        map((r) => {
          return r.permissions;
        }),
        catchError(this.errorCatcher)
      );
  }

  public getPermissionsAllowed(
    guildID: string,
    memberID: string
  ): Observable<string[]> {
    const allowed = this.cachePermissionsAllowed.get(`${guildID}:${memberID}`);
    if (allowed) return of(allowed);

    return this.http
      .get<ListReponse<string>>(
        this.rcGuildMembersPermissionsAllowed(guildID, memberID),
        this.defopts()
      )
      .pipe(
        map((l) => l.data),
        tap((l) =>
          this.cachePermissionsAllowed.put(`${guildID}:${memberID}`, l)
        ),
        catchError(this.errorCatcher)
      );
  }

  public getReports(
    guildID: string,
    memberID: string = null,
    offset: number = 0,
    limit: number = 100
  ): Observable<Report[]> {
    const uri = memberID
      ? this.rcGuildMemberReports(guildID, memberID)
      : this.rcGuildReports(guildID);

    const opts = this.defopts({
      params: new HttpParams()
        .set('sortBy', 'created')
        .set('offset', offset.toString())
        .set('limit', limit.toString()),
    });

    return this.http.get<ListReponse<Report>>(uri, opts).pipe(
      map((lr) => lr.data),
      catchError(this.errorCatcher)
    );
  }

  public getReport(reportID: string): Observable<Report> {
    return this.http
      .get<Report>(this.rcReports(reportID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getReportsCount(
    guildID: string,
    memberID: string = null
  ): Observable<number> {
    const uri = memberID
      ? this.rcGuildMemberReportsCount(guildID, memberID)
      : this.rcGuildReportsCount(guildID);

    return this.http.get<Count>(uri, this.defopts()).pipe(
      map((c) => c.count),
      catchError(this.errorCatcher)
    );
  }

  public postReportRevoke(reportID: string, reason: string): Observable<any> {
    return this.http
      .post(this.rcReports(reportID, 'revoke'), { reason }, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettings(guildID: string): Observable<GuildSettings> {
    return this.http
      .get<GuildSettings>(this.rcGuildSettings(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettings(
    guildID: string,
    settings: GuildSettings
  ): Observable<any> {
    return this.http
      .post(this.rcGuildSettings(guildID), settings, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildPermissions(
    guildID: string
  ): Observable<Map<string, string[]>> {
    return this.http
      .get<Map<string, string[]>>(
        this.rcGuildPermissions(guildID),
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildPermissions(
    guildID: string,
    update: PermissionsUpdate
  ): Observable<any> {
    return this.http
      .post(this.rcGuildPermissions(guildID), update, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postReport(
    guildID: string,
    memberID: string,
    rep: ReportRequest
  ): Observable<Report> {
    return this.http
      .post<Report>(
        this.rcGuildMemberReports(guildID, memberID),
        rep,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postKick(
    guildID: string,
    memberID: string,
    rep: ReasonRequest
  ): Observable<Report> {
    return this.http
      .post<Report>(
        this.rcGuildMemberKick(guildID, memberID),
        rep,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postBan(
    guildID: string,
    memberID: string,
    rep: ReasonRequest,
    anonymous: boolean = false
  ): Observable<Report> {
    const opts = this.defopts({
      params: new HttpParams().set('anonymous', anonymous ? '1' : '0'),
    });
    return this.http
      .post<Report>(this.rcGuildMemberBan(guildID, memberID), rep, opts)
      .pipe(catchError(this.errorCatcher));
  }

  public postMute(
    guildID: string,
    memberID: string,
    rep: ReasonRequest
  ): Observable<Report> {
    return this.http
      .post<Report>(
        this.rcGuildMemberMute(guildID, memberID),
        rep,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postUnmute(
    guildID: string,
    memberID: string,
    rep: ReasonRequest
  ): Observable<any> {
    return this.http
      .post<any>(
        this.rcGuildMemberUnmute(guildID, memberID),
        rep,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildBackupToggle(
    guildID: string,
    enabled: boolean
  ): Observable<any> {
    return this.http
      .post(this.rcGuildBackups(guildID, 'toggle'), { enabled }, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildBackups(
    guildID: string
  ): Observable<ListReponse<GuildBackup>> {
    return this.http
      .get<ListReponse<GuildBackup>>(
        this.rcGuildBackups(guildID),
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildBackupsDownload(
    guildID: string,
    backupID: string
  ): Observable<AccessTokenModel> {
    return this.http
      .post(
        `${this.rcGuildBackups(guildID)}/${backupID}/download`,
        null,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildInviteBlock(
    guildID: string,
    enabled: boolean
  ): Observable<any> {
    return this.http
      .post(this.rcGuildInviteBlock(guildID), { enabled }, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildScoreboard(
    guildID: string,
    limit: number = 25
  ): Observable<ListReponse<GuildScoreboardEntry>> {
    return this.http
      .get<Presence>(
        `${this.rcGuildScoreboard(guildID)}?limit=${limit}`,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildStarboard(
    guildID: string,
    sort: string = 'latest',
    limit: number = 20,
    offset: number = 0
  ): Observable<ListReponse<GuildStarboardEntry>> {
    return this.http
      .get<Presence>(
        `${this.rcGuildStarboard(
          guildID
        )}?limit=${limit}&offset=${offset}&sort=${sort}`,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getPresence(): Observable<Presence> {
    return this.http
      .get<Presence>(this.rcSetting('presence'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postPresence(p: Presence): Observable<Presence> {
    return this.http
      .post<Presence>(this.rcSetting('presence'), p, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getInviteSettings(): Observable<InviteSettingsResponse> {
    return this.http
      .get<InviteSettingsResponse>(
        this.rcSetting('noguildinvite'),
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postInviteSettings(s: InviteSettingsRequest): Observable<any> {
    return this.http
      .post(this.rcSetting('noguildinvite'), s, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getSystemInfo(): Observable<SystemInfo> {
    return this.http
      .get(this.rcAPI('sysinfo'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getAPIToken(ignoreError: boolean = false): Observable<APIToken> {
    let req = this.http.get<APIToken>(this.rcAPI('token'), this.defopts());
    if (!ignoreError) {
      req = req.pipe(catchError(this.errorCatcher));
    }
    return req;
  }

  public postAPIToken(): Observable<APIToken> {
    return this.http
      .post(this.rcAPI('token'), null, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public deleteAPIToken(): Observable<any> {
    return this.http
      .delete(this.rcAPI('token'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getCommandInfos(): Observable<ListReponse<CommandInfo>> {
    return this.http
      .get(this.rcUtil('slashcommands'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsKarma(guildID: string): Observable<KarmaSettings> {
    return this.http
      .get(this.rcGuildSettingsKarma(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsKarma(
    guildID: string,
    settings: KarmaSettings
  ): Observable<any> {
    return this.http
      .post(this.rcGuildSettingsKarma(guildID), settings, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsKarmaBlocklist(
    guildID: string
  ): Observable<ListReponse<Member>> {
    return this.http
      .get(this.rcGuildSettingsKarmaBlocklist(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public putGuildSettingsKarmaBlocklist(
    guildID: string,
    userIdent: string
  ): Observable<any> {
    return this.http
      .put(
        this.rcGuildSettingsKarmaBlocklist(guildID, userIdent),
        null,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public deleteGuildSettingsKarmaBlocklist(
    guildID: string,
    userID: string
  ): Observable<any> {
    return this.http
      .delete(
        this.rcGuildSettingsKarmaBlocklist(guildID, userID),
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsKarmaRules(
    guildID: string
  ): Observable<ListReponse<KarmaRule>> {
    return this.http
      .get(this.rcGuildSettingsKarmaRules(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public createGuildSettingsKarmaRules(rule: KarmaRule): Observable<KarmaRule> {
    return this.http
      .post(this.rcGuildSettingsKarmaRules(rule.guildid), rule, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public updateGuildSettingsKarmaRules(rule: KarmaRule): Observable<KarmaRule> {
    return this.http
      .post(
        this.rcGuildSettingsKarmaRules(rule.guildid, rule.id),
        rule,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public deleteGuildSettingsKarmaRules(rule: KarmaRule): Observable<KarmaRule> {
    return this.http
      .delete(
        this.rcGuildSettingsKarmaRules(rule.guildid, rule.id),
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsAntiraid(
    guildID: string
  ): Observable<AntiraidSettings> {
    return this.http
      .get(this.rcGuildSettingsAntiraid(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsAntiraid(
    guildID: string,
    settings: AntiraidSettings
  ): Observable<any> {
    return this.http
      .post(this.rcGuildSettingsAntiraid(guildID), settings, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildAntiraidJoinlog(
    guildID: string
  ): Observable<ListReponse<JoinlogEntry>> {
    return this.http
      .get(this.rcGuildAntiraidJoinlog(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public deleteGuildAntiraidJoinlog(guildID: string): Observable<any> {
    return this.http
      .delete(this.rcGuildAntiraidJoinlog(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildAntiraidAction(
    guildID: string,
    action: AntiraidAction
  ): Observable<any> {
    return this.http
      .post(this.rcGuildSettingsAntiraidAction(guildID), action, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getLandingPageInfo(): Observable<any> {
    return this.http
      .get(this.rcUtil('landingpageinfo'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildUnbanrequests(
    guildId: string
  ): Observable<ListReponse<UnbanRequest>> {
    return this.http
      .get(this.rcGuildUnbanRequest(guildId), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildMemberUnbanrequests(
    guildId: string,
    memberId: string
  ): Observable<ListReponse<UnbanRequest>> {
    return this.http
      .get(this.rcGuildMemberUnbanRequest(guildId, memberId), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildUnbanrequest(
    guildId: string,
    id: string
  ): Observable<ListReponse<UnbanRequest>> {
    return this.http
      .get(this.rcGuildUnbanRequest(guildId, id), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildUnbanrequestCount(
    guildId: string,
    stateFilter: number = -1
  ): Observable<Count> {
    const opts = this.defopts({
      params: new HttpParams().set('state', stateFilter.toString()),
    });
    return this.http
      .get(this.rcGuildUnbanRequest(guildId, 'count'), opts)
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildUnbanrequest(
    guildId: string,
    request: UnbanRequest
  ): Observable<ListReponse<UnbanRequest>> {
    return this.http
      .post(
        this.rcGuildUnbanRequest(guildId, request.id),
        request,
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getUnbanrequestBannedguilds(): Observable<ListReponse<Guild>> {
    return this.http
      .get(this.rcUnbanRequests('bannedguilds'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getUnbanrequests(): Observable<ListReponse<UnbanRequest>> {
    return this.http
      .get(this.rcUnbanRequests(), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postUnbanrequests(request: UnbanRequest): Observable<UnbanRequest> {
    return this.http
      .post(this.rcUnbanRequests(), request, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getUserSettingsOTA(): Observable<UserSettingsOTA> {
    return this.http
      .get(this.rcUserSettings('ota'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postUserSettingsOTA(
    ota: UserSettingsOTA
  ): Observable<UserSettingsOTA> {
    return this.http
      .post(this.rcUserSettings('ota'), ota, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsLogs(
    guildID: string,
    limit = 50,
    offset = 0,
    severity = -1
  ): Observable<ListReponse<GuildLogEntry>> {
    const opts = this.defopts({
      params: new HttpParams()
        .set('severity', severity.toString())
        .set('offset', offset.toString())
        .set('limit', limit.toString()),
    });

    return this.http
      .get(this.rcGuildSettingsLogs(guildID), opts)
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsLogsCount(
    guildID: string,
    severity = -1
  ): Observable<Count> {
    const opts = this.defopts({
      params: new HttpParams().set('severity', severity.toString()),
    });

    return this.http
      .get(this.rcGuildSettingsLogs(guildID, 'count'), opts)
      .pipe(catchError(this.errorCatcher));
  }

  public deleteGuildSettingsLogs(
    guildID: string,
    id?: string
  ): Observable<any> {
    return this.http
      .delete(this.rcGuildSettingsLogs(guildID, id), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsLogsState(guildID: string): Observable<State> {
    return this.http
      .get(this.rcGuildSettingsLogs(guildID, 'state'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsLogsState(
    guildID: string,
    state: boolean
  ): Observable<State> {
    return this.http
      .post(
        this.rcGuildSettingsLogs(guildID, 'state'),
        { state },
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsFlushGuildData(
    guildID: string,
    validation: string,
    leave_after: boolean
  ): Observable<State> {
    return this.http
      .post(
        this.rcGuildSettingsFlushData(guildID),
        { validation, leave_after },
        this.defopts()
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getSearch(query: string, limit = 50): Observable<SearchResult> {
    const opts = this.defopts({
      params: new HttpParams()
        .set('query', query)
        .set('limit', limit.toString()),
    });
    return this.http
      .get(this.rcSearch(), opts)
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsApi(guildID: string): Observable<GuildSettingsApi> {
    return this.http
      .get(this.rcGuildSettingsApi(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsApi(
    guildID: string,
    settings: GuildSettingsApi
  ): Observable<GuildSettingsApi> {
    return this.http
      .post(this.rcGuildSettingsApi(guildID), settings, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getChannels(
    guildID: string
  ): Observable<ListReponse<ChannelWithPermissions>> {
    return this.http
      .get(this.rcChannels(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postChannels(
    guildID: string,
    chanID: string,
    embed: MessageEmbed
  ): Observable<any> {
    return this.http
      .post(this.rcChannels(guildID, chanID), embed, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postChannelsMessage(
    guildID: string,
    chanID: string,
    msgID: string,
    embed: MessageEmbed
  ): Observable<any> {
    return this.http
      .post(this.rcChannels(guildID, chanID, msgID), embed, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getUser(userID: string): Observable<User> {
    return this.http
      .get(this.rcUsers(userID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getVerificationSiteKey(): Observable<VerificationSiteKey> {
    return this.http
      .get(this.rcVerification('sitekey'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postVerifyUser(token: string): Observable<any> {
    return this.http
      .post(this.rcVerification('verify'), { token }, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsVerification(
    guildID: string
  ): Observable<GuildSettingsVerification> {
    return this.http
      .get(this.rcGuildSettingsVerification(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsVerification(
    guildID: string,
    settings: GuildSettingsVerification
  ): Observable<any> {
    return this.http
      .post(this.rcGuildSettingsVerification(guildID), settings, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettingsCodeExec(
    guildID: string
  ): Observable<CodeExecSettings> {
    return this.http
      .get(this.rcGuildSettingsCodeExec(guildID), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettingsCodeExec(
    guildID: string,
    settings: CodeExecSettings
  ): Observable<any> {
    return this.http
      .post(this.rcGuildSettingsCodeExec(guildID), settings, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postUserSettingsFlush(): Observable<object> {
    return this.http
      .post(this.rcUserSettings('flush'), {}, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }
}
