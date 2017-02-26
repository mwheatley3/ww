import React, { Component, PropTypes } from 'react';
import cx from 'classnames';

const { string } = PropTypes;

export default class Icon extends Component {
    static propTypes = {
        icon: string.isRequired,
        className: string,
    };

    render() {
        const { icon, className, ...rest } = this.props;
        const cls = cx('icon', 'fa', 'fa-' + icon, className);

        return (
            <i className={ cls } { ...rest } />
        );
    }
}
