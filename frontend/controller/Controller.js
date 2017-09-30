class Controller {
  constructor(model, view) {
    this.model = model;
    this.view = view;

    view.more$.subscribe(this.onMore.bind(this));
    view.ready$.subscribe(this.onReady.bind(this));
  }

  onMore({ explicit }) {
    this.model.items.loadMore({ explicit });
  }

  onReady() {
    this.model.items.loadMore({ explicit: false });
  }
}

export default Controller;
