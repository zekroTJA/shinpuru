import { Injectable } from '@angular/core';
import { UpdateInfoResponse } from '../api/api.models';
import { APIService } from '../api/api.service';
import { ToastService } from '../components/toast/toast.service';
import LocalStorageUtil from '../utils/localstorage';

@Injectable({
  providedIn: 'root',
})
export class UpdateService {
  constructor(private api: APIService, private toasts: ToastService) {}

  async check() {
    const self = await this.api.getSelfUser().toPromise();
    if (!self.bot_owner) return;

    const res = await this.api.getUpdateInfo().toPromise();
    if (!res.isold || LocalStorageUtil.get(lsKey(res), false)) return;
    const toast = this.toasts.push(
      `A new version of shinpuru is available! Current version is ${res.current_str} and latest version is ${res.latest_str}.`,
      'New Update available!',
      'yellow',
      999999999
    );

    const sub = this.toasts.onRemove.subscribe((t) => {
      if (t === toast) {
        LocalStorageUtil.set(lsKey(res), true);
        sub.unsubscribe();
      }
    });
  }
}

function lsKey(res: UpdateInfoResponse): string {
  return `UPDATEINFO_${res.current_str}_${res.latest_str}`;
}
