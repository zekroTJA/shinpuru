/** @format */

import { Pipe, PipeTransform } from '@angular/core';
import { UnbanRequestState } from '../api/api.models';

@Pipe({
  name: 'unbanrequestState',
})
export class UnbanrequestStatePipe implements PipeTransform {
  transform(state: UnbanRequestState): string {
    switch (state) {
      case UnbanRequestState.PENDING:
        return 'pending';
      case UnbanRequestState.DECLINED:
        return 'declined';
      case UnbanRequestState.ACCEPTED:
        return 'accepted';
      default:
        return 'unknown';
    }
  }
}
