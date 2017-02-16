var base = require('./config/base'),
    debug = require('./config/debug'),
    production = require('./config/production');

var nodeEnv = process.env.NODE_ENV === 'production' ? 'production' : 'development';

var fns = [base];

// debug
if (nodeEnv === 'production') {
    fns.push(production);
} else {
    fns.push(debug);
}

module.exports = function(conf) {
    fns.push(conf);
    return fns.reduce(function(last, fn) {
        return fn(last);
    }, {});
};
