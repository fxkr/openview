import { Subject } from 'rxjs/Subject';

export default class NotificationModel {
  constructor() {
    this.notification$ = new Subject();
  }

  send(data) {
    this.notification$.next(data);
  }
}
