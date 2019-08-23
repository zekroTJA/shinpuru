/** @format */

import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { catchError, map, share } from 'rxjs/operators';
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
} from './api.models';
import { environment } from 'src/environments/environment';
import { ToastService } from '../components/toast/toast.service';
import { CacheBucket } from './api.cache';
import { isObject } from 'util';

/** @format */

@Injectable({
  providedIn: 'root',
})
export class APIService {
  private rootURL = '';

  private readonly cacheMembers = new CacheBucket<string, Member>(
    10 * 60 * 1000
  );
  private readonly cacheUsers = new CacheBucket<string, User>(10 * 60 * 1000);
  private readonly cacheGuilds = new CacheBucket<string, Guild>(30 * 1000);

  private readonly defopts = (obj?: object) => {
    const defopts = {
      withCredentials: true,
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

  private readonly rcReports = (reportID: string) =>
    `${this.rcAPI('reports')}/${reportID}`;

  private readonly rcGuildSettings = (guildID: string) =>
    `${this.rcGuilds(guildID)}/settings`;

  private readonly rcGuildPermissions = (guildID: string) =>
    `${this.rcGuilds(guildID)}/permissions`;

  private readonly rcGuildMemberKick = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/kick`;

  private readonly rcGuildMemberBan = (guildID: string, memberID: string) =>
    `${this.rcGuildMembers(guildID, memberID)}/ban`;

  private readonly rcSetting = (rc: string = '') =>
    `${this.rcAPI('settings')}${rc ? '/' + rc : ''}`;

  private readonly errorCatcher = (err) => {
    console.error(err);
    this.toasts.push(err.message, 'Request Error', 'error', 10000);
    return of(null);
  };

  constructor(private http: HttpClient, private toasts: ToastService) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }

  public logout(): Observable<any> {
    return this.http
      .post<any>(this.rcAPI('logout'), this.defopts())
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

    return this.http.get<Guild>(this.rcGuilds(id), this.defopts()).pipe(
      this.cacheGuilds.putFromPipe(id),
      catchError(this.errorCatcher)
    );
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
    const m = this.cacheMembers.get(memberID);
    if (m) {
      return of(m);
    }

    return this.http
      .get<Member>(this.rcGuildMembers(guildID, memberID), this.defopts())
      .pipe(
        this.cacheMembers.putFromPipe(memberID),
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
    return this.http
      .get<ListReponse<string>>(
        this.rcGuildMembersPermissionsAllowed(guildID, memberID),
        this.defopts()
      )
      .pipe(
        map((l) => l.data),
        catchError(this.errorCatcher)
      );
  }

  public getReports(
    guildID: string,
    memberID: string = null,
    offset: number = 0,
    limit: number = 0
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
    rep: ReasonRequest
  ): Observable<Report> {
    return this.http
      .post<Report>(
        this.rcGuildMemberBan(guildID, memberID),
        rep,
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
}
