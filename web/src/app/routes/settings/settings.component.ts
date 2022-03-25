/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import {
  Presence,
  InviteSettingsResponse,
  InviteSettingsRequest,
  Guild,
} from 'src/app/api/api.models';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
})
export class SettingsComponent {
  public presence: Presence;
  public inviteSettings: InviteSettingsResponse;
  public inviteSettingsFields: InviteSettingsRequest =
    {} as InviteSettingsRequest;
  public guilds: Guild[];

  constructor(private api: APIService, private toasts: ToastService) {
    this.api.getPresence().subscribe((presence) => {
      this.presence = presence;
    });

    this.api.getGuilds().subscribe((guilds) => {
      this.guilds = guilds;
    });

    this.api.getInviteSettings().subscribe((inviteSettings) => {
      if (inviteSettings) {
        const invUrlSplit = inviteSettings.invite_url
          ? inviteSettings.invite_url.split('/')
          : [];

        this.inviteSettings = inviteSettings;
        this.inviteSettingsFields.guild_id = inviteSettings.guild
          ? inviteSettings.guild.id
          : '';
        this.inviteSettingsFields.invite_code =
          invUrlSplit[invUrlSplit.length - 1];
        this.inviteSettingsFields.message = inviteSettings.message;
      }
    });
  }

  public updatePresence() {
    this.api.postPresence(this.presence).subscribe((presence) => {
      if (presence) {
        this.presence = presence;
        this.toasts.push(
          'Updated bot presence.',
          'Updated',
          'success',
          6000,
          true
        );
      }
    });
  }

  public updateInvite() {
    this.api.postInviteSettings(this.inviteSettingsFields).subscribe((res) => {
      if (res.code === 200) {
        this.toasts.push(
          'Updated guild invite.',
          'Updated',
          'success',
          6000,
          true
        );
      }
    });
  }

  public resetInvite() {
    this.inviteSettingsFields = {
      guild_id: '',
      invite_code: '',
      message: '',
    };
    this.api.postInviteSettings(this.inviteSettingsFields).subscribe((res) => {
      if (res.code === 200) {
        this.toasts.push(
          'Reset guild invite.',
          'Updated',
          'success',
          6000,
          true
        );
      }
    });
  }
}
