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

  constructor(private api: APIService) {}

  async ngOnInit() {
    try {
      this.commands = (await this.api.getCommandInfos().toPromise()).data;
    } catch {}
  }
}
