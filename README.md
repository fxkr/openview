# openview - A modern image gallery

Design criteria:

* Modern look, low distraction, high usability user interface. 
* Low-effort: Serves images from directory. Easy to deploy.
* Read-only, no database, only caches.
* Modern language.
* Avoid reinventing the wheel.
* Focus on stable and secure source code.
* Self-hostable.

The frontend was inspired by [Tumblr](https://www.tumblr.com/)'s low-distraction image wall, [flickr](http://www.flickr.com/)'s elegant layout and [fgallery](https://www.thregr.org/~wavexx/software/fgallery/)'s minimal UI.

License: MIT (see LICENSE file)


# Contributing

Thank you for your interest! Contributions to this project are very welcome. Here's an overview about our technology stack:

The frontend is written in JavaScript, of course.
We use [Justified Gallery](https://miromannino.github.io/Justified-Gallery/) to show an image wall and [PhotoSwipe](http://photoswipe.com) as a slideshow viewer.
We do use [jQuery](https://jquery.com/), which is *not* modern, but lightweight, very stable and it's a transitive dependency anyway.
Building/bundling is done by [webpack](https://webpack.js.org/).
The [yarn](https://yarnpkg.com/) package manager is highly recommended over npm.
We follow the [AirBNB coding style](https://github.com/airbnb/javascript) (enforced by [ESLint](https://eslint.org/)).

The backend is written in [Go](https://golang.org/), using the [Chi](https://github.com/go-chi/chi) web framework.
Server-side image operations are done using Go [bindings](https://github.com/gographics/imagick) for [ImageMagick](https://www.imagemagick.org/).
There are cache implementations based on files (for thumbnails) and [redigo](https://github.com/garyburd/redigo) (for metadata).
The latter can be used either with a real [redis](https://redis.io/) server or an embedded [miniredis](https://github.com/alicebob/miniredis) server.

In the backend we use the type system to maximize safety and security by avoiding basic types wherever possible.
For example, throughout the codebase absolute paths and relative paths are represented through safe absolute and relative path types, instead of strings.
This makes certain categories of issues like for example [path traversal](https://www.owasp.org/index.php/Path_Traversal) very unlikely.

We use [CircleCI 2.0](https://circleci.com/) for automated builds, tests and packaging.
The builds are based on a [Docker](https://www.docker.com/) container built and hosted on [Quay](https://quay.io/).
Release packages are uploaded to [packagecloud](https://packagecloud.io/).
