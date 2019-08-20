/** @format */

import { Component, forwardRef, Output, EventEmitter } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

export const CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR: any = {
  provide: NG_VALUE_ACCESSOR,
  // tslint:disable-next-line: no-use-before-declare
  useExisting: forwardRef(() => SpinnerButtonComponent),
  multi: true,
};

@Component({
  selector: 'app-spinner-button',
  templateUrl: './spinnerButton.component.html',
  styleUrls: ['./spinnerButton.component.sass'],
  providers: [CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR],
})
export class SpinnerButtonComponent implements ControlValueAccessor {
  private _loading = false;

  @Output() click: EventEmitter<any> = new EventEmitter<any>();

  private onTouchedCallback: () => void = () => {};
  private onChangeCallback: (_: any) => void = () => {};

  public onClick(event: any) {
    this.click.emit(event);
  }

  public get loading(): boolean {
    return this._loading;
  }

  public set loading(v: boolean) {
    if (v !== this._loading) {
      this._loading = v;
      this.onChangeCallback(v);
    }
  }

  public onBlur() {
    this.onTouchedCallback();
  }

  public writeValue(v: boolean): void {
    if (v !== this._loading) {
      this._loading = v;
    }
  }

  public registerOnChange(fn: any): void {
    this.onChangeCallback = fn;
  }

  public registerOnTouched(fn: any): void {
    this.onTouchedCallback = fn;
  }
}
