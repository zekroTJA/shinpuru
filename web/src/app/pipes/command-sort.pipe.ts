/** @format */

import { Pipe, PipeTransform } from '@angular/core';
import { CommandInfo } from 'src/app/api/api.models';

@Pipe({
  name: 'commandSort',
})
export class CommandSortPipe implements PipeTransform {
  transform(list: CommandInfo[]): CommandInfo[] {
    return list.sort((a, b) => (a.name < b.name ? -1 : 1));
  }
}
