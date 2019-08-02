/** @format */

import { Component, Input } from '@angular/core';
import { Role } from 'src/app/api/api.models';
import { toHexClr } from 'src/app/utils/utils';

@Component({
  selector: 'app-rolebadge',
  templateUrl: './rolebadge.component.html',
  styleUrls: ['./rolebadge.component.sass'],
})
export class RoleBadgeComponent {
  @Input() public role: Role;

  public toHexClr = toHexClr;
}
