/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { Guild, InviteSettingsResponse } from 'src/app/api/api.models';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import LocalStorageUtil from 'src/app/utils/localstorage';
import { NextLoginRedirect } from 'src/app/utils/objects';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  public guilds: Guild[] = null;
  public inviteSettings: InviteSettingsResponse;

  constructor(private api: APIService) {
    this.api.getGuilds().subscribe((guilds) => {
      this.guilds = guilds;
      if (this.guilds?.length < 1) {
        this.api.getInviteSettings().subscribe((inviteSettings) => {
          this.inviteSettings = inviteSettings;
        });
      }
    });
  }
}
