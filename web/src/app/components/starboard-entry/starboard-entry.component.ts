/** @format */

import { Component, Input, OnInit } from '@angular/core';
import { GuildStarboardEntry } from 'src/app/api/api.models';

@Component({
  selector: 'app-starboard-entry',
  templateUrl: './starboard-entry.component.html',
  styleUrls: ['./starboard-entry.component.scss'],
})
export class StarboardEntryComponent implements OnInit {
  @Input() public entry: GuildStarboardEntry;

  constructor() {}

  ngOnInit(): void {}
}
