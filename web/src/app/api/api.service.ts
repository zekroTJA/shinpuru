/** @format */

import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { User, ListReponse, Guild } from './api.models';
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

  private errorChatcher = (err) => {
    console.error(err);
    this.toasts.push(err.message, 'Request Error', 'error', 10000);
    return of(null);
  };

  constructor(private http: HttpClient, private toasts: ToastService) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }

  public getSelfUser(): Observable<User> {
    return this.http.get<User>(this.rootURL + '/api/me', this.defopts).pipe(
      catchError((err) => {
        if (err.status !== 401) {
          return this.errorChatcher(err);
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
        catchError(this.errorChatcher)
      );
  }

  public getGuild(id: string): Observable<Guild> {
    return this.http
      .get<Guild>(this.rootURL + '/api/guilds/' + id, this.defopts)
      .pipe(catchError(this.errorChatcher));
  }
}
