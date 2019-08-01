/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import { ActivatedRoute } from '@angular/router';
import { Guild, Role, Member } from 'src/app/api/api.models';

@Component({
  selector: 'app-guild',
  templateUrl: './guild.component.html',
  styleUrls: ['./guild.component.sass'],
})
export class GuildComponent {
  public guild: Guild;
  public permLvl: number;
  public members: Member[];

  constructor(
    public api: APIService,
    public spinner: SpinnerService,
    private route: ActivatedRoute
  ) {
    const guildID = this.route.snapshot.paramMap.get('id');
    this.api.getGuild(guildID).subscribe((guild) => {
      this.guild = guild;
      this.members = this.guild.members.filter(
        (m) => m.user.id !== this.guild.self_member.user.id
      );
      this.spinner.stop('spinner-load-guild');
    });
    this.api.getPermissionLvl(guildID).subscribe((lvl) => {
      this.permLvl = lvl;
    });
  }

  public get userRoles(): Role[] {
    const userRoleIDs = this.guild.self_member.roles;
    return this.guild.roles
      .filter((r) => userRoleIDs.includes(r.id))
      .sort((a, b) => b.position - a.position);
  }

  public permLvlColor(lvl: number): string {
    if (lvl < 1) {
      return '#424242';
    } else if (lvl < 3) {
      return '#0288D1';
    } else if (lvl < 5) {
      return '#689F38';
    } else if (lvl < 7) {
      return '#FFA000';
    } else if (lvl < 9) {
      return '#E64A19';
    } else if (lvl < 11) {
      return '#d32f2f';
    }

    return '#F50057';
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
}
