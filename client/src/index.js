'use strict';

require('font-awesome/css/font-awesome.css');
require('../styles/main.scss');

const Elm = require('./Main.elm');
const mountNode = document.getElementById('main');

// .embed() can take an optional second argument. This would be an object describing the data we need to start a program, i.e. a userID or some token
const app = Elm.Main.embed(mountNode, window.options);
