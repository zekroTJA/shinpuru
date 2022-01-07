/** @format */

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './routes/login/login.component';
import { HomeComponent } from './routes/home/home.component';
import { GuildComponent } from './routes/guild/guild.component';
import { MemberRouteComponent } from './routes/member/member.component';
import { SettingsComponent } from './routes/settings/settings.component';
import { InfoComponent } from './routes/info/info.component';
import { UserSettingsComponent } from './routes/usersettings/usersettingscomponent';
import { ScoreboardComponent } from './routes/scoreboard/scoreboard.component';
import { CommandsComponent } from './routes/commands/commands.component';
import { GuildAdminKarmaComponent } from './routes/guild-admin/karma/ga-karma.component';
import { GuildAdminAntiraidComponent } from './routes/guild-admin/antiraid/ga-antiraid.component';
import { GuildUnbanRequestComponent } from './routes/guild-unbanrequests/guild-unbanrequests.component';
import { MemberReportsComponent } from './routes/member-reports/member-reports.component';
import { UnbanmeComponent } from './routes/unbanme/unbanme.component';
import { GuildAdminGeneralComponent } from './routes/guild-admin/general/ga-general.component';
import { GuildAdminLogsComponent } from './routes/guild-admin/logs/ga-logs.component';
import { GuildAdminDataComponent } from './routes/guild-admin/data/ga-data.component';
import { DebugComponent } from './routes/debug/debug.component';
import { GuildAdminApiComponent } from './routes/guild-admin/api/ga-api.component';
import { EmbedsComponent } from './routes/utils/embeds/embeds.component';
import { VerifyRouteComponent } from './routes/verify/verify.component';
import { GuildAdminVerificationComponent } from './routes/guild-admin/verification/ga-verification.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: '/guilds',
    pathMatch: 'full',
  },
  {
    path: 'guilds',
    component: HomeComponent,
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: 'guilds/:id',
    component: GuildComponent,
  },
  {
    path: 'guilds/:guildid/scoreboard',
    component: ScoreboardComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/general',
    component: GuildAdminGeneralComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/antiraid',
    component: GuildAdminAntiraidComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/karma',
    component: GuildAdminKarmaComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/logs',
    component: GuildAdminLogsComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/data',
    component: GuildAdminDataComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/api',
    component: GuildAdminApiComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin/verification',
    component: GuildAdminVerificationComponent,
  },
  {
    path: 'guilds/:guildid/guildadmin',
    redirectTo: 'guilds/:guildid/guildadmin/general',
    pathMatch: 'full',
  },
  {
    path: 'guilds/:guildid/unbanrequests',
    component: GuildUnbanRequestComponent,
  },
  {
    path: 'guilds/:guildid/:memberid',
    component: MemberRouteComponent,
  },
  {
    path: 'guilds/:guildid/:memberid/reports',
    component: MemberReportsComponent,
  },
  {
    path: 'guilds/:guildid/utils/embeds',
    component: EmbedsComponent,
  },
  {
    path: 'settings',
    component: SettingsComponent,
  },
  {
    path: 'info',
    component: InfoComponent,
  },
  {
    path: 'usersettings',
    component: UserSettingsComponent,
  },
  {
    path: 'commands',
    component: CommandsComponent,
  },
  {
    path: 'unbanme',
    component: UnbanmeComponent,
  },
  {
    path: 'verify',
    component: VerifyRouteComponent,
  },
  {
    path: 'pogchamp',
    component: DebugComponent,
  },
  {
    path: '**',
    redirectTo: '/guilds',
    pathMatch: 'full',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' })],
  exports: [RouterModule],
})
export class AppRoutingModule {}
