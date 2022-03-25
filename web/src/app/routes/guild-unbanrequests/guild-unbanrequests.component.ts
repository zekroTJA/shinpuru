/** @format */

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';
import { UnbanRequest, UnbanRequestState } from 'src/app/api/api.models';
import { ToastService } from 'src/app/components/toast/toast.service';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-guild-unbanrequests',
  templateUrl: './guild-unbanrequests.component.html',
  styleUrls: ['./guild-unbanrequests.component.scss'],
})
export class GuildUnbanRequestComponent implements OnInit {
  @ViewChild('modalProcess') private modalProcess: TemplateRef<any>;

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

  public unbanRequests: UnbanRequest[];
  public isAccept = false;
  public processMessage: string;
  public selected: UnbanRequest;

  private guildID: string;

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
    private toats: ToastService,
    private router: Router,
    private modal: NgbModal
  ) {}

  public async ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;
      this.fetch();
    });
  }

  public async fetch() {
    this.unbanRequests = (
      await this.api.getGuildUnbanrequests(this.guildID).toPromise()
    ).data.sort(this.requestsSortFunc);
  }

  public onReports(r: UnbanRequest) {
    this.router.navigate(['guilds', this.guildID, r.user_id, 'reports']);
  }

  public onAccept(r: UnbanRequest) {
    this.isAccept = true;
    this.process(r);
  }

  public onDecline(r: UnbanRequest) {
    this.isAccept = false;
    this.process(r);
  }

  private async process(r: UnbanRequest) {
    this.selected = r;
    this.processMessage = '';
    const res = await this.modal.open(this.modalProcess, {
      windowClass: 'dark-modal',
    }).result;
    if (res) {
      await this.api
        .postGuildUnbanrequest(this.guildID, {
          id: this.selected.id,
          status: this.isAccept
            ? UnbanRequestState.ACCEPTED
            : UnbanRequestState.DECLINED,
          processed_message: this.processMessage,
        } as UnbanRequest)
        .toPromise();
      this.fetch();
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
