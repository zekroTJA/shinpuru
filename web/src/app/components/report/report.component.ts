/** @format */

import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';
import { FlatUser, Report } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';

const typeColors = ['#d81b60', '#e53935', '#009688', '#fb8c00', '#8e24aa'];

@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.scss'],
})
export class ReportComponent implements OnInit {
  @Input() public report: Report;
  @Input() public allowRevoke: boolean;
  @Input() public fetchData: boolean;

  @Output() public revoke = new EventEmitter<any>();

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

  constructor(private api: APIService) {}

  async ngOnInit() {
    try {
      if (!this.report.executor && this.fetchData)
        this.report.executor = await this.api
          .getUser(this.report.executor_id)
          .toPromise();

      if (!this.report.victim && this.fetchData)
        this.report.victim = await this.api
          .getUser(this.report.victim_id)
          .toPromise();
    } catch {}
  }

  public isDiscordAttachment(url: string): boolean {
    return url.startsWith('https://cdn.discordapp.com/attachments/');
  }

  public get typeColor(): string {
    return typeColors[this.report.type] || 'gray';
  }
}
