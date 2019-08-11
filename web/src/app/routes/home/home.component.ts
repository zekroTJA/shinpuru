/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { Guild, InviteSettingsResponse } from 'src/app/api/api.models';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.sass'],
})
export class HomeComponent {
  public guilds: Guild[] = [];
  public inviteSettings: InviteSettingsResponse;

  constructor(private api: APIService, private spinner: SpinnerService) {
    this.api.getGuilds().subscribe((guilds) => {
      this.guilds = guilds;

      if (this.guilds.length < 1) {
        this.api.getInviteSettings().subscribe((inviteSettings) => {
          this.inviteSettings = inviteSettings;
          this.spinner.stop('spinner-load-guilds');
        });
      } else {
        this.spinner.stop('spinner-load-guilds');
      }
    });
  }
}
