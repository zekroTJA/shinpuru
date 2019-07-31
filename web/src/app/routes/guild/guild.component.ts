/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
// import { Guild } from 'src/app/api/api.models';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import { ActivatedRoute } from '@angular/router';
import { Guild } from 'src/app/api/api.models';

@Component({
  selector: 'app-guild',
  templateUrl: './guild.component.html',
  styleUrls: ['./guild.component.sass'],
})
export class GuildComponent {
  public guild: Guild;

  constructor(
    public api: APIService,
    public spinner: SpinnerService,
    private route: ActivatedRoute
  ) {
    const guildID = this.route.snapshot.paramMap.get('id');
    this.api.getGuild(guildID).subscribe((guild) => {
      this.guild = guild;
      this.spinner.stop('spinner-load-guild');
    });
  }
}
