/** @format */

import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
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
} from './api.models';
import { environment } from 'src/environments/environment';
import { ToastService } from '../components/toast/toast.service';

/** @format */

@Injectable({
  providedIn: 'root',
})
export class APIService {
  private rootURL = '';
  private defopts = {
    withCredentials: true,
  };

  private errorCatcher = (err) => {
    console.error(err);
    this.toasts.push(err.message, 'Request Error', 'error', 10000);
    return of(null);
  };

  constructor(private http: HttpClient, private toasts: ToastService) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }

  public logout(): Observable<any> {
    return this.http
      .post<any>(this.rootURL + '/api/logout', this.defopts)
      .pipe(catchError(this.errorCatcher));
  }

  public getSelfUser(): Observable<User> {
    return this.http.get<User>(this.rootURL + '/api/me', this.defopts).pipe(
      catchError((err) => {
        if (err.status !== 401) {
          return this.errorCatcher(err);
        }
      })
    );
  }

  public getGuilds(): Observable<Guild[]> {
    return this.http
      .get<ListReponse<Guild>>(this.rootURL + '/api/guilds', this.defopts)
      .pipe(
        map((lr) => {
          return lr.data;
        }),
        catchError(this.errorCatcher)
      );
  }

  public getGuild(id: string): Observable<Guild> {
    return this.http
      .get<Guild>(this.rootURL + '/api/guilds/' + id, this.defopts)
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildMember(
    guildID: string,
    memberID: string,
    ignoreError: boolean = false
  ): Observable<Member> {
    return this.http
      .get<Member>(
        this.rootURL + '/api/guilds/' + guildID + '/' + memberID,
        this.defopts
      )
      .pipe(catchError(ignoreError ? (err) => of(null) : this.errorCatcher));
  }

  public getPermissions(guildID: string, userID: string): Observable<string[]> {
    return this.http
      .get<PermissionResponse>(
        this.rootURL + '/api/guilds/' + guildID + '/' + userID + '/permissions',
        this.defopts
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
    userID: string
  ): Observable<string[]> {
    return this.http
      .get<ListReponse<string>>(
        this.rootURL +
          '/api/guilds/' +
          guildID +
          '/' +
          userID +
          '/permissions/allowed',
        this.defopts
      )
      .pipe(
        map((l) => l.data),
        catchError(this.errorCatcher)
      );
  }

  public getReports(
    guildID: string,
    memberID: string = null
  ): Observable<Report[]> {
    const uri = memberID
      ? this.rootURL + '/api/guilds/' + guildID + '/' + memberID + '/reports'
      : this.rootURL + '/api/guilds/' + guildID + '/reports';

    const opts = {
      withCredentials: this.defopts.withCredentials,
      params: new HttpParams().set('sortBy', 'created'),
    };

    return this.http.get<ListReponse<Report>>(uri, opts).pipe(
      map((lr) => lr.data),
      catchError(this.errorCatcher)
    );
  }

  public getReport(reportID: string): Observable<Report> {
    return this.http
      .get<Report>(this.rootURL + '/api/reports/' + reportID, this.defopts)
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildSettings(guildID: string): Observable<GuildSettings> {
    return this.http
      .get<GuildSettings>(
        this.rootURL + '/api/guilds/' + guildID + '/settings',
        this.defopts
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildSettings(
    guildID: string,
    settings: GuildSettings
  ): Observable<any> {
    return this.http
      .post(
        this.rootURL + '/api/guilds/' + guildID + '/settings',
        settings,
        this.defopts
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getGuildPermissions(
    guildID: string
  ): Observable<Map<string, string[]>> {
    return this.http
      .get<Map<string, string[]>>(
        this.rootURL + '/api/guilds/' + guildID + '/permissions',
        this.defopts
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postGuildPermissions(
    guildID: string,
    update: PermissionsUpdate
  ): Observable<any> {
    return this.http
      .post(
        this.rootURL + '/api/guilds/' + guildID + '/permissions',
        update,
        this.defopts
      )
      .pipe(catchError(this.errorCatcher));
  }

  public postReport(
    guildID: string,
    memberID: string,
    rep: ReportRequest
  ): Observable<Report> {
    return this.http
      .post<Report>(
        this.rootURL + '/api/guilds/' + guildID + '/' + memberID + '/reports',
        rep,
        this.defopts
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
        this.rootURL + '/api/guilds/' + guildID + '/' + memberID + '/kick',
        rep,
        this.defopts
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
        this.rootURL + '/api/guilds/' + guildID + '/' + memberID + '/ban',
        rep,
        this.defopts
      )
      .pipe(catchError(this.errorCatcher));
  }

  public getPresence(): Observable<Presence> {
    return this.http
      .get<Presence>(this.rootURL + '/api/settings/presence', this.defopts)
      .pipe(catchError(this.errorCatcher));
  }

  public postPresence(p: Presence): Observable<Presence> {
    return this.http
      .post<Presence>(this.rootURL + '/api/settings/presence', p, this.defopts)
      .pipe(catchError(this.errorCatcher));
  }
}
