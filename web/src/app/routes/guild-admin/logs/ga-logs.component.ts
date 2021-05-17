/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  AntiraidSettings,
  GuildLogEntry,
  JoinlogEntry,
  State,
} from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import dateFormat from 'dateformat';

interface Severity {
  name: string;
  color: string;
}

@Component({
  selector: 'app-ga-logs',
  templateUrl: './ga-logs.component.html',
  styleUrls: ['./ga-logs.component.sass'],
})
export class GuildAdminLogsComponent implements OnInit {
  public state: State;
  public entries: GuildLogEntry[];
  private guildID: string;

  public dateFormat = dateFormat;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.state = await this.api
        .getGuildSettingsLogsState(this.guildID)
        .toPromise();

      this.fetchEntries();
    });
  }

  severity(v: number): Severity {
    switch (v) {
      case -1:
        return { name: 'all', color: '' };
      case 0:
        return { name: 'debug', color: '#FF9800' };
      case 1:
        return { name: 'info', color: '#2196F3' };
      case 2:
        return { name: 'warn', color: '#FF5722' };
      case 3:
        return { name: 'error', color: '#f44336' };
      case 4:
        return { name: 'fatal', color: '#9C27B0' };
    }
  }

  async onEnabledChanged() {
    await this.api
      .postGuildSettingsLogsState(this.guildID, this.state.state)
      .toPromise();
  }

  private async fetchEntries(limit = 50, offset = 0, severity = -1) {
    this.entries = (
      await this.api
        .getGuildSettingsLogs(this.guildID, limit, offset, severity)
        .toPromise()
    ).data;
    console.log(this.entries);
  }
}
