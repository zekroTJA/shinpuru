/** @format */

import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { UnbanRequest, UnbanRequestState, User } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';

// const typeColors = ['#d81b60', '#e53935', '#009688', '#fb8c00', '#8e24aa'];
const typeColors = ['#fb8c00', '#d81b60', '#8BC34A'];

@Component({
  selector: 'app-unbanrequest',
  templateUrl: './unbanrequest.component.html',
  styleUrls: ['./unbanrequest.component.scss'],
})
export class UnbanrequestComponent implements OnInit {
  @Input() public request: UnbanRequest;
  @Input() public showControls: boolean = false;

  @Output() public accept = new EventEmitter<any>();
  @Output() public decline = new EventEmitter<any>();
  @Output() public reports = new EventEmitter<any>();

  public processedBy: User;

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);
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
      .getUser(this.request.processed_by)
      .toPromise();
  }
}
