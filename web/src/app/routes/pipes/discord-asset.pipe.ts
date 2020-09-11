/** @format */

import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'discordAsset',
})
export class DiscordAssetPipe implements PipeTransform {
  transform(
    assetUrl: string,
    alt: string,
    res?: number,
    cond?: boolean
  ): string {
    if (cond === false || !assetUrl) {
      return alt;
    }

    if (res) {
      assetUrl = `${assetUrl}?size=${res}`;
    }

    return assetUrl;
  }
}
