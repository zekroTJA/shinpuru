/** @format */

import { Component, TemplateRef, ViewChild } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { ActivatedRoute, Router } from '@angular/router';
import {
  Member,
  Guild,
  Role,
  Report,
  ReportRequest,
} from 'src/app/api/api.models';
import { format, formatDistance } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { ToastService } from 'src/app/components/toast/toast.service';
import { padNumber, rolePosDiff } from 'src/app/utils/utils';
import { MaxLengthValidator } from '@angular/forms';

@Component({
  selector: 'app-member-route',
  templateUrl: './member.component.html',
  styleUrls: ['./member.component.scss'],
})
export class MemberRouteComponent {
  public member: Member;
  public guild: Guild;
  public perm: string[];

  public reports: Report[];
  public permissionsAllowed: string[] = [];

  public roleDiff: number;

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

  public sinceFormat = (d: string | Date) =>
    formatDistance(new Date(d), new Date(), {
      addSuffix: true,
      includeSeconds: true,
    });

  @ViewChild('modalReport') private modalReport: TemplateRef<any>;
  @ViewChild('modalKick') private modalKick: TemplateRef<any>;
  @ViewChild('modalBan') private modalBan: TemplateRef<any>;
  @ViewChild('modalRevoke') private modalRevoke: TemplateRef<any>;
  @ViewChild('modalMute') private modalMute: TemplateRef<any>;
  @ViewChild('modalUnmute') private modalUnmute: TemplateRef<any>;

  public repModalType = 3;
  public repModalReason = '';
  public repModalAttachment = '';
  private repModalTimeout = '';

  public canRevoke = false;

  constructor(
    public modal: NgbModal,
    private api: APIService,
    private toasts: ToastService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    // return;
    const guildID = this.route.snapshot.paramMap.get('guildid');
    const memberID = this.route.snapshot.paramMap.get('memberid');

    this.api.getGuildMember(guildID, memberID).subscribe((member) => {
      this.member = member;
      if (this.guild) {
        this.ready();
      }
    });

    this.api.getGuild(guildID).subscribe((guild) => {
      this.guild = guild;
      if (this.member) {
        this.ready();
      }

      this.api
        .getPermissionsAllowed(guildID, guild.self_member.user.id)
        .subscribe((perms) => {
          this.permissionsAllowed = perms;
          this.canRevoke = this.hasPermission('sp.guild.mod.report');
        });
    });

    this.api.getPermissions(guildID, memberID).subscribe((perm) => {
      this.perm = perm;
    });

    this.fetchReports(guildID, memberID);
  }

  private ready() {
    this.roleDiff = rolePosDiff(
      this.guild.roles,
      this.guild.self_member,
      this.member
    );
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

  public hasPermission(perm: string): boolean {
    return this.permissionsAllowed.includes(perm);
  }

  public async report() {
    this.openModal(this.modalReport)
      .then((res) => {
        if (res && this.checkReason()) {
          this.api
            .postReport(this.guild.id, this.member.user.id, {
              attachment: this.repModalAttachment,
              reason: this.repModalReason,
              type: this.repModalType,
            })
            .subscribe((resRep) => {
              if (resRep) {
                this.fetchReports(this.guild.id, this.member.user.id);
                this.toasts.push(
                  'Report created.',
                  'Executed',
                  'success',
                  5000,
                  true
                );
              }
            });
        }
      })
      .finally(() => this.clearReportModalModels());
  }

  public kick() {
    this.openModal(this.modalKick)
      .then((res) => {
        if (res && this.checkReason()) {
          this.api
            .postKick(this.guild.id, this.member.user.id, {
              attachment: this.repModalAttachment,
              reason: this.repModalReason,
            })
            .subscribe((resRep) => {
              if (resRep) {
                this.router
                  .navigate(['../'], { relativeTo: this.route })
                  .then(() => {
                    this.toasts.push(
                      'Member kicked.',
                      'Executed',
                      'success',
                      5000,
                      true
                    );
                  });
              }
            });
        }
      })
      .finally(() => this.clearReportModalModels());
  }

  public ban() {
    this.openModal(this.modalBan)
      .then((res) => {
        if (res && this.checkReason()) {
          this.api
            .postBan(this.guild.id, this.member.user.id, {
              attachment: this.repModalAttachment,
              reason: this.repModalReason,
              timeout: this.formattedRepTimeout,
            })
            .subscribe((resRep) => {
              if (resRep) {
                this.router
                  .navigate(['../'], { relativeTo: this.route })
                  .then(() => {
                    this.toasts.push(
                      'Member banned.',
                      'Executed',
                      'success',
                      5000,
                      true
                    );
                  });
              }
            });
        }
      })
      .finally(() => this.clearReportModalModels());
  }

  public revokeReport(report: Report) {
    this.openModal(this.modalRevoke)
      .then((res) => {
        if (res && this.checkReason()) {
          this.api
            .postReportRevoke(report.id, this.repModalReason)
            .subscribe((revRes) => {
              if (revRes) {
                const i = this.reports.indexOf(report);
                if (i >= 0) {
                  this.reports.splice(i, 1);
                }
                this.toasts.push(
                  'Report revoked.',
                  'Revoked',
                  'success',
                  5000,
                  true
                );
              }
            });
        }
      })
      .finally(() => this.clearReportModalModels());
  }

  public canPerform(perm: string): boolean {
    return (
      this.hasPermission(perm) &&
      (!this.roleDiff ||
        this.roleDiff > 0 ||
        this.guild.self_member.dominance >= this.member.dominance)
    );
  }

  public canPerformAny(...perm: string[]) {
    return !!perm.find((p) => this.canPerform(p));
  }

  public muteUnmute() {
    const muted = this.member.chat_muted;
    const modal = muted ? this.modalUnmute : this.modalMute;
    const apiCall = muted ? this.api.postUnmute : this.api.postMute;
    const message = muted ? 'Revoked chat mute.' : 'Member chat muted.';

    this.openModal(modal)
      .then((res) => {
        if (res) {
          apiCall
            .call(this.api, this.guild.id, this.member.user.id, {
              attachment: this.repModalAttachment,
              reason: this.repModalReason,
              timeout: this.formattedRepTimeout,
            })
            .subscribe((resRep) => {
              if (resRep) {
                this.member.chat_muted = !muted;
                this.toasts.push(message, 'Executed', 'success', 5000, true);
              }
            });
        }
      })
      .finally(() => this.clearReportModalModels());
  }

  private openModal(modal: TemplateRef<any>): Promise<any> {
    return this.modal.open(modal, { windowClass: 'dark-modal' }).result;
  }

  private checkReason(): boolean {
    if (this.repModalReason.length < 3) {
      this.toasts.push(
        'A valid reason must be given.',
        'Error',
        'error',
        8000,
        true
      );
      return false;
    }

    return true;
  }

  private clearReportModalModels() {
    this.repModalAttachment = this.repModalReason = '';
    this.repModalType = 3;
  }

  private fetchReports(guildID: string, memberID: string) {
    this.reports = null;
    this.api.getReports(guildID, memberID).subscribe((reports) => {
      this.reports = reports || [];
    });
  }

  private get formattedRepTimeout(): string | null {
    if (!this.repModalTimeout) return null;
    const offsetHours = new Date().getTimezoneOffset() / -60;
    const offset =
      offsetHours !== 0
        ? `${offsetHours < 0 ? '-' : '+'}${padNumber(
            Math.floor(offsetHours),
            2
          )}:00`
        : 'Z';
    return `${this.repModalTimeout}:00${offset}`;
  }
}
