/** @format */

import { Injectable } from '@angular/core';

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

  push(
    text: string,
    heading: string,
    type: string = '',
    delay: number = null,
    closable: boolean = true
  ) {
    this.toasts.push({ text, heading, type, delay, closable });
  }

  remove(toast) {
    this.toasts = this.toasts.filter((t) => t !== toast);
  }
}
