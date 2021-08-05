/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { GuildSettingsApi } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-ga-api',
  templateUrl: './ga-api.component.html',
  styleUrls: ['./ga-api.component.scss'],
})
export class GuildAdminApiComponent implements OnInit {
  state: GuildSettingsApi;

  private guildID: string;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;
      this.state = await this.api.getGuildSettingsApi(this.guildID).toPromise();
    });
  }

  async saveSettings() {
    this.state = await this.api
      .postGuildSettingsApi(this.guildID, this.state)
      .toPromise();
    this.toasts.push('Settings were saved.', '', 'success');
  }

  async resetToken() {
    this.state = await this.api
      .postGuildSettingsApi(this.guildID, {
        enabled: this.state.enabled,
        allowed_origins: this.state.allowed_origins,
        reset_token: true,
      } as GuildSettingsApi)
      .toPromise();
    this.toasts.push('Token was removed.', '', 'success');
  }

  get apiUrl(): string {
    return `${window.location.protocol}//${
      window.location.host
    }/api/public/guilds/${this.guildID}${
      this.state.protected ? '?token={token}' : ''
    }`;
  }
}
