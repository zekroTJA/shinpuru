/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import dateFormat from 'dateformat';
import { APIToken } from 'src/app/api/api.models';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-apitoken',
  templateUrl: './apitoken.component.html',
  styleUrls: ['./apitoken.component.sass'],
})
export class APITokenComponent {
  public dateFormat = dateFormat;

  public token: APIToken;
  public notGenerated = false;
  public revealToken = false;

  constructor(private api: APIService, private toats: ToastService) {
    this.fetch();
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

  private fetch() {
    this.api
      .getAPIToken(true)
      .pipe(catchError((err) => of(null)))
      .subscribe((data) => {
        this.notGenerated = data == null;
        this.token = data;
      });
  }
}
