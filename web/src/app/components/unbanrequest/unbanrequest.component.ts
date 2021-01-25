/** @format */

import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import {
  Report,
  Member,
  UnbanRequest,
  UnbanRequestState,
} from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import dateFormat from 'dateformat';

// const typeColors = ['#d81b60', '#e53935', '#009688', '#fb8c00', '#8e24aa'];
const typeColors = ['#fb8c00', '#d81b60', '#8BC34A'];

@Component({
  selector: 'app-unbanrequest',
  templateUrl: './unbanrequest.component.html',
  styleUrls: ['./unbanrequest.component.sass'],
})
export class UnbanrequestComponent implements OnInit {
  @Input() public request: UnbanRequest;
  @Input() public showControls: boolean = false;

  @Output() public accept = new EventEmitter<any>();
  @Output() public decline = new EventEmitter<any>();
  @Output() public reports = new EventEmitter<any>();

  public processedBy: Member;

  public dateFormat = dateFormat;
  public UnbanRequestState = UnbanRequestState;

  constructor(private api: APIService) {}

  async ngOnInit() {
    if (this.request.status !== UnbanRequestState.PENDING) {
      this.fetchProcessedBy();
    }
  }

  public get typeColor(): string {
    return typeColors[this.request.status] || 'gray';
  }

  public async fetchProcessedBy() {
    this.processedBy = await this.api
      .getGuildMember(this.request.guild_id, this.request.processed_by)
      .toPromise();
  }
}
