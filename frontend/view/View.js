import 'justifiedGallery/dist/css/justifiedGallery.css';
import 'justifiedGallery/dist/js/jquery.justifiedGallery';
import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.png';
import 'photoswipe/dist/default-skin/preloader.gif';
import 'photoswipe/dist/default-skin/default-skin.svg';
import 'photoswipe/dist/default-skin/default-skin.css';
import { Subject } from 'rxjs/Subject';
import 'noty/lib/noty.css';
import Noty from 'noty';

const PLACEHOLDER_IMAGE = 'data:image/gif;base64,R0lGODlhAQABAAAAACwAAAAAAQABAAA=';

const IMAGE_DATA_KEY = 'image-data';
const PHOTOSWIPE_IMAGE_DATA_KEY = 'photoswipe-image-data';

class View {
  constructor(model) {
    this.model = model;

    model.ready$.subscribe(() => {
      this.onReadyInitNavbar();
      this.onReadyInitLoadMore();
      this.onReadyInitScroll();
      this.onReadyInitGallery();
      this.ready$.next();
    });

    model.items.newImages$.subscribe(this.onNewImagesExtendGallery.bind(this));
    model.items.newDirectories$.subscribe(this.onNewDirectoriesShowDirectories.bind(this));
    model.notifications.notification$.subscribe(this.onNotificationNotify.bind(this));
    model.state.changed$.subscribe(this.onStateChanged.bind(this));

    this.more$ = new Subject();
    this.ready$ = new Subject();

    this.navbar = $('.navbar-directory');
    this.gallery = $('.gallery');
    this.directoryList = $('.dir-list');
    this.photoswipe = $('.pswp');
    this.moreButton = $('.load-more');

    this.pswpImages = [];
  }

  onReadyInitNavbar() {
    const homeFragment = this.model.fragments.getHomeFragment();

    this.navbar.append($('<a>', {
      class: 'breadcrumb',
      href: homeFragment.url,
    }).text(homeFragment.name));

    const fragments = this.model.fragments.getFragments();
    $(fragments).each((i, fragment) => {
      const divider = $('<span>', {
        class: 'breadcrumb-divider',
      }).text('/');

      const crumb = $('<a>', {
        class: 'breadcrumb',
        href: fragment.url,
      }).text(fragment.name);

      this.navbar.append(divider);
      this.navbar.append(crumb);
    });

    this.model.ready$.unsubscribe(this.onReadyInitNavbar);
  }

  onReadyInitLoadMore() {
    $('.load-more-btn').click(() => {
      $('.load-more-btn').blur(); // Keep spacebar working
      this.more$.next({ explicit: true });
    });

    this.model.ready$.unsubscribe(this.onReadyInitLoadMore);
  }

  onReadyInitScroll() {
    $(window).on('scroll', () => {
      if (this.isGalleryAtBottom()) {
        this.more$.next({ explicit: false });
      }
    });

    this.model.ready$.unsubscribe(this.onReadyInitScroll);
  }

  onReadyInitGallery() {
    this.gallery.justifiedGallery({
      rowHeight: 200,
      maxRowHeight: 400,
      margins: 3,

      thumbnailPath: this.getThumbnailPath.bind(this),

      cssAnimation: true,
      waitThumbnailsLoad: false,
    });

    this.gallery.justifiedGallery().on('jg.complete jg.resize', () => {
      if (this.isGalleryAtBottom()) {
        this.more$.next({ explicit: false });
      }
    });

    this.model.ready$.unsubscribe(this.onReadyInitGallery);
  }

  onNotificationNotify({ text, type }) { // eslint-disable-line class-methods-use-this
    new Noty({
      text,
      type,
      layout: 'bottomCenter',
    }).show();
  }

