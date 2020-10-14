/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { KarmaSettings } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-ga-karma',
  templateUrl: './ga-karma.component.html',
  styleUrls: ['./ga-karma.component.sass'],
})
export class GuildAdminKarmaComponent implements OnInit {
  public karmaSettings: KarmaSettings;
  private guildID: string;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.karmaSettings = await this.api
        .getGuildSettingsKarma(this.guildID)
        .toPromise();
    });
  }

  public async onSave() {
    try {
      await this.api
        .postGuildSettingsKarma(this.guildID, this.karmaSettings)
        .toPromise();
      this.toasts.push(
        'Karma settings saved.',
        'Settings saved',
        'green',
        4000
      );
    } catch {}
  }

  public onIncChange(event: any) {
    this.karmaSettings.emotes_increase = event.target.value
      .split(',')
      .map((v: string) => v.trim());
  }

  public onDecChange(event: any) {
    this.karmaSettings.emotes_decrease = event.target.value
      .split(',')
      .map((v: string) => v.trim());
  }
}
