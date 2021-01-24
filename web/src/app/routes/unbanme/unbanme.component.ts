/** @format */

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { Guild, UnbanRequest } from 'src/app/api/api.models';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

interface ExtendedUnbanRequest extends UnbanRequest {
  guild: Guild;
}

@Component({
  selector: 'app-unbanme',
  templateUrl: './unbanme.component.html',
  styleUrls: ['./unbanme.component.sass'],
})
export class UnbanmeComponent implements OnInit {
  @ViewChild('modalRequest') private modalRequest: TemplateRef<any>;

  public bannedGuilds: Guild[] = null;
  public requests: UnbanRequest[] = [];
  public newRequest = {} as ExtendedUnbanRequest;

  constructor(private api: APIService, private modal: NgbModal) {}

  public async ngOnInit() {
    this.fetch();
  }

  public async fetch() {
    this.bannedGuilds = (
      await this.api.getUnbanrequestBannedguilds().toPromise()
    ).data;
    this.requests = (await this.api.getUnbanrequests().toPromise()).data;
  }

  public async createRequest(guild: Guild) {
    this.newRequest.guild_id = guild.id;
    this.newRequest.guild = guild;

    const res = await this.modal.open(this.modalRequest, {
      windowClass: 'dark-modal',
    }).result;
  }
}
