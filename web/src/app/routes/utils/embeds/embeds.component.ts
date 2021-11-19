import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  Channel,
  ChannelType,
  ChannelWithPermissions,
  Guild,
  MessageEmbed,
  MessageEmbedField,
  MessageEmbedFooter,
  MessageEmbedImage,
  MessageEmbedThumbnail,
  MessageEmbedVideo,
} from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';

@Component({
  selector: 'app-embeds',
  templateUrl: './embeds.component.html',
  styleUrls: ['./embeds.component.scss'],
})
export class EmbedsComponent implements OnInit {
  embed = {
    fields: [],
    footer: {} as MessageEmbedFooter,
    image: {} as MessageEmbedImage,
    thumbnail: {} as MessageEmbedThumbnail,
    video: {} as MessageEmbedVideo,
  } as MessageEmbed;

  guild: Guild;
  erroneousJson = false;
  channelId: string;
  editMessageId: string;
  channels: ChannelWithPermissions[] = [];

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guild = await this.api.getGuild(params.guildid).toPromise();
      this.channels = (
        await this.api.getChannels(params.guildid).toPromise()
      ).data;
    });
  }

  get textChannels(): Channel[] {
    return this.channels.filter((c) => c.can_read && c.can_write);
  }

  addEmbedField() {
    this.embed.fields.push({} as MessageEmbedField);
  }

  removeEmbedField(i: number) {
    this.embed.fields.splice(i, 1);
  }

  jsonInput(v: string) {
    try {
      this.embed = JSON.parse(v);
      this.embed.color_hex = '#' + this.embed.color.toString(16);
      if (!this.embed.fields) this.embed.fields = [];
      if (!this.embed.footer) this.embed.footer = {} as MessageEmbedFooter;
      if (!this.embed.image) this.embed.image = {} as MessageEmbedImage;
      if (!this.embed.thumbnail)
        this.embed.thumbnail = {} as MessageEmbedThumbnail;
      if (!this.embed.video) this.embed.video = {} as MessageEmbedVideo;
      this.erroneousJson = false;
    } catch {
      this.erroneousJson = true;
    }
  }

  updateColor(v: string) {
    this.embed.color = parseInt(v.substr(1), 16);
  }

  async sendMessage() {
    if (this.editMessageId) {
      await this.api
        .postChannelsMessage(
          this.guild.id,
          this.channelId,
          this.editMessageId,
          this.embed
        )
        .toPromise();
      this.toasts.push('Embed successfully updated.', '', 'success', 6000);
    } else {
      this.editMessageId = (
        await this.api
          .postChannels(this.guild.id, this.channelId, this.embed)
          .toPromise()
      ).id;
      this.toasts.push('Embed successfully sent.', '', 'success', 6000);
    }
  }
}
