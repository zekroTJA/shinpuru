/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import { ActivatedRoute } from '@angular/router';
import {
  Guild,
  Role,
  Member,
  Report,
  GuildSettings,
  Channel,
} from 'src/app/api/api.models';
import { ToastService } from 'src/app/components/toast/toast.service';
import { toHexClr } from '../../utils/utils';

interface Perms {
  id: string;
  role: Role;
  perms: string[];
}

@Component({
  selector: 'app-guild',
  templateUrl: './guild.component.html',
  styleUrls: ['./guild.component.sass'],
})
export class GuildComponent {
  public guild: Guild;
  public members: Member[];
  public reports: Report[];
  public allowed: string[];
  public settings: GuildSettings;
  public updatedSettings: GuildSettings = {} as GuildSettings;

  public guildSettingsAllowed: string[] = [];

  public addPermissionPerm: string;
  public addPermissionRoles: string[] = [''];

  public guildToggle = false;
  public modlogToggle = false;
  public guildSettingsToggle = false;
  public permissionsToggle = true;

  public toHexClr = toHexClr;

  constructor(
    public api: APIService,
    public spinner: SpinnerService,
    private route: ActivatedRoute,
    private toasts: ToastService
  ) {
    const guildID = this.route.snapshot.paramMap.get('id');
    this.api.getGuild(guildID).subscribe((guild) => {
      this.guild = guild;
      this.members = this.guild.members.filter(
        (m) => m.user.id !== this.guild.self_member.user.id
      );
      this.api
        .getPermissionsAllowed(guildID, guild.self_member.user.id)
        .subscribe((allowed) => {
          this.allowed = allowed;
          this.guildSettingsAllowed = this.allowed.filter((a) =>
            a.startsWith('sp.guild.config')
          );
        });
      this.spinner.stop('spinner-load-guild');
    });

    this.api.getGuildSettings(guildID).subscribe((settings) => {
      this.settings = settings;
    });

    this.api.getReports(guildID).subscribe((reports) => {
      this.reports = reports;
    });
  }

  public get userRoles(): Role[] {
    const userRoleIDs = this.guild.self_member.roles;
    return this.guild.roles
      .filter((r) => userRoleIDs.includes(r.id))
      .sort((a, b) => b.position - a.position);
  }

  public searchInput(e: any) {
    const val = e.target.value.toLowerCase();

    if (val === '') {
      this.members = this.guild.members.filter(
        (m) => m.user.id !== this.guild.self_member.user.id
      );
    } else {
      this.members = this.guild.members.filter(
        (m) =>
          m.user.id !== this.guild.self_member.user.id &&
          ((m.nick && m.nick.toLowerCase().includes(val)) ||
            m.user.username.toLowerCase().includes(val) ||
            m.user.id.includes(val))
      );
    }
  }

  public guildSettingsContains(str: string): boolean {
    return this.guildSettingsAllowed.includes(str);
  }

  public saveGuildSettings() {
    console.log(this.updatedSettings);
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

  public getSelectedValue(e: any): string {
    const t = e.target;
    return t.options[t.selectedIndex].value
      .split(' ')
      .slice(1)
      .join(' ');
  }

  public channelsByType(a: Channel[], type: number): Channel[] {
    return a.filter((c) => c.type === type);
  }

  public fetchGuildPermissions() {
    this.api.getGuildPermissions(this.guild.id).subscribe((perms) => {
      this.settings.perms = perms;
    });
  }

  public getRoleByID(roleID: string): Role {
    return this.guild.roles.find((r) => r.id === roleID);
  }

  public objectAsArray(obj: any): Perms[] {
    if (!obj) {
      return [];
    }

    return Object.keys(obj).map<Perms>((k) => {
      return { id: k, role: this.getRoleByID(k), perms: obj[k] };
    });
  }

  public removePermission(p: Perms, perm: string) {
    const prefix = perm.startsWith('+') ? '-' : '+';
    perm = prefix + perm.substr(1);
    this.api
      .postGuildPermissions(this.guild.id, { role_ids: [p.id], perm })
      .subscribe(() => {
        this.fetchGuildPermissions();
      });
  }

  public inputAddPermissionRole(val: string, i: number) {
    console.log(val);
    const wasEmpty = this.addPermissionRoles[i] === '';
    this.addPermissionRoles[i] = val;
    if (i + 1 === this.addPermissionRoles.length && wasEmpty) {
      this.addPermissionRoles.push('');
    }
  }
}
