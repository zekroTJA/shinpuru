/** @format */

import { Component, OnInit, forwardRef } from '@angular/core';
import { NG_VALUE_ACCESSOR, ControlValueAccessor } from '@angular/forms';

export const CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR: any = {
  provide: NG_VALUE_ACCESSOR,
  // tslint:disable-next-line: no-use-before-declare
  useExisting: forwardRef(() => SpoilerComponent),
  multi: true,
};

@Component({
  selector: 'app-spoiler',
  templateUrl: './spoiler.component.html',
  styleUrls: ['./spoiler.component.scss'],
  providers: [CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR],
})
export class SpoilerComponent implements ControlValueAccessor {
  private _value = false;

  private onTouchedCallback: () => void = () => {};
  private onChangeCallback: (_: any) => void = () => {};

  constructor() {}

  public get rotation(): string {
    return `rotate(${this._value ? 90 : -90}deg)`;
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

  public toggle() {
    this.value = !this.value;
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
