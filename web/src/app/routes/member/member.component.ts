/** @format */

import { Component, ViewChildren, TemplateRef, ViewChild } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import {
  Member,
  Guild,
  Role,
  Report,
  ReportRequest,
} from 'src/app/api/api.models';
import dateFormat from 'dateformat';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { ToastService } from 'src/app/components/toast/toast.service';
import { rolePosDiff } from 'src/app/utils/utils';

@Component({
  selector: 'app-member-route',
  templateUrl: './member.component.html',
  styleUrls: ['./member.component.sass'],
})
export class MemberRouteComponent {
  public member: Member;
  public guild: Guild;
  public perm: string[];

  public reports: Report[];
  public permissionsAllowed: string[] = [];

  public roleDiff: number;

  public dateFormat = dateFormat;

  @ViewChild('modalReport') private modalReport: TemplateRef<any>;
  @ViewChild('modalKick') private modalKick: TemplateRef<any>;
  @ViewChild('modalBan') private modalBan: TemplateRef<any>;

  public repModalType = 3;
  public repModalReason = '';
  public repModalAttachment = '';

  constructor(
    public modal: NgbModal,
    private api: APIService,
    private spinner: SpinnerService,
    private toasts: ToastService,
    private route: ActivatedRoute,
    private router: Router
  ) {
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

  public report() {
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
        this.clearReportModalModels();
      })
      .catch(() => this.clearReportModalModels());
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
        this.clearReportModalModels();
      })
      .catch(() => this.clearReportModalModels());
  }

  public ban() {
    this.openModal(this.modalBan)
      .then((res) => {
        if (res && this.checkReason()) {
          this.api
            .postBan(this.guild.id, this.member.user.id, {
              attachment: this.repModalAttachment,
              reason: this.repModalReason,
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
        this.clearReportModalModels();
      })
      .catch(() => this.clearReportModalModels());
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
}
