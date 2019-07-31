/** @format */

import { Component } from '@angular/core';
import { ToastService } from './components/toast/toast.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass'],
})
export class AppComponent {
  title = 'shinpuru Web Interface';

  constructor(public toasts: ToastService) {}
}
