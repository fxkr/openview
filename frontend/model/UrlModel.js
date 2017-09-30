export default class UrlModel {
  static getDirectory() {
    return window.location.pathname;
  }

  static getDirectoryInfoURL(nextPageToken) {
    const a = document.createElement('a');
    a.href = window.location.href;
    a.pathname = UrlModel.getDirectory();
    a.search = $.param({
      action: 'info',
      page_token: nextPageToken,
    });
    a.hash = '';
    return a.href;
  }

  static getDirectoryURL(relativePath) {
    const a = document.createElement('a');
    a.href = window.location.href;
    a.pathname = relativePath;
    a.search = '';
    a.hash = '';
    return a.href;
  }

  static getImageURL(relativePath) {
    const a = document.createElement('a');
    a.href = window.location.href;
    a.pathname = relativePath;
    a.search = '';
    a.hash = '';
    return a.href;
  }

  static getImageThumbnailURL(relativePath, thumbSize) {
    const a = document.createElement('a');
    a.href = window.location.href;
    a.pathname = relativePath;
    a.search = $.param({
      size: thumbSize,
    });
    a.hash = '';
    return a.href;
  }
}
