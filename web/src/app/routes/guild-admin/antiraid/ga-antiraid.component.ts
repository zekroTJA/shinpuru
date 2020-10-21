/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AntiraidSettings, KarmaSettings } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-ga-antiraid',
  templateUrl: './ga-antiraid.component.html',
  styleUrls: ['./ga-antiraid.component.sass'],
})
export class GuildAdminAntiraidComponent implements OnInit {
  public antiraidSettings: AntiraidSettings;
  private guildID: string;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.antiraidSettings = await this.api
        .getGuildSettingsAntiraid(this.guildID)
        .toPromise();
    });
  }

  public async onSave() {
    try {
      await this.api
        .postGuildSettingsAntiraid(this.guildID, this.antiraidSettings)
        .toPromise();
      this.toasts.push(
        'Antiraid settings saved.',
        'Settings saved',
        'green',
        4000
      );
    } catch {}
  }
}
