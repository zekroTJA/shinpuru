/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { GuildSettingsVerification } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-ga-verification',
  templateUrl: './ga-verification.component.html',
  styleUrls: ['./ga-verification.component.scss'],
})
export class GuildAdminVerificationComponent implements OnInit {
  state: GuildSettingsVerification;

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
        .getGuildSettingsVerification(this.guildID)
        .toPromise();
    });
  }

  async saveSettings() {
    this.state = await this.api
      .postGuildSettingsVerification(this.guildID, this.state)
      .toPromise();
    this.toasts.push('Settings were saved.', '', 'success');
  }
}
