/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { Guild } from 'src/app/api/api.models';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.sass'],
})
export class HomeComponent {
  public guilds: Guild[] = [];

  constructor(private api: APIService, public spinner: SpinnerService) {
    this.api.getGuilds().subscribe((guilds) => {
      this.guilds = guilds;
      this.spinner.stop('spinner-load-guilds');
    });
  }
}
