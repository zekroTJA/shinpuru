/** @format */

import { Pipe, PipeTransform } from '@angular/core';
import { KeyValue } from '@angular/common';

@Pipe({
  name: 'keyValue',
})
export class KeyValuePipe implements PipeTransform {
  transform<T>(
    map: { [key: string]: T },
    ...args: any[]
  ): KeyValue<string, T>[] {
    return Object.keys(map).map(
      (k) => ({ key: k, value: map[k] } as KeyValue<string, T>)
    );
  }
}
