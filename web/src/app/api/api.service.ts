/** @format */

import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { User } from './api.models';
import { environment } from 'src/environments/environment';

/** @format */

@Injectable({
  providedIn: 'root',
})
export class APIService {
  private rootURL = '';
  private defopts = {
    withCredentials: true,
  };

  constructor(private http: HttpClient) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }

  public getSelfUser(): Observable<User> {
    return this.http.get<User>(this.rootURL + '/api/me', this.defopts);
  }
}
