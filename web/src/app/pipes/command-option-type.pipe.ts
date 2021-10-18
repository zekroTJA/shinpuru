import { Pipe, PipeTransform } from '@angular/core';
import { CommandOption, CommandOptionType } from '../api/api.models';

@Pipe({
  name: 'commandOptionType',
})
export class CommandOptionTypePipe implements PipeTransform {
  transform(opt: CommandOption): string {
    return (Object.values(CommandOptionType)[opt.type] as string).toLowerCase();
  }
}
