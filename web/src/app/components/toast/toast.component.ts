/** @format */

import { Component, Input, OnInit, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-toast',
  templateUrl: './toast.component.html',
  styleUrls: ['./toast.component.scss'],
})
export class ToastComponent implements OnInit {
  @Input() public heading: string;
  @Input() public closable = true;
  @Input() public color: string;
  @Input() public delay: number;

  @Output() public hidden: EventEmitter<any> = new EventEmitter();

  public colorClass: string;
  public dark = true;
  public visible = true;
  public fadeOut = false;

  constructor() {}

  ngOnInit() {
    setTimeout(() => this.hide(false), (this.delay ?? 6000) + 600);

    switch (this.color) {
      case 'cyan':
        this.colorClass = 'c-cyan';
        this.dark = false;
        break;

      case 'red':
      case 'error':
        this.colorClass = 'c-red';
        this.dark = false;
        break;

      case 'yellow':
      case 'warning':
        this.colorClass = 'c-yellow';
        this.dark = true;
        break;

      case 'green':
      case 'success':
        this.colorClass = 'c-green';
        this.dark = false;
        break;
    }
  }

  public hide(active: boolean) {
    if (!this.visible) {
      return;
    }

    this.fadeOut = true;
    setTimeout(() => {
      this.visible = false;
      this.hidden.emit(active);
    }, 300);
  }
}
