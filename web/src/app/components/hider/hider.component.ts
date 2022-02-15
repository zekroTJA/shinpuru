import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';

@Component({
  selector: 'app-hider',
  templateUrl: './hider.component.html',
  styleUrls: ['./hider.component.scss'],
})
export class HiderComponent implements OnInit {
  @Input() value: string;
  @Input() hiddenValue: string = 'Hover to show';
  @Input() styleHidden: { [key: string]: any };
  @Input() styleShown: { [key: string]: any };

  @ViewChild('valueBox') valueBox: ElementRef;

  showCbInfo = false;

  private _isHidden = true;

  constructor() {}

  ngOnInit() {}

  get containerWidth(): string {
    const width =
      this.hiddenValue.length > this.value.length
        ? this.hiddenValue.length
        : this.value.length;
    return `${width + 2}ch`;
  }

  get isHidden(): boolean {
    return this._isHidden;
  }

  set isHidden(v: boolean) {
    this._isHidden = v;

    if (v) {
      this.valueBox.nativeElement.value = this.hiddenValue;
      this.showCbInfo = false;
    } else {
      this.valueBox.nativeElement.value = this.value;
      this.valueBox.nativeElement.focus();
      this.valueBox.nativeElement.setSelectionRange(0, this.value.length);
      navigator.clipboard.writeText(this.value).then(() => {
        this.showCbInfo = true;
      });
    }
  }
}
