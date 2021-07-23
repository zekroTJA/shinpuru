/** @format */

import { Component, Input, OnInit } from '@angular/core';
import LocalStorageUtil from 'src/app/utils/localstorage';

@Component({
  selector: 'app-protip',
  templateUrl: './protip.component.html',
  styleUrls: ['./protip.component.scss'],
})
export class ProtipComponent implements OnInit {
  @Input() public uid: string | number;
  @Input() public top: string | number | undefined = '70px';
  @Input() public right: string | number | undefined = '20px';
  @Input() public left: string | number | undefined;
  @Input() public bottom: string | number | undefined;
  @Input() public delay = 750;

  public style = {
    display: 'none',
    top: undefined,
    left: undefined,
    bottom: undefined,
    right: undefined,
  };

  public classes = {
    show: false,
  };

  constructor() {}

  ngOnInit() {
    this.style.top = this.top;
    this.style.left = this.left;
    this.style.bottom = this.bottom;
    this.style.right = this.right;

    if (!LocalStorageUtil.get(this.localStorageKey, false)) {
      this.style.display = 'block';
      setTimeout(() => (this.classes.show = true), this.delay);
    }
  }

  public onClose() {
    this.classes.show = false;
    setTimeout(() => (this.style.display = 'none'), 750);
    LocalStorageUtil.set(this.localStorageKey, true);
  }

  private get localStorageKey(): string {
    return `PROTIP_DISMISSED_${this.uid}`;
  }
}
