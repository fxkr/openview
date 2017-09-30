import { Subject } from 'rxjs/Subject';

import FragmentModel from './FragmentModel';
import ItemModel from './ItemModel';
import StateModel from './StateModel';
import NotificationModel from './NotificationModel';
import UrlModel from './UrlModel';

export default class Model {
  constructor() {
    this.urls = UrlModel;
    this.state = new StateModel();
    this.notifications = new NotificationModel();
    this.items = new ItemModel(this.state, this.notifications, this.urls);
    this.fragments = new FragmentModel(this.urls);

    this.ready$ = new Subject();
  }

  ready() {
    this.ready$.next();
  }
}
