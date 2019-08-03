/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { ActivatedRoute } from '@angular/router';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import { Member, Guild, Role, Report } from 'src/app/api/api.models';
import dateFormat from 'dateformat';
import { permLvlColor } from 'src/app/utils/utils';

@Component({
  selector: 'app-member-route',
  templateUrl: './member.component.html',
  styleUrls: ['./member.component.sass'],
})
export class MemberRouteComponent {
  public member: Member;
  public guild: Guild;
  public permLvl: number;

  public reports: Report[] = [];

  public dateFormat = dateFormat;
  public permLvlColor = permLvlColor;

  constructor(
    public api: APIService,
    public spinner: SpinnerService,
    private route: ActivatedRoute
  ) {
    const guildID = this.route.snapshot.paramMap.get('guildid');
    const memberID = this.route.snapshot.paramMap.get('memberid');

    this.api.getGuildMember(guildID, memberID).subscribe((member) => {
      console.log(member);
      this.member = member;
      if (this.guild) {
        this.spinner.stop('spinner-load-member');
      }
    });

    this.api.getGuild(guildID).subscribe((guild) => {
      this.guild = guild;
      if (this.member) {
        this.spinner.stop('spinner-load-member');
      }
    });

    this.api.getPermissionLvl(guildID, memberID).subscribe((permLvl) => {
      this.permLvl = permLvl;
    });

    this.api.getReports(guildID, memberID).subscribe((reports) => {
      this.reports = reports || [];
    });

    // let remWatcher: NodeJS.Timer;
    // remWatcher = setInterval(() => {
    //   if (this.guild && this.member && this.reports) {
    //     this.spinner.stop('spinner-load-reports');
    //     clearInterval(remWatcher);
    //   }
    // }, 100);
  }

  public get memberRoles(): Role[] {
    const rls = this.guild.roles
      .filter((r) => this.member.roles.includes(r.id))
      .sort((a, b) => b.position - a.position);
    return rls;
  }
}
