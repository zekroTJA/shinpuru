/** @format */

import { EventEmitter, Injectable } from '@angular/core';

export interface Toast {
  text: string;
  heading: string;
  type: string;
  delay: number;
  closable: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class ToastService {
  toasts: Toast[] = [];

  onRemove = new EventEmitter<Toast>();

  push(
    text: string,
    heading: string,
    type: 'cyan' | 'red' | 'error' | 'yellow' | 'warning' | 'green' | 'success',
    delay: number = null,
    closable: boolean = true
  ): Toast {
    const t = { text, heading, type, delay, closable };
    this.toasts.push(t);
    return t;
  }

  remove(toast: Toast) {
    this.toasts = this.toasts.filter((t) => t !== toast);
    this.onRemove.emit(toast);
  }
}
