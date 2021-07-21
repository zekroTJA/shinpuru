/** @format */

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { Guild, UnbanRequest, UnbanRequestState } from 'src/app/api/api.models';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { ToastService } from 'src/app/components/toast/toast.service';

interface ExtendedUnbanRequest extends UnbanRequest {
  guild: Guild;
}

@Component({
  selector: 'app-unbanme',
  templateUrl: './unbanme.component.html',
  styleUrls: ['./unbanme.component.scss'],
})
export class UnbanmeComponent implements OnInit {
  @ViewChild('modalRequest') private modalRequest: TemplateRef<any>;

  public bannedGuilds: Guild[] = null;
  public requests: UnbanRequest[] = [];
  public newRequest = {} as ExtendedUnbanRequest;

  constructor(
    private api: APIService,
    private modal: NgbModal,
    private toasts: ToastService
  ) {}

  public async ngOnInit() {
    this.fetch();
  }

  public async fetch() {
    this.bannedGuilds = (
      await this.api.getUnbanrequestBannedguilds().toPromise()
    ).data;
    this.requests = (await this.api.getUnbanrequests().toPromise()).data.sort(
      this.requestsSortFunc
    );
  }

  public async createRequest(guild: Guild) {
    this.newRequest.guild_id = guild.id;
    this.newRequest.guild = guild;

    const res = await this.modal.open(this.modalRequest, {
      windowClass: 'dark-modal',
    }).result;
    if (res) {
      const res = await this.api.postUnbanrequests(this.newRequest).toPromise();
      this.requests = [res].concat(this.requests);
      this.toasts.push(
        'Unban request was submitted. An administrator will process the request as soon as possible.',
        'Unban request submitted',
        'success',
        6000,
        true
      );
    }
  }

  private requestsSortFunc(a: UnbanRequest, b: UnbanRequest) {
    if (a.status === UnbanRequestState.PENDING) return -1;
    if (b.status === UnbanRequestState.PENDING) return 1;
    if (typeof a.created === 'string') a.created = new Date(a.created);
    if (typeof b.created === 'string') b.created = new Date(b.created);
    return b.created.getTime() - a.created.getTime();
  }
}
