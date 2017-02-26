// http://bulma.io/documentation/elements/table/

import React, { Component, PropTypes } from 'react';
import cx from 'classnames';

const { bool, string } = PropTypes;

export default class Table extends Component {
    static propTypes = {
        borders: bool,
        stripes: bool,
        narrow: bool,
        className: string,
    };

    render() {
        const { borders, stripes, narrow, className, ...rest } = this.props;
        const cls = cx('table', {
            'is-bordered': borders,
            'is-striped': stripes,
            'is-narrow': narrow,
        }, className);

        return <table className={ cls } { ...rest } />;
    }
}
