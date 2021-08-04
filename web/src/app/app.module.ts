/** @format */

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { MarkdownModule } from 'ngx-markdown';

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
import { SpoilerComponent } from './components/spoiler/spoiler.component';
import { FormsModule } from '@angular/forms';
import { TagsInputComponent } from './components/tagsinput/tagsinput.component';
import { SettingsComponent } from './routes/settings/settings.component';
import { SpinnerButtonComponent } from './components/spinner-button/spinner-button.component';
import { InfoComponent } from './routes/info/info.component';
import { SliderComponent } from './components/slider/slider.component';
import { UserSettingsComponent } from './routes/usersettings/usersettingscomponent';
import { KarmaTileComponent } from './components/karma-tile/karma-tile.component';
import { ScoreboardComponent } from './routes/scoreboard/scoreboard.component';
import { KarmaScoreboardComponent } from './components/karma-scoreboard/karma-scoreboard.component';
import { CommandsComponent } from './routes/commands/commands.component';
import { KeyValuePipe } from './pipes/key-value.pipe';
import { CommandSortPipe } from './pipes/command-sort.pipe';
import { DiscordAssetPipe } from './pipes/discord-asset.pipe';
import { GuildAdminKarmaComponent } from './routes/guild-admin/karma/ga-karma.component';
import { GuildAdminAntiraidComponent } from './routes/guild-admin/antiraid/ga-antiraid.component';
import { GuildAdminNavbarComponent } from './components/guild-admin-navbar/guild-admin-navbar.component';
import { GuildUnbanRequestComponent } from './routes/guild-unbanrequests/guild-unbanrequests.component';
import { UnbanrequestComponent } from './components/unbanrequest/unbanrequest.component';
import { UnbanrequestStatePipe } from './pipes/unbanrequest-state.pipe';
import { MemberReportsComponent } from './routes/member-reports/member-reports.component';
import { UnbanmeComponent } from './routes/unbanme/unbanme.component';
import { AutocompleteInputComponent } from './components/autocomplete-input/autocomplete-input.component';
import { ProtipComponent } from './components/protip/protip.component';
import { StarboardEntryComponent } from './components/starboard-entry/starboard-entry.component';
import { StarboardComponent } from './components/starboard/starboard.component';
import { CapTextPipe } from './pipes/cap-text.pipe';
import AuthInterceptor from './api/auth.interceptor';
import { GuildAdminGeneralComponent } from './routes/guild-admin/general/ga-general.component';
import { GuildAdminLogsComponent } from './routes/guild-admin/logs/ga-logs.component';
import { GuildAdminDataComponent } from './routes/guild-admin/data/ga-data.component';
import LoadingInterceptor from './api/loading.interceptor';
import { SkeletonTileComponent } from './components/skeleton-tile/skeleton-tile.component';
import { DebugComponent } from './routes/debug/debug.component';
import { GlobalSearchComponent } from './components/global-search/global-search.component';
import { GuildAdminApiComponent } from './routes/guild-admin/api/ga-api.component';

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
    SpoilerComponent,
    TagsInputComponent,
    SettingsComponent,
    SpinnerButtonComponent,
    InfoComponent,
    SliderComponent,
    UserSettingsComponent,
    KarmaTileComponent,
    ScoreboardComponent,
    KarmaScoreboardComponent,
    CommandsComponent,
    KeyValuePipe,
    CommandSortPipe,
    DiscordAssetPipe,
    GuildAdminKarmaComponent,
    GuildAdminAntiraidComponent,
    GuildAdminNavbarComponent,
    GuildAdminGeneralComponent,
    GuildAdminLogsComponent,
    GuildAdminDataComponent,
    GuildUnbanRequestComponent,
    UnbanrequestComponent,
    UnbanrequestStatePipe,
    MemberReportsComponent,
    UnbanmeComponent,
    AutocompleteInputComponent,
    ProtipComponent,
    StarboardComponent,
    StarboardEntryComponent,
    CapTextPipe,
    SkeletonTileComponent,
    DebugComponent,
    GlobalSearchComponent,
    GuildAdminApiComponent,
  ],
  imports: [
    NgbModule,
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule,
    MarkdownModule.forRoot(),
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: LoadingInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
