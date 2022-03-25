/** @format */

import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-karma-tile',
  templateUrl: './karma-tile.component.html',
  styleUrls: ['./karma-tile.component.scss'],
})
export class KarmaTileComponent implements OnInit {
  @Input() public value: number;
  @Input() public title: string;
  @Input() public small: boolean;

  constructor() {}

  ngOnInit(): void {}

  public get classes(): object {
    return {
      bad: this.value < 0,
      good: this.value > 0,
      supreme: this.value >= 1000,
      small: this.small,
    };
  }
}
