import { Component, OnInit } from '@angular/core';
import { User } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-verify',
  templateUrl: './verify.component.html',
  styleUrls: ['./verify.component.scss'],
})
export class VerifyRouteComponent implements OnInit {
  selfUser: User;
  captchaSiteKey: string;

  constructor(private api: APIService, private toasts: ToastService) {}

  ngOnInit() {
    this.api.getSelfUser().subscribe((user) => {
      this.selfUser = user;
    });
    this.api.getVerificationSiteKey().subscribe((siteKey) => {
      this.captchaSiteKey = siteKey.sitekey;
    });
  }

  async onCaptchaError(error: any) {
    this.toasts.push(error, 'Failed verifying captcha', 'error');
  }

  async onCaptchaVerify(token: string) {
    await this.api.postVerifyUser(token).toPromise();
    this.toasts.push(
      'Your account has successfully been veified!',
      'Verification successful',
      'success'
    );
    this.selfUser.captcha_verified = true;
  }

  onCaptchaExpired(response: any) {
    window.location.reload();
  }
}
