var webpack = require('webpack');

module.exports = function(config) {
    config.devtool = 'cheap-module-source-map';

    config.plugins.push(
        new webpack.LoaderOptionsPlugin({
            debug: true,
            minimize: false,
        })
    );

    return config;
};
