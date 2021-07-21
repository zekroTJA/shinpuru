/** @format */

import { Component, OnInit } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { CommandInfo, SubPermission } from 'src/app/api/api.models';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-commands',
  templateUrl: './commands.component.html',
  styleUrls: ['./commands.component.scss'],
})
export class CommandsComponent implements OnInit {
  public commands: CommandInfo[];
  public groupMap: { [key: string]: CommandInfo[] } = {};

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

  public getCommandSubPermTerm(cmd: CommandInfo, sp: SubPermission): string {
    if (sp.term.startsWith('/')) return sp.term.substr(1);
    return cmd.domain_name + '.' + sp.term;
  }

  public scrollTo(selector: string) {
    const el = document.querySelector(selector);
    if (el) {
      el.scrollIntoView({
        block: 'center',
      });
    }
  }

  public onScrollToTop() {
    window.scrollTo({
      top: 0,
    });
  }

  public onSearchBarChange(e: InputEvent) {
    const val = (e.currentTarget as HTMLInputElement).value;
    this.fetchGroups(val);
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

    if (c.invokes.find((i) => i.toLowerCase().includes(f))) {
      return true;
    }

    if (c.domain_name.toLowerCase().includes(f)) {
      return true;
    }

    if (c.description.toLowerCase().includes(f)) {
      return true;
    }

    return false;
  }
}
