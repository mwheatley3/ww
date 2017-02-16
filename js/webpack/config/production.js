var webpack = require('webpack');

module.exports = function(config) {
    config.devtool = false;

    config.plugins.push(
        new webpack.LoaderOptionsPlugin({
            debug: false,
            minimize: true,
        }),
        new webpack.optimize.UglifyJsPlugin({ sourceMap: false })
    );

    config.module.rules.push({
        test: /js\/common\/perf\.js$/,
        use: ['null-loader'],
    });

    return config;
};
