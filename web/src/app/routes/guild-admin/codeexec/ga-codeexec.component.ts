/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CodeExecSettings } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-ga-codeexec',
  templateUrl: './ga-codeexec.component.html',
  styleUrls: ['./ga-codeexec.component.scss'],
})
export class GuildAdminCodeExecComponent implements OnInit {
  state: CodeExecSettings;

  private guildID: string;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;
      this.state = await this.api
        .getGuildSettingsCodeExec(this.guildID)
        .toPromise();
    });
  }

  async saveSettings() {
    this.state = await this.api
      .postGuildSettingsCodeExec(this.guildID, this.state)
      .toPromise();
    this.toasts.push('Settings were saved.', '', 'success');
  }

  get isJdoodle(): boolean {
    return this.state?.type === 'jdoodle';
  }

  get emptyCreds(): boolean {
    return !this.state.jdoodle_clientid && !this.state.jdoodle_clientsecret;
  }

  resetCredentials() {
    this.state.jdoodle_clientid = this.state.jdoodle_clientsecret = null;
  }
}
