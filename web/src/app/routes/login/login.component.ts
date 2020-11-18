/** @format */

import { Component } from '@angular/core';
import { LandingPageInfo } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.sass'],
})
export class LoginComponent {
  public info: LandingPageInfo;

  constructor(private api: APIService) {
    this.api.getLandingPageInfo().subscribe((info) => {
      this.info = info;
    });
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
