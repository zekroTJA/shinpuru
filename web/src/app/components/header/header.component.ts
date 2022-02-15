/** @format */

import { Component, ViewChild, TemplateRef, OnInit } from '@angular/core';
import { APIService } from '../../api/api.service';
import { User } from '../../api/api.models';
import { PopupElement } from '../popup/popup.component';
import { Router } from '@angular/router';
import { LoadingBarService } from './loadingbar.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnInit {
  @ViewChild('logout', { static: true })
  private logoutTemplate: TemplateRef<any>;

  @ViewChild('settings', { static: true })
  private settingsTemplate: TemplateRef<any>;

  @ViewChild('sysinfo', { static: true })
  private sysinfoTemplate: TemplateRef<any>;

  @ViewChild('usersettings', { static: true })
  private usersettingsTemplate: TemplateRef<any>;

  @ViewChild('documentation', { static: true })
  private documentationTemplate: TemplateRef<any>;

  @ViewChild('commands', { static: true })
  private commandsTemplate: TemplateRef<any>;

  public selfUser: User;

  public popupVisible = false;
  public popupElements = [];

  constructor(
    private api: APIService,
    private router: Router,
    public loadingBar: LoadingBarService
  ) {}

  ngOnInit() {
    this.api.getSelfUser().subscribe((user) => {
      this.selfUser = user;
      if (user.bot_owner) {
        this.popupElements.push({
          el: this.settingsTemplate,
          action: this.settings.bind(this),
        } as PopupElement);
      }
    });

    this.popupElements.push(
      {
        el: this.logoutTemplate,
        action: this.logout.bind(this),
      } as PopupElement,
      {
        el: this.sysinfoTemplate,
        action: this.sysinfo.bind(this),
      } as PopupElement,
      {
        el: this.documentationTemplate,
        action: this.documentaiton.bind(this),
      } as PopupElement,
      {
        el: this.commandsTemplate,
        action: this.commands.bind(this),
      } as PopupElement,
      {
        el: this.usersettingsTemplate,
        action: this.usersettings.bind(this),
      } as PopupElement
    );
  }

  public get routes(): string[][] {
    const rts = this.router.url.split('/').filter((e) => e.length > 0);
    let path = '';
    return rts.map((r) => [r, (path += '/' + r)]);
  }

  public popupClose(e: any) {
    const target = e.target as HTMLElement;
    if (target.id !== 'user-info' && target.parentElement?.id !== 'user-info') {
      this.popupVisible = false;
    }
  }

  public onLogin() {
    this.router.navigate(['/login']);
  }

  public get isLoginPage(): boolean {
    return window.location.pathname.startsWith('/login');
  }

  private logout() {
    this.api.logout().subscribe(() => {
      window.location.assign('/');
    });
  }

  private settings() {
    this.router.navigate(['/settings']);
  }

  private sysinfo() {
    this.router.navigate(['/info']);
  }

  private usersettings() {
    this.router.navigate(['/usersettings']);
  }

  private documentaiton() {
    window.open('https://github.com/zekroTJA/shinpuru/wiki', '_blank');
  }

  private commands() {
    this.router.navigate(['/commands']);
  }
}
