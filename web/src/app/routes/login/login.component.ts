/** @format */

import { Component } from '@angular/core';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.sass'],
})
export class LoginComponent {
  public scrollTo(id: string) {
    const e = document.querySelector('#' + id);
    if (e) {
      e.scrollIntoView({
        behavior: 'smooth',
      });
    }
  }
}
