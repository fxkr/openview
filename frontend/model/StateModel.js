import { Subject } from 'rxjs/Subject';

export default class StateModel {
  constructor() {
    this.value = StateModel.Values.READY;
    this.nextPageToken = undefined;

    this.changed$ = new Subject();
  }

  setLoading(explicit) {
    if (this.canLoadMore(explicit)) {
      this.value = StateModel.Values.LOADING;
      return true;
    }
    return false;
  }

  setSuccess(nextPageToken) {
    this.nextPageToken = nextPageToken;
    if (nextPageToken) {
      this.value = StateModel.Values.READY;
    } else {
      this.value = StateModel.Values.END;
    }
    this.changed$.next(this);
  }

  setFailure() {
    this.value = StateModel.Values.FAILED;
    this.changed$.next(this);
  }

  canLoadMore(explicit) {
    switch (this.value) {
      case StateModel.Values.READY:
        return true;
      case StateModel.Values.FAILED:
        return !!explicit;
      default:
        return false;
    }
  }

  isEnd() {
    return this.value === StateModel.Values.END;
  }

  getNextPageToken() {
    return this.nextPageToken;
  }
}

StateModel.Values = Object.freeze({
  // Ready to attempt to load more items.
  READY: Symbol('ready'),

  // A request for more items is currently in-flight.
  LOADING: Symbol('loading'),

  // A request has failed. No more auto-loading, but manual retry possible.
  FAILED: Symbol('failed'),

  // There are no more items to load.
  END: Symbol('end'),
});
