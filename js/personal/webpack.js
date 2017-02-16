var path = require('path');
var config = require('../webpack');

module.exports = config(function(conf) {
    conf.entry = ['babel-polyfill', path.join('js', path.basename(__dirname), 'index.js')];
    conf.output.path = path.join(conf.output.path, path.basename(__dirname));

    return conf;
});
