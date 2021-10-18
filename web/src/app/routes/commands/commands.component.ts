/** @format */

import { Component, OnInit } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import {
  CommandInfo,
  CommandOptionType,
  SubPermission,
} from 'src/app/api/api.models';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-commands',
  templateUrl: './commands.component.html',
  styleUrls: ['./commands.component.scss'],
})
export class CommandsComponent implements OnInit {
  commands: CommandInfo[];
  groupMap: { [key: string]: CommandInfo[] } = {};
  lastSelected: string;

  constructor(private api: APIService, private route: ActivatedRoute) {}

  async ngOnInit() {
    try {
      this.commands = (await this.api.getCommandInfos().toPromise()).data;
      this.fetchGroups();

      this.route.fragment.subscribe((fragment) => {
        setTimeout(() => this.scrollTo(`#${fragment}`), 500);
      });
    } catch (err) {
      console.error(err);
    }
  }

  getCommandSubPermTerm(cmd: CommandInfo, sp: SubPermission): string {
    if (sp.term.startsWith('/')) return sp.term.substr(1);
    return cmd.domain + '.' + sp.term;
  }

  scrollTo(selector: string) {
    const el = document.querySelector(selector);
    if (el) {
      el.scrollIntoView({
        block: 'center',
      });
    }
  }

  onScrollToTop() {
    window.scrollTo({
      top: 0,
    });
  }

  onSearchBarChange(e: InputEvent) {
    const val = (e.currentTarget as HTMLInputElement).value;
    this.fetchGroups(val);
  }

  hasSubCommands(cmd: CommandInfo) {
    return (
      cmd.options.length != 0 &&
      cmd.options[0].type === CommandOptionType.SUBCOMMAND
    );
  }

  private fetchGroups(filter?: string) {
    const groups: { [key: string]: CommandInfo[] } = {};

    this.commands
      .filter((c) => !filter || this.commandFilterFunc(c, filter))
      .forEach((c) => {
        if (!(c.group in groups)) {
          groups[c.group] = [];
        }
        groups[c.group].push(c);
      });
    this.groupMap = groups;
  }

  private commandFilterFunc(c: CommandInfo, f: string): boolean {
    f = f.toLowerCase();

    if (c.name.toLowerCase().includes(f)) {
      return true;
    }

    if (c.domain.toLowerCase().includes(f)) {
      return true;
    }

    if (c.description.toLowerCase().includes(f)) {
      return true;
    }

    return false;
  }
}
