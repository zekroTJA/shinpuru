/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { ActivatedRoute } from '@angular/router';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import { Member, Guild, Role, Report } from 'src/app/api/api.models';
import dateFormat from 'dateformat';

@Component({
  selector: 'app-member-route',
  templateUrl: './member.component.html',
  styleUrls: ['./member.component.sass'],
})
export class MemberRouteComponent {
  public member: Member;
  public guild: Guild;
  public perm: string[];

  public reports: Report[] = [];

  public dateFormat = dateFormat;

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

    this.api.getPermissions(guildID, memberID).subscribe((perm) => {
      this.perm = perm;
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

  public getPerms(allowed: boolean): string[] {
    return this.perm.filter((p) => p.startsWith(allowed ? '+' : '-'));
  }
}
