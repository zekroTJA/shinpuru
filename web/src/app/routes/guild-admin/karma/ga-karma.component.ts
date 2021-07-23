/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { KarmaRule, KarmaSettings, Member } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-ga-karma',
  templateUrl: './ga-karma.component.html',
  styleUrls: ['./ga-karma.component.scss'],
})
export class GuildAdminKarmaComponent implements OnInit {
  public karmaSettings: KarmaSettings;
  public blocklist: Member[] = null;
  public rules: KarmaRule[] = null;
  public newRule = {
    trigger: 0,
    action: 'TOGGLE_ROLE',
    value: 0,
  } as KarmaRule;
  private guildID: string;

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.karmaSettings = await this.api
        .getGuildSettingsKarma(this.guildID)
        .toPromise();

      await this.fetchRules();
      await this.fetchBlocklist();
    });
  }

  private async fetchBlocklist() {
    this.blocklist = (
      await this.api.getGuildSettingsKarmaBlocklist(this.guildID).toPromise()
    ).data;
  }

  private async fetchRules() {
    this.rules = (
      await this.api.getGuildSettingsKarmaRules(this.guildID).toPromise()
    ).data;
  }

  public async onSave() {
    try {
      this.karmaSettings.emotes_increase =
        this.karmaSettings.emotes_increase.filter((e) => !!e);
      this.karmaSettings.emotes_decrease =
        this.karmaSettings.emotes_decrease.filter((e) => !!e);
      await this.api
        .postGuildSettingsKarma(this.guildID, this.karmaSettings)
        .toPromise();
      this.toasts.push(
        'Karma settings saved.',
        'Settings saved',
        'green',
        4000
      );
    } catch {}
  }

  public onIncChange(event: any) {
    this.karmaSettings.emotes_increase = event.target.value
      .split(',')
      .map((v: string) => v.trim());
  }

  public onDecChange(event: any) {
    this.karmaSettings.emotes_decrease = event.target.value
      .split(',')
      .map((v: string) => v.trim());
  }

  public async onMemberBlock(id: string) {
    if (!id) return;
    try {
      await this.api
        .putGuildSettingsKarmaBlocklist(this.guildID, id)
        .toPromise();
      await this.fetchBlocklist();
      this.toasts.push(
        'Member added to karma blocklist.',
        'Member blocked',
        'green'
      );
    } catch {}
  }

  public async onMemberUnblock(id: string) {
    try {
      await this.api
        .deleteGuildSettingsKarmaBlocklist(this.guildID, id)
        .toPromise();
      const i = this.blocklist.findIndex((m) => m.user.id === id);
      if (i >= 0) {
        this.blocklist.splice(i, 1);
      }
      this.toasts.push(
        'Member removed from karma blocklist.',
        'Member unblocked',
        'green'
      );
    } catch {}
  }

  public onRuleTriggerChange(v: string) {
    this.newRule.trigger = parseInt(v);
  }

  public async applyRule() {
    if (
      this.newRule.action !== 'KICK' &&
      this.newRule.action !== 'BAN' &&
      !this.newRule.argument
    )
      return;
    this.newRule.guildid = this.guildID;
    const res = await this.api
      .createGuildSettingsKarmaRules(this.newRule)
      .toPromise();
    this.rules.push(res);
    this.newRule = {
      trigger: 0,
      action: 'TOGGLE_ROLE',
      value: 0,
    } as KarmaRule;
  }

  public ruleTrigger(v: number): string {
    switch (v) {
      case 0:
        return 'drops below';
      case 1:
        return 'rises above';
    }
  }

  public ruleAction(v: string): string {
    switch (v) {
      case 'TOGGLE_ROLE':
        return 'toggle role';
      case 'KICK':
        return 'kick member';
      case 'BAN':
        return 'ban member';
      case 'SEND_MESSAGE':
        return 'send message';
    }
  }

  public async removeRule(r: KarmaRule) {
    await this.api.deleteGuildSettingsKarmaRules(r).toPromise();
    const i = this.rules.indexOf(r);
    if (i >= 0) {
      this.rules.splice(i, 1);
    }
  }
}
