/** @format */

import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';
import { Report, Member } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';

const typeColors = ['#d81b60', '#e53935', '#009688', '#fb8c00', '#8e24aa'];

@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.scss'],
})
export class ReportComponent implements OnInit {
  @Input() public report: Report;
  @Input() public victim: Member;
  @Input() public executor: Member;
  @Input() public allowRevoke: boolean;

  @Output() public revoke = new EventEmitter<any>();

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

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
