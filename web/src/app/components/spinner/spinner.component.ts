/** @format */

import { Component, Input, OnInit } from '@angular/core';
import { SpinnerService } from './spinner.service';

@Component({
  selector: 'app-spinner',
  templateUrl: './spinner.component.html',
  styleUrls: ['./spinner.component.scss'],
})
export class SpinnerComponent implements OnInit {
  @Input() public id: string;
  @Input() public started = false;
  @Input() public small: boolean;

  public className = 'spinner-border';

  constructor(public spinnerService: SpinnerService) {}

  ngOnInit() {
    if (this.started) {
      this.spinnerService.start(this.id);
    }

    if (this.small) {
      this.className = 'spinner-border spinner-border-sm';
    }
  }
}
