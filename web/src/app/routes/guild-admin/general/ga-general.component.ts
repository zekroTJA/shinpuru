/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  AntiraidSettings,
  Channel,
  Guild,
  GuildSettings,
  JoinlogEntry,
  Role,
} from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';

@Component({
  selector: 'app-ga-general',
  templateUrl: './ga-general.component.html',
  styleUrls: ['./ga-general.component.scss'],
})
export class GuildAdminGeneralComponent implements OnInit {
  public antiraidSettings: AntiraidSettings;
  public joinlog: JoinlogEntry[] = [];
  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);
  public settings: GuildSettings;
  public updatedSettings = {} as GuildSettings;
  public autoRoles: Role[] = [];

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
    // return;
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
      this.autoRoles = this.settings.autoroles.map(
        (rid) =>
          this.guild.roles.find((r) => r.id === rid) ??
          ({
            id: rid,
          } as Role)
      );
    });
  }

  public guildSettingsContains(str: string): boolean {
    return this.guildSettingsAllowed.includes(str);
  }

  public guildSettingsContainsAny(str: string[]): boolean {
    return !!str.find((s) => this.guildSettingsContains(s));
  }

  public channelsByType(a: Channel[], type: number): Channel[] {
    return a?.filter((c) => c.type === type) ?? [];
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
    this.updatedSettings.autoroles = this.autoRoles.map((r) => r.id);
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

  roleInputFilter(v: Role, inpt: string): boolean {
    if (v.id === inpt) {
      return true;
    }

    return (
      v.name !== '@everyone' &&
      v.name.toLowerCase().includes(inpt.toLowerCase())
    );
  }

  roleNameFormatter(r: Role): string {
    return r.name ?? `<deleted Role> ${r.id}`;
  }

  roleInvalidFilter(r: Role): boolean {
    return !r.name;
  }
}
