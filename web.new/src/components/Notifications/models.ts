export enum NotificationType {
  INFO,
  SUCCESS,
  WARNING,
  ERROR,
}

export type NotificationMeta = {
  type?: NotificationType;
  delay?: number;
  hide?: boolean;
  uid?: number;
};

export type Notification = NotificationMeta & {
  heading?: string;
  message: string | JSX.Element;
};
