/** @format */

import { Component } from '@angular/core';
import { APIService } from '../api/api.service';
import { User } from '../api/api.models';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.sass'],
})
export class HeaderComponent {
  public selfUser: User;

  constructor(private api: APIService) {
    this.api.getSelfUser().subscribe((user) => {
      this.selfUser = user;
    });
  }
}
