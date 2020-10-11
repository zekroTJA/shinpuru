/** @format */

import { Component, OnInit } from '@angular/core';

interface Route {
  route: string;
  icon: string;
  displayname: string;
}

@Component({
  selector: 'app-guild-admin-navbar',
  templateUrl: './guild-admin-navbar.component.html',
  styleUrls: ['./guild-admin-navbar.component.sass'],
})
export class GuildAdminNavbarComponent implements OnInit {
  public routes: Route[] = [
    {
      route: 'karma',
      icon: 'karma.svg',
      displayname: 'Karma',
    },
  ];

  public currentPath: string;

  constructor() {}

  ngOnInit(): void {
    const path = window.location.pathname.split('/');
    this.currentPath = path[path.length - 1];
  }
}
