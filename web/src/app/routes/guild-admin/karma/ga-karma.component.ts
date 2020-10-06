/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { KarmaSettings } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';

@Component({
  selector: 'app-ga-karma',
  templateUrl: './ga-karma.component.html',
  styleUrls: ['./ga-karma.component.sass'],
})
export class GuildAdminKarmaComponent implements OnInit {
  public karmaSettings: KarmaSettings;
  private guildID: string;

  constructor(private route: ActivatedRoute, private api: APIService) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.karmaSettings = await this.api
        .getGuildSettingsKarma(this.guildID)
        .toPromise();
    });
  }
}
