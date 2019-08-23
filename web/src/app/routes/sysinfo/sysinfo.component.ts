/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { SystemInfo } from 'src/app/api/api.models';
import dateFormat from 'dateformat';

@Component({
  selector: 'app-sysinfo',
  templateUrl: './sysinfo.component.html',
  styleUrls: ['./sysinfo.component.sass'],
})
export class SysInfoComponent {
  public sysinfo: SystemInfo;
  public uptime: number;

  public dateFormat = dateFormat;

  constructor(private api: APIService) {
    this.refresh();
  }

  public refresh() {
    this.api.getSystemInfo().subscribe((sysinfo) => {
      this.sysinfo = sysinfo;
      this.uptime = sysinfo.uptime;
    });
  }

  public byteCountFormatter(bc: number) {
    const k = 1024;
    const fix = 2;
    if (bc < k) {
      return `${bc} B`;
    }
    if (bc < k * k) {
      return `${(bc / k).toFixed(fix)} kiB`;
    }
    if (bc < k * k * k) {
      return `${(bc / k / k).toFixed(fix)} MiB`;
    }
    if (bc < k * k * k * k) {
      return `${(bc / k / k / k).toFixed(fix)} GiB`;
    }
    return `${(bc / k / k / k / k).toFixed(fix)} TiB`;
  }

  public toDDHHMMSS(secs: number) {
    const dd = this.padFront(Math.floor(secs / 86400), 2, 0);
    const hh = this.padFront(Math.floor((secs % 86400) / 3600), 2, 0);
    const mm = this.padFront(Math.floor(((secs % 86400) % 3600) / 60), 2, 0);
    const ss = this.padFront(Math.floor(((secs % 86400) % 3600) % 60), 2, 0);
    return `${dd}:${hh}:${mm}:${ss}`;
  }

  private padFront(num: any, len: number, char: any) {
    num = num.toString();
    char = char.toString();
    while (num.length < len) {
      num = char + num;
    }
    return num;
  }
}
