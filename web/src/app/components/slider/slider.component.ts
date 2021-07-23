/** @format */

import { Component, forwardRef, Output, EventEmitter } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

export const CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR: any = {
  provide: NG_VALUE_ACCESSOR,
  // tslint:disable-next-line: no-use-before-declare
  useExisting: forwardRef(() => SliderComponent),
  multi: true,
};

@Component({
  selector: 'app-slider',
  templateUrl: './slider.component.html',
  styleUrls: ['./slider.component.scss'],
  providers: [CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR],
})
export class SliderComponent implements ControlValueAccessor {
  private _value = false;

  @Output() public switch: EventEmitter<any> = new EventEmitter();

  private onTouchedCallback: () => void = () => {};
  private onChangeCallback: (_: any) => void = () => {};

  public onClick() {
    this.value = !this.value;
    this.switch.emit(this.value);
  }

  public get value(): boolean {
    return this._value;
  }

  public set value(v: boolean) {
    if (v !== this._value) {
      this._value = v;
      this.onChangeCallback(v);
    }
  }

  public onBlur() {
    this.onTouchedCallback();
  }

  public writeValue(v: boolean): void {
    if (v !== this._value) {
      this._value = v;
    }
  }

  public registerOnChange(fn: any): void {
    this.onChangeCallback = fn;
  }

  public registerOnTouched(fn: any): void {
    this.onTouchedCallback = fn;
  }
}
