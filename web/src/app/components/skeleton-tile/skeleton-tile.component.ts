import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-skeleton-tile',
  templateUrl: './skeleton-tile.component.html',
  styleUrls: ['./skeleton-tile.component.scss'],
})
export class SkeletonTileComponent implements OnInit {
  @Input() width: string | number = '100%';
  @Input() height: string | number = '100%';
  @Input() margin: string | number = 0;
  @Input() delay: string | number = 0;

  constructor() {}

  ngOnInit(): void {}
}
