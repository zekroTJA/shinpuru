/** @format */

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { Guild } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

interface Severity {
  name: string;
  color: string;
}

@Component({
  selector: 'app-ga-data',
  templateUrl: './ga-data.component.html',
  styleUrls: ['./ga-data.component.sass'],
})
export class GuildAdminDataComponent implements OnInit {
  validator: string;
  leaveAfter: boolean;

  private guildID: string;
  private guild: Guild;

  @ViewChild('modalConfirm') private modalConfirm: TemplateRef<any>;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private api: APIService,
    private toasts: ToastService,
    private modals: NgbModal
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;
      this.guild = await this.api.getGuild(this.guildID).toPromise();
    });
  }

  async onGuildDataDelete() {
    this.validator = '';
    const res = await this.modals.open(this.modalConfirm, {
      windowClass: 'dark-modal',
    }).result;
    if (res) {
      if (this.validator !== this.guild.name) {
        this.toasts.push(
          'The entered guild name does not match!',
          'Validation failed',
          'error',
          10000
        );
        return;
      }
      await this.api
        .postGuildSettingsFlushGuildData(
          this.guildID,
          this.validator,
          this.leaveAfter
        )
        .toPromise();
      this.toasts.push(
        'Guild data successfully removed',
        'Guild data removed',
        'success',
        10000
      );
      if (this.leaveAfter) {
        this.router.navigate(['/guilds']);
      }
    }
  }
}
