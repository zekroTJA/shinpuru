/** @format */

import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class SpinnerService {
  public running: string[] = [];

  public start(id: string) {
    if (!this.running.includes(id)) {
      this.running.push(id);
    }
  }

  public stop(id: string) {
    this.running = this.running.filter((r) => r !== id);
  }
}
