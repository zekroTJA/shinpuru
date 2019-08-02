/** @format */

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './routes/login/login.component';
import { HeaderComponent } from './components/header/header.component';
import { ToastComponent } from './components/toast/toast.component';
import { HomeComponent } from './routes/home/home.component';
import { SpinnerComponent } from './components/spinner/spinner.component';
import { GuildComponent } from './routes/guild/guild.component';
import { PopupComponent } from './components/popup/popup.component';
import { RoleBadgeComponent } from './components/rolebadge/rolebadge.component';
import { MemberComponent } from './components/member/member.component';
import { MemberRouteComponent } from './routes/member/member.component';
import { ReportComponent } from './components/report/report.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    HeaderComponent,
    ToastComponent,
    HomeComponent,
    SpinnerComponent,
    GuildComponent,
    PopupComponent,
    RoleBadgeComponent,
    MemberComponent,
    MemberRouteComponent,
    ReportComponent,
  ],
  imports: [NgbModule, BrowserModule, AppRoutingModule, HttpClientModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
