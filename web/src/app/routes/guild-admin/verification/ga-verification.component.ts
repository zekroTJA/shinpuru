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
  changeInfo: string;

  private guildID: string;
  private _enabled: boolean;

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
      this.enabled = this.state.enabled;
    });
  }

  get enabled(): boolean {
    return this._enabled;
  }

  set enabled(v: boolean) {
    this._enabled = v;

    if (this.state.enabled && !this.enabled) {
      this.changeInfo =
        'When disabling the user verification, current users which are timed out due to ' +
        'pending verification will be released. This is not reversable after applying the change!';
    } else if (!this.state.enabled && this.enabled) {
      this.changeInfo =
        'Currently unverified members on the guild will not be timed out. This is only applied to new guild members.';
    } else {
      this.changeInfo = null;
    }
  }

  async saveSettings() {
    this.state.enabled = this.enabled;
    this.state = await this.api
      .postGuildSettingsVerification(this.guildID, this.state)
      .toPromise();
    this.toasts.push('Settings were saved.', '', 'success');
    this.changeInfo = null;
  }
}
