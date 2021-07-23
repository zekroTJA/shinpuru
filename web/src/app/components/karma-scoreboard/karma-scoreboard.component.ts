/** @format */

import { Component, Input } from '@angular/core';
import { GuildScoreboardEntry, Member } from 'src/app/api/api.models';

@Component({
  selector: 'app-karma-scoreboard',
  templateUrl: './karma-scoreboard.component.html',
  styleUrls: ['./karma-scoreboard.component.scss'],
})
export class KarmaScoreboardComponent {
  @Input() public guildID: string;
  @Input() public scoreboard: GuildScoreboardEntry[];
  @Input() public self: Member;

  get selfInScoreboard(): boolean {
    return !!this.scoreboard.find(
      (e) => e.member.user.id === this.self.user.id
    );
  }
}
