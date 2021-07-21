/** @format */

import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

export type FilterFunc = (input: string, selection: string) => boolean;
export type FormatterFunc = (input: string, val: string) => string;

@Component({
  selector: 'app-autocomplete-input',
  templateUrl: './autocomplete-input.component.html',
  styleUrls: ['./autocomplete-input.component.scss'],
})
export class AutocompleteInputComponent implements OnInit {
  @Input() value: string;
  @Output() valueChange = new EventEmitter<string>();

  @Input() selection: string[];

  @Input() filterFunc: FilterFunc = (input, selection) =>
    selection.toLowerCase().startsWith(input.toLowerCase());
  @Input() formatterFunc: FormatterFunc = (_, value) => value;

  public selectedIndex = -1;

  constructor() {}

  ngOnInit() {}

  public get displaySelection(): Set<string> {
    return new Set(
      !this.value
        ? []
        : this.selection
            .filter((s) => this.filterFunc(this.value, s))
            .map((s) => this.formatterFunc(this.value, s))
            .filter((s) => s !== this.value)
            .slice(0, 10)
    );
  }

  public onInput(event: InputEvent) {
    if (!this.value) {
      this.selectedIndex = -1;
    }
    this.valueChange.emit((event.target as HTMLInputElement).value);
  }

  public onSelect(val: string | number) {
    if (typeof val === 'number') {
      val = Array.from(this.displaySelection)[val];
      if (!val) return;
    }
    this.valueChange.emit(val as string);
    this.selectedIndex = -1;
  }

  public onPreselect(delta: number) {
    if (delta <= -1 && this.selectedIndex <= 0) this.selectedIndex = 0;
    else if (
      delta >= 1 &&
      this.selectedIndex >= this.displaySelection.size - 1
    ) {
    } else this.selectedIndex += delta;
  }
}
