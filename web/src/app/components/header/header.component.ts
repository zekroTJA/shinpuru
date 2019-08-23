/** @format */

import { Component, ViewChild, TemplateRef, OnInit } from '@angular/core';
import { APIService } from '../../api/api.service';
import { User } from '../../api/api.models';
import { PopupElement } from '../popup/popup.component';
import { Router } from '@angular/router';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.sass'],
})
export class HeaderComponent implements OnInit {
  @ViewChild('logout') private logoutTemplate: TemplateRef<any>;
  @ViewChild('settings') private settingsTemplate: TemplateRef<any>;
  @ViewChild('sysinfo') private sysinfoTemplate: TemplateRef<any>;

  public selfUser: User;

  public popupVisible = false;
  public popupElements = [];

  constructor(private api: APIService, private router: Router) {
    this.api.getSelfUser().subscribe((user) => {
      this.selfUser = user;
      if (user.bot_owner) {
        this.popupElements.push({
          el: this.settingsTemplate,
          action: this.settings.bind(this),
        } as PopupElement);
      }
    });
  }

  ngOnInit() {
    this.popupElements.push(
      {
        el: this.logoutTemplate,
        action: this.logout.bind(this),
      } as PopupElement,
      {
        el: this.sysinfoTemplate,
        action: this.sysinfo.bind(this),
      } as PopupElement
    );
  }

  public get routes(): string[][] {
    const rts = this.router.url.split('/').filter((e) => e.length > 0);
    let path = '';
    return rts.map((r) => [r, (path += '/' + r)]);
  }

  public popupClose(e: any) {
    if (e.target.className !== 'logout-btn') {
      this.popupVisible = false;
    }
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
    this.router.navigate(['/sysinfo']);
  }
}
