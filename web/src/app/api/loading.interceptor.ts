import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { finalize } from 'rxjs/operators';
import { LoadingBarService } from '../components/header/loadingbar.service';

@Injectable()
export default class LoadingInterceptor implements HttpInterceptor {
  constructor(private loadingBar: LoadingBarService) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    this.loadingBar.init();
    return next.handle(req).pipe(finalize(() => this.loadingBar.finished()));
  }
}
