/** @format */

import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { catchError, switchMap } from 'rxjs/operators';
import { APIService } from './api.service';

@Injectable()
export default class AuthInterceptor implements HttpInterceptor {
  constructor(private api: APIService) {}

  intercept(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    if (this.api.getStoredAccessToken()) {
      request = this.addToken(request, this.api.getStoredAccessToken().token);
    }

    return next.handle(request).pipe(
      catchError((error) => {
        if (
          error instanceof HttpErrorResponse &&
          error.status === 401 &&
          error.error?.error === 'invalid access token'
        )
          return this.handleInvalidAccessToken(request, next);

        return throwError(error);
      })
    );
  }

  private addToken(request: HttpRequest<any>, token: string) {
    return request.clone({
      setHeaders: {
        Authorization: `accessToken ${token}`,
      },
    });
  }

  private handleInvalidAccessToken(
    request: HttpRequest<any>,
    next: HttpHandler
  ) {
    return this.api.getAndSetAccessToken().pipe(
      switchMap((token) => {
        return next.handle(this.addToken(request, token?.token ?? ''));
      })
    );
  }
}
