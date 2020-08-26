/** @format */

import { Component, OnInit } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { CommandInfo } from 'src/app/api/api.models';

@Component({
  selector: 'app-commands',
  templateUrl: './commands.component.html',
  styleUrls: ['./commands.component.sass'],
})
export class CommandsComponent implements OnInit {
  public commands: CommandInfo[];
  public groupMap: { [key: string]: CommandInfo[] } = {};

  constructor(private api: APIService) {}

  async ngOnInit() {
    try {
      const groups: { [key: string]: CommandInfo[] } = {};

      this.commands = (await this.api.getCommandInfos().toPromise()).data;
      this.commands.forEach((c) => {
        if (!(c.group in groups)) {
          groups[c.group] = [];
        }
        groups[c.group].push(c);
      });
      this.groupMap = groups;
    } catch (err) {
      console.error(err);
    }
  }

  public scrollTo(selector: string) {
    const el = document.querySelector(selector);
    if (el) {
      el.scrollIntoView();
      window.scrollBy({
        top: -60,
      });
    }
  }
}
