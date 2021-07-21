/** @format */

import { Component, EventEmitter, Input, Output } from '@angular/core';
import { GuildStarboardEntry } from 'src/app/api/api.models';

@Component({
  selector: 'app-starboard',
  templateUrl: './starboard.component.html',
  styleUrls: ['./starboard.component.scss'],
})
export class StarboardComponent {
  @Input() public sortOrder: string;
  @Input() public guildID: string;
  @Input() public starboard: GuildStarboardEntry[];
  @Output() public sortOrderChange = new EventEmitter<any>();
}
