/** @format */

import { Component, OnInit, forwardRef, Input } from '@angular/core';
import { NG_VALUE_ACCESSOR, ControlValueAccessor } from '@angular/forms';

export const CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR: any = {
  provide: NG_VALUE_ACCESSOR,
  // tslint:disable-next-line: no-use-before-declare
  useExisting: forwardRef(() => TagsInputComponent),
  multi: true,
};

@Component({
  selector: 'app-tags',
  templateUrl: './tagsinput.component.html',
  styleUrls: ['./tagsinput.component.scss'],
  providers: [CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR],
})
export class TagsInputComponent implements ControlValueAccessor {
  private _value: any[] = [];

  public suggested: any[] = [];
  public inputTxt: string;

  @Input() public available: any[] = [];
  @Input() public placeholder = '';
  @Input() public formatter: (v: any) => string = (v: any) => v.toString();
  @Input() public filter: (v: any, inpt: string) => boolean = (
    v: any,
    inpt: string
  ) => {
    return this.formatter(v).toLowerCase().includes(inpt.toLowerCase());
  };
  @Input() public invalidFilter: (v: any) => boolean = (v: any) => false;

  private onTouchedCallback: () => void = () => {};
  private onChangeCallback: (_: any) => void = () => {};

  constructor() {}

  public onInput(e: any) {
    const val = e.target.value;
    if (val.length === 0) {
      this.suggested = [];
    } else {
      this.suggested = this.available.filter((el) => {
        return this.filter(el, val) && !this.value.includes(el);
      });
    }
  }

  public onAdd(e: any) {
    if (this.invalidFilter(e)) return;
    this.value.push(e);
    this.suggested = [];
    this.inputTxt = '';
  }

  public onRemove(e: any) {
    const i = this.value.findIndex((el) => el === e);
    if (i > -1) {
      this.value.splice(i, 1);
    }
  }

  public get value(): any[] {
    return this._value;
  }

  public set value(v: any[]) {
    if (v !== this._value) {
      this._value = v;
      this.onChangeCallback(v);
    }
  }

  public onBlur() {
    this.onTouchedCallback();
  }

  public writeValue(v: any[]): void {
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
