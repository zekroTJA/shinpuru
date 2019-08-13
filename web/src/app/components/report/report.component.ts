/** @format */

import { Component, Input, OnInit } from '@angular/core';
import { Report, Member } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import dateFormat from 'dateformat';

const typeColors = ['#d81b60', '#e53935', '#009688', '#fb8c00', '#8e24aa'];

@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.sass'],
})
export class ReportComponent implements OnInit {
  @Input() public report: Report;
  @Input() public victim: Member;
  @Input() public executor: Member;

  public dateFormat = dateFormat;

  constructor(private api: APIService) {}

  ngOnInit() {
    if (!this.executor) {
      this.api
        .getGuildMember(this.report.guild_id, this.report.executor_id, true)
        .subscribe((u) => {
          this.executor = u;
        });
    }

    if (!this.victim) {
      this.api
        .getGuildMember(this.report.guild_id, this.report.victim_id, true)
        .subscribe((u) => {
          this.victim = u;
        });
    }
  }

  public isDiscordAttachment(url: string): boolean {
    return url.startsWith('https://cdn.discordapp.com/attachments/');
  }

  public get typeColor(): string {
    return typeColors[this.report.type] || 'gray';
  }
}
