/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-ga-karma',
  templateUrl: './ga-karma.component.html',
  styleUrls: ['./ga-karma.component.sass'],
})
export class GuildAdminKarmaComponent implements OnInit {
  constructor(private route: ActivatedRoute) {}

  private guildID: string;

  ngOnInit() {
    this.route.params.subscribe((params) => {
      this.guildID = params.guildid;
    });
  }
}
