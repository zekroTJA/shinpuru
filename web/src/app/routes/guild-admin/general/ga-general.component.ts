/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  AntiraidSettings,
  Channel,
  Guild,
  GuildSettings,
  JoinlogEntry,
} from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import dateFormat from 'dateformat';

@Component({
  selector: 'app-ga-general',
  templateUrl: './ga-general.component.html',
  styleUrls: ['./ga-general.component.sass'],
})
export class GuildAdminGeneralComponent implements OnInit {
  public antiraidSettings: AntiraidSettings;
  public joinlog: JoinlogEntry[] = [];
  public dateFormat = dateFormat;
  public settings: GuildSettings;
  public updatedSettings = {} as GuildSettings;

  private guildID: string;
  private guild: Guild;
  private allowed: string[];
  private guildSettingsAllowed: string[] = [];

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.guild = await this.api.getGuild(this.guildID).toPromise();
      this.allowed = await this.api
        .getPermissionsAllowed(this.guildID, this.guild.self_member.user.id)
        .toPromise();
      this.guildSettingsAllowed = this.allowed.filter((a) =>
        a.startsWith('sp.guild.')
      );
      this.settings = await this.api.getGuildSettings(this.guildID).toPromise();
    });
  }

  public guildSettingsContains(str: string): boolean {
    return this.guildSettingsAllowed.includes(str);
  }

  public guildSettingsContainsAny(str: string[]): boolean {
    return !!str.find((s) => this.guildSettingsContains(s));
  }

  public channelsByType(a: Channel[], type: number): Channel[] {
    return a.filter((c) => c.type === type);
  }

  public getSelectedValue(e: any): string {
    const t = e.target;
    const val = t.options[t.selectedIndex].value;
    if (val.match(/\d+:\s.+/g)) {
      return val.split(' ').slice(1).join(' ');
    }
    return val;
  }

  public saveGuildSettings() {
    this.api
      .postGuildSettings(this.guild.id, this.updatedSettings)
      .subscribe((res) => {
        if (res.code === 200) {
          this.toasts.push(
            'Guild settings updated.',
            'Guild Settings Update',
            'success',
            6000,
            true
          );
        }
      });
  }
}
