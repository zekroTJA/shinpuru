/** @format */

import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'capText',
})
export class CapTextPipe implements PipeTransform {
  transform(value: string, max: number): string {
    return value.length > max ? value.substring(0, max - 3) + '...' : value;
  }
}
