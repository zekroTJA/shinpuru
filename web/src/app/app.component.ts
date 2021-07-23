/** @format */

import { Component, OnInit } from '@angular/core';
import { LoadingBarService } from './components/header/loadingbar.service';
import { ToastService } from './components/toast/toast.service';
import LocalStorageUtil from './utils/localstorage';
import { NextLoginRedirect } from './utils/objects';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  title = 'shinpuru Web Interface';

  constructor(public toasts: ToastService) {}

  ngOnInit() {
    const nlr = LocalStorageUtil.get<NextLoginRedirect>('NEXT_LOGIN_REDIRECT');
    if (nlr && nlr.deadline >= Date.now()) {
      LocalStorageUtil.remove('NEXT_LOGIN_REDIRECT');
      window.location.replace(nlr.destination);
    }
  }
}
