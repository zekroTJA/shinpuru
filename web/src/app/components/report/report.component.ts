/** @format */

import { Component, Input, OnInit } from '@angular/core';
import { Report, Member } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import dateFormat from 'dateformat';

@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.sass'],
})
export class ReportComponent implements OnInit {
  @Input() public report: Report;
  @Input() public victim: Member;

  public executor: Member;

  public dateFormat = dateFormat;

  constructor(public api: APIService) {}

  ngOnInit() {
    this.api
      .getGuildMember(this.report.guild_id, this.report.executor_id)
      .subscribe((u) => {
        this.executor = u;
      });

    if (!this.victim) {
      this.api
        .getGuildMember(this.report.guild_id, this.report.victim_id)
        .subscribe((u) => {
          this.victim = u;
        });
    }
  }
}
