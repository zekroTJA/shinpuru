/** @format */

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';
import { APIToken, User, UserSettingsOTA } from 'src/app/api/api.models';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { ToastService } from 'src/app/components/toast/toast.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-usersettings',
  templateUrl: './usersettings.component.html',
  styleUrls: ['./usersettings.component.scss'],
})
export class UserSettingsComponent implements OnInit {
  @ViewChild('modalConfirm') private modalConfirm: TemplateRef<any>;

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

  public token: APIToken;
  public notGenerated = false;
  public revealToken = false;

  public ota: UserSettingsOTA;

  public selfUser: User;
  public validator: string;

  constructor(
    private api: APIService,
    private toats: ToastService,
    private modals: NgbModal,
    private toasts: ToastService
  ) {}

  public async ngOnInit() {
    this.fetch();
    this.selfUser = await this.api.getSelfUser().toPromise();
  }

  public resetToken() {
    this.api.deleteAPIToken().subscribe(() => {
      this.fetch();
    });
  }

  public generateToken() {
    this.api.postAPIToken().subscribe((data) => {
      this.notGenerated = false;
      this.token = data;
    });
  }

  public copyTokenToClipboard() {
    const selBox = document.createElement('textarea');
    selBox.style.position = 'fixed';
    selBox.style.left = '0';
    selBox.style.top = '0';
    selBox.style.opacity = '0';
    selBox.value = this.token.token;
    document.body.appendChild(selBox);
    selBox.focus();
    selBox.select();
    document.execCommand('copy');
    document.body.removeChild(selBox);
    this.toats.push(
      'Copied token to clipboard.',
      'Token copied',
      'success',
      4000,
      true
    );
  }

  public onOTASave() {
    this.api.postUserSettingsOTA(this.ota).subscribe((data) => {
      this.ota = data;
      this.toats.push(
        'One Time Auth settings successfully updated.',
        'OTA Settings Updated',
        'green',
        3000
      );
    });
  }

  public async flushData() {
    this.validator = '';
    try {
      const res = await this.modals.open(this.modalConfirm, {
        windowClass: 'dark-modal',
      }).result;
      if (!res) return;
      if (this.validator !== this.selfUser.username) {
        this.toasts.push(
          'The entered user name does not match!',
          'Validation failed',
          'error',
          10000
        );
        return;
      }
      const fres = await this.api.postUserSettingsFlush().toPromise();
      const fresAssembled = Object.keys(fres)
        .map((k) => `${k}: ${fres[k]}, `)
        .join('\n');
      this.toasts.push(
        'Removed entries: ' + fresAssembled,
        'User data removed!',
        'success'
      );
    } catch {}
  }

  private fetch() {
    this.api
      .getAPIToken(true)
      .pipe(catchError((err) => of(null)))
      .subscribe((data) => {
        this.notGenerated = data == null;
        this.token = data;
      });

    this.api.getUserSettingsOTA().subscribe((data) => {
      console.log(data);
      this.ota = data;
    });
  }
}
