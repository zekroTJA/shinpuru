/** @format */

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { GuildLogEntry, State } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

interface Severity {
  name: string;
  color: string;
}

@Component({
  selector: 'app-ga-logs',
  templateUrl: './ga-logs.component.html',
  styleUrls: ['./ga-logs.component.scss'],
})
export class GuildAdminLogsComponent implements OnInit {
  public state: State;
  public entries: GuildLogEntry[];
  public entriesCount: number;
  public entriesSelected: GuildLogEntry[] = [];
  private guildID: string;

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

  public readonly limit = 100;
  public offset = 0;
  public severity = -1;

  @ViewChild('modalDeleteAll') private modalDeleteAll: TemplateRef<any>;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService,
    private modals: NgbModal
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

  getSeverity(v: number): Severity {
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

  async onEnabledChanged(v: boolean) {
    console.log(v);
    await this.api.postGuildSettingsLogsState(this.guildID, v).toPromise();
  }

  async pageDial(v: number) {
    let offset = this.offset + v * this.limit;
    if (offset < 0) offset = 0;
    if (offset >= this.entriesCount) return;
    this.offset = offset;
    this.fetchEntries();
  }

  async onSeverityChange(v: string) {
    this.severity = parseInt(v);
    this.fetchEntries();
  }

  selectUnselectEntry(v: GuildLogEntry) {
    const i = this.entriesSelected.indexOf(v);
    if (i >= 0) this.entriesSelected.splice(i, 1);
    else this.entriesSelected.push(v);
  }

  async deleteEntries() {
    if (this.entriesSelected.length > 0) {
      this.entriesSelected.forEach((e) => {
        this.api.deleteGuildSettingsLogs(this.guildID, e.id).toPromise();
        const i = this.entries.indexOf(e);
        if (i >= 0) this.entries.splice(i, 1);
      });
    } else {
      const res = await this.modals.open(this.modalDeleteAll, {
        windowClass: 'dark-modal',
      }).result;
      if (res) {
        await this.api.deleteGuildSettingsLogs(this.guildID).toPromise();
        this.entries = [];
      }
    }
  }

  get pageEnd(): number {
    const n = (this.offset + 1) * this.limit;
    if (n > this.entriesCount) return this.entriesCount;
    return n;
  }

  private async fetchEntries() {
    this.entries = (
      await this.api
        .getGuildSettingsLogs(
          this.guildID,
          this.limit,
          this.offset,
          this.severity
        )
        .toPromise()
    ).data;

    this.entriesCount = (
      await this.api
        .getGuildSettingsLogsCount(this.guildID, this.severity)
        .toPromise()
    ).count;
  }
}
