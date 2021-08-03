/** @format */

import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
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
  isSearch = false;

  constructor(public toasts: ToastService, private router: Router) {}

  ngOnInit() {
    const nlr = LocalStorageUtil.get<NextLoginRedirect>('NEXT_LOGIN_REDIRECT');
    if (nlr && nlr.deadline >= Date.now()) {
      LocalStorageUtil.remove('NEXT_LOGIN_REDIRECT');
      window.location.replace(nlr.destination);
    }

    window.onkeydown = (e: KeyboardEvent) => {
      if (e.ctrlKey && e.key === 'f') {
        e.preventDefault();
        this.isSearch = true;
      }

      if (e.key === 'Escape' && this.isSearch) {
        e.preventDefault;
        this.isSearch = false;
      }
    };
  }

  onSearchNavigate(route: string[]) {
    this.isSearch = false;
    this.router.navigate(route);
  }

  onSearchBgClick(e: MouseEvent) {
    if ((e.target as HTMLElement).id !== 'search-bar-container') return;
    e.preventDefault();
    this.isSearch = false;
  }
}
