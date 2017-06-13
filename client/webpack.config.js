var path = require("path");
var HtmlWebpackPlugin = require('html-webpack-plugin');

function getEntry() {
  const result = [];

  // the entry point of our app
  result.push('./src/index.js');

  return result;
}

function getSettings() {
  let settings = {};
  const isProd = process.env.NODE_ENV === 'production';

    settings.apiUrl = isProd ? '/api' : 'http://localhost:3030/api';

  return settings;
}

function getPlugins() {
  return [
    new HtmlWebpackPlugin(({
      template: './src/index.html',
      settings: JSON.stringify(getSettings())
    }))
  ];
}

const config = function(env) {
  return {
    entry: getEntry(),

    output: {
      path: path.resolve(__dirname + '/dist'),
      filename: 'app.js',
      publicPath: '/'
    },

    plugins: getPlugins(),

    module: {
      rules: [{
          test: /\.(css|scss)$/,
          use: [
            'style-loader',
            'css-loader',
              'sass-loader',
          ]
        },
        {
          test: /\.elm$/,
          exclude: [/elm-stuff/, /node_modules/],
          loader: 'elm-webpack-loader?verbose=true&warn=true',
        },
        {
          test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
          loader: 'url-loader?limit=10000&mimetype=application/font-woff',
        },
        {
          test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
          loader: 'file-loader',
        },
      ],

      noParse: /\.elm$/,
    },

    devServer: {
      inline: true,
      stats: {
        colors: true
      },
    },


  };
};

module.exports = function(env) {
  if (!env)
    env = {};

  console.log('Node Env: ' + process.env.NODE_ENV)
  console.log(env);
  
  return config(env);
};
