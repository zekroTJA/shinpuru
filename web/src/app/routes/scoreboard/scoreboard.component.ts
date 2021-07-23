/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { APIService } from 'src/app/api/api.service';
import { GuildScoreboardEntry } from 'src/app/api/api.models';

@Component({
  selector: 'app-scoreboard-route',
  templateUrl: './scoreboard.component.html',
  styleUrls: ['./scoreboard.component.scss'],
})
export class ScoreboardComponent implements OnInit {
  public scoreboard: GuildScoreboardEntry[];
  public guildID: string;

  constructor(private route: ActivatedRoute, private api: APIService) {}

  public async ngOnInit() {
    this.guildID = this.route.snapshot.paramMap.get('guildid');

    try {
      this.scoreboard = (
        await this.api.getGuildScoreboard(this.guildID).toPromise()
      )?.data;
    } catch {}
  }
}
