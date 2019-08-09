/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { Presence } from 'src/app/api/api.models';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.sass'],
})
export class SettingsComponent {
  public presence: Presence;

  constructor(private api: APIService) {
    this.api.getPresence().subscribe((presence) => {
      this.presence = presence;
    });
  }

  public updatePresence() {
    this.api.postPresence(this.presence).subscribe((presence) => {
      if (presence) {
        this.presence = presence;
      }
    });
  }
}
