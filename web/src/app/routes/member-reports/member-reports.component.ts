/** @format */

import { Component, OnInit } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import dateFormat from 'dateformat';
import { Report } from 'src/app/api/api.models';
import { ToastService } from 'src/app/components/toast/toast.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-member-reports',
  templateUrl: './member-reports.component.html',
  styleUrls: ['./member-reports.component.sass'],
})
export class MemberReportsComponent implements OnInit {
  public dateFormat = dateFormat;

  public reports: Report[];

  private guildID: string;
  private memberID: string;

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
    private toats: ToastService
  ) {}

  public async ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;
      this.memberID = params.memberid;
      this.reports = await this.api
        .getReports(this.guildID, this.memberID, 0, 100)
        .toPromise();
    });
  }
}