  onNewImagesExtendGallery(newImages) {
    newImages.forEach((newImageData) => {
      const { width, height } = newImageData;

      const img = $('<img>', {
        src: PLACEHOLDER_IMAGE, // Don't temporarily load wrong thumbnail size.
        width,
        height,
        alt: newImageData.name,
      });

      // Store metadata on image. We'll use it to determine thumbnail URLs.
      img.data(IMAGE_DATA_KEY, newImageData);

      // Tell JustifiedGallery about ratio.
      // Without, it would infer the ratio from the 1x1 placeholder image.
      img.data('width', width);
      img.data('height', height);

      const imgLink = $('<a>', {
        href: newImageData.url,
        width: newImageData.width,
        height: newImageData.height,
      });
      imgLink.append(img);
      imgLink.appendTo(this.gallery);

      // When image is clicked, open slideshow at same position
      imgLink.on('click', 'a,img', img, (event) => {
        event.preventDefault();
        this.showSlideshow(event.data);
      });

      // Serve thumbnail, not raw image.
      // Benefits:
      // 1. Very large pictures don't slow slideshow down.
      // 2. Thumbnails already have corrected orientation (difficult to get right in browser).
      // 3. Later, we can extend this with optimizations for mobile and high-dpi screens.
      const maxSize = 2048;
      let pswpWidth = width;
      let pswpHeight = height;
      if (width > maxSize || height > maxSize) {
        if (width > height) {
          pswpWidth = maxSize;
          pswpHeight /= (width / maxSize);
        } else {
          pswpWidth /= (height / maxSize);
          pswpHeight = maxSize;
        }
      }
      const pswpSrc = this.model.urls.getImageThumbnailURL(newImageData.relativePath, maxSize);

      // Append image to PhotoSwipe gallery.
      // Always use the same array so it'll work even if PhotoSwipe already open.
      const pswpImage = {
        index: newImageData.index,
        w: pswpWidth,
        h: pswpHeight,
        src: pswpSrc,
        title: newImageData.name,
        elem: img[0],
      };
      this.pswpImages.push(pswpImage);

      // Store PhotoSwipe data on image so we can later make PhotoSwipe use the same
      // thumbnail image JustifiedGallery has already loaded as placeholder.
      img.data(PHOTOSWIPE_IMAGE_DATA_KEY, pswpImage);
    });

    // Layout new images if necessary. (Don't re-layout existing rows.)
    if (newImages) {
      this.gallery.justifiedGallery('norewind');
    }
  }

  showSlideshow(currentImg) {
    const currentImgData = currentImg.data(IMAGE_DATA_KEY);

    const options = {
      index: currentImgData.index,
      bgOpacity: 1.0,
      showHideOpacity: true,
      loop: false,
      closeOnScroll: false,
      history: false,
    };

    const pswp = new PhotoSwipe(
      this.photoswipe[0],
      PhotoSwipeUI_Default,
      this.pswpImages,
      options,
    );

    // Scroll page so gallery keeps up with slideshow.
    // Useful side-effect: loads more images as necessary via infinite scrolling.
    // Note: with beforeChange, it would be noticeable during the slideshow opening animation.
    pswp.listen('afterChange', () => {
      pswp.currItem.elem.scrollIntoView();
    });

    pswp.init();
  }

  onNewDirectoriesShowDirectories(newDirectories) {
    $(newDirectories).each((index, directory) => {
      const elemLink = $('<a>', {
        href: directory.url,
        class: 'dir-list-item',
      });

      const elem = $('<li>', {
        class: 'dir-list-item',
      }).text(directory.name);

      elemLink.append(elem);
      elemLink.appendTo(this.directoryList);
    });

    if (newDirectories) {
      this.directoryList.removeClass('dir-list-hidden');
    }
  }

  onStateChanged(state) {
    if (state.isEnd()) {
      this.moreButton.remove();
    }
  }

  isGalleryAtBottom() {
    const controller = this.gallery.data('jg.controller');
    const { settings } = controller;

    const bottomOfGallery = this.gallery.offset().top + this.gallery.height();
    const focusPoint = bottomOfGallery - settings.rowHeight;
    const bottomOfWindow = $(window).scrollTop() + $(window).height();

    return focusPoint <= bottomOfWindow;
  }

  getThumbnailPath(currentPath, width, height, image) {
    const availableSizes = [100, 240, 360, 500, 800, 1024, 1600, 2048];

    const displayedSize = Math.max(width, height);
    const neededSize = displayedSize * 2; // Nyquist
    const requestedSize = availableSizes.find(x => x >= neededSize)
      || availableSizes[availableSizes.length - 1];

    const { relativePath } = $(image).data(IMAGE_DATA_KEY);

    const thumbnailUrl = this.model.urls.getImageThumbnailURL(relativePath, requestedSize);

    // Use thumbnail as placeholder image in PhotoSwipe since it'll already be loaded.
    $(image).data(PHOTOSWIPE_IMAGE_DATA_KEY).msrc = thumbnailUrl;

    return thumbnailUrl;
  }
}

export default View;
