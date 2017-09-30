import './app.styl';
import './index.html';

import View from './view/View';
import Model from './model/Model';
import Controller from './controller/Controller';

function main() {
  const model = new Model();
  const view = new View(model);
  const controller = new Controller(model, view);

  // Expose to web developer tools
  window.debug = {
    model,
    view,
    controller,
  };

  model.ready();
}

$(window).on('load', main);
