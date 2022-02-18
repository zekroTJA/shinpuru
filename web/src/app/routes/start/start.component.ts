/** @format */

import { Component, OnInit } from '@angular/core';
import { LandingPageInfo } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import LocalStorageUtil from 'src/app/utils/localstorage';
import { NextLoginRedirect } from 'src/app/utils/objects';

@Component({
  selector: 'app-start',
  templateUrl: './start.component.html',
  styleUrls: ['./start.component.scss'],
})
export class StartComponent implements OnInit {
  public info: LandingPageInfo;

  constructor(private api: APIService) {
    this.api.getLandingPageInfo().subscribe((info) => {
      this.info = info;
    });
  }

  ngOnInit() {
    const params = new URLSearchParams(window.location.search);
    const redirect = params.get('redirect');
    console.log('login', redirect);

    if (redirect) {
      LocalStorageUtil.set('NEXT_LOGIN_REDIRECT', {
        destination: redirect,
        deadline: Date.now() + 5 * 60 * 1000,
      } as NextLoginRedirect);
    }
  }

  public scrollTo(id: string) {
    const e = document.querySelector('#' + id);
    if (e) {
      e.scrollIntoView({
        behavior: 'smooth',
      });
    }
  }

  public get currentYear(): number {
    return new Date().getFullYear();
  }
}
