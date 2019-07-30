/** @format */

export class User {
  constructor(
    public id: string,
    public username: string,
    public avatar: string,
    public locale: string,
    public discriminator: string,
    public verified: boolean,
    public bot: boolean,
    // tslint:disable-next-line: variable-name
    public avatar_url: string
  ) {}
}
