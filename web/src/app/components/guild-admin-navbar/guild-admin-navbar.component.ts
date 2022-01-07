/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { APIService } from 'src/app/api/api.service';

interface Route {
  route: string;
  icon: string;
  displayname: string;
  perm?: string;
  permAny?: string[];
}

const ROUTES: Route[] = [
  {
    route: 'general',
    icon: 'settings-ol.svg',
    displayname: 'General',
    permAny: [
      'sp.guild.config.autorole',
      'sp.guild.config.prefix',
      'sp.guild.config.modlog',
      'sp.guild.config.voicelog',
      'sp.guild.config.announcements',
    ],
  },
  {
    route: 'antiraid',
    icon: 'antiraid.svg',
    displayname: 'Antiraid',
    perm: 'sp.guild.config.antiraid',
  },
  {
    route: 'verification',
    icon: 'verification.svg',
    displayname: 'Verification',
    perm: 'sp.guild.config.verification',
  },
  {
    route: 'karma',
    icon: 'karma.svg',
    displayname: 'Karma',
    perm: 'sp.guild.config.karma',
  },
  {
    route: 'logs',
    icon: 'logs.svg',
    displayname: 'Logs',
    perm: 'sp.guild.config.logs',
  },
  {
    route: 'data',
    icon: 'data.svg',
    displayname: 'Data',
    perm: 'sp.guild.admin.flushdata',
  },
  {
    route: 'api',
    icon: 'cloud.svg',
    displayname: 'API',
    perm: 'sp.guild.config.api',
  },
];

@Component({
  selector: 'app-guild-admin-navbar',
  templateUrl: './guild-admin-navbar.component.html',
  styleUrls: ['./guild-admin-navbar.component.scss'],
})
export class GuildAdminNavbarComponent implements OnInit {
  public routes: Route[] = [];

  public currentPath: string;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private api: APIService
  ) {
    route.params.subscribe(async (params) => {
      const guildID = params.guildid;
      const self = await this.api.getSelfUser().toPromise();
      const allowed = await api
        .getPermissionsAllowed(guildID, self.id)
        .toPromise();
      this.routes = ROUTES.filter((r) => {
        if (r.perm) return allowed.includes(r.perm);
        if (r.permAny) return !!r.permAny.find((p) => allowed.includes(p));
      });
    });
  }

  ngOnInit(): void {
    const path = window.location.pathname.split('/');
    this.currentPath = path[path.length - 1];
  }

  public navigate(r: Route) {
    const path = this.route.snapshot.url.map((u) => u.path);
    path[path.length - 1] = r.route;
    this.router.navigate(path);
  }
}
