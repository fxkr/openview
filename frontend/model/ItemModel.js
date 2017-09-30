import { Subject } from 'rxjs/Subject';

export default class ItemModel {
  constructor(stateModel, notificationModel, urlModel) {
    this.state = stateModel;
    this.notificationModel = notificationModel;
    this.urlModel = urlModel;

    this.images = [];
    this.directories = [];

    this.newImages$ = new Subject();
    this.newDirectories$ = new Subject();
  }

  loadMore({ explicit }) {
    if (!this.state.setLoading(explicit)) {
      return;
    }

    $.get({
      url: this.urlModel.getDirectoryInfoURL(this.state.getNextPageToken()),
      dataType: 'json',
    }).done(this.loadMoreDone.bind(this))
      .fail(this.loadMoreFailed.bind(this));
  }

  loadMoreDone(data) {
    const newDirectories = [];
    const newImages = [];

    data.directories.forEach((directoryData) => {
      newDirectories.push({
        name: directoryData.name,
        relativePath: directoryData.relative_path,

        url: this.urlModel.getDirectoryURL(directoryData.relative_path),
      });
    });

    const baseIndex = this.images.length;
    data.images.forEach((imageData, i) => {
      newImages.push({
        name: imageData.name,
        relativePath: imageData.relative_path,

        width: imageData.width,
        height: imageData.height,

        url: this.urlModel.getImageURL(imageData.relative_path),

        index: baseIndex + i,
      });
    });

    if (newDirectories) {
      this.directories.push(...newDirectories);
      this.newDirectories$.next(newDirectories);
    }
    if (newImages) {
      this.images.push(...newImages);
      this.newImages$.next(newImages);
    }

    this.state.setSuccess(data.next_page_token);
  }

  loadMoreFailed() {
    this.state.setFailure();
    this.notificationModel.send({
      text: 'Sorry, something went wrong... try again?',
      type: 'error',
    });
  }
}
