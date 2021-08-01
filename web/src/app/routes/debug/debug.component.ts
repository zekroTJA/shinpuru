/** @format */

import { Component } from '@angular/core';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-debug',
  templateUrl: './debug.component.html',
  styleUrls: ['./debug.component.scss'],
})
export class DebugComponent {
  constructor(public toasts: ToastService) {}
}
