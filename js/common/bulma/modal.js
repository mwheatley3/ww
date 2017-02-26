// http://bulma.io/documentation/components/modal/

import React, { Component, PropTypes } from 'react';
import cx from 'classnames';

const { node, func, string } = PropTypes;

export default class Modal extends Component {
    static propTypes = {
        children: node,
        onCloseClick: func,
        className: string,
    };

    render() {
        const { children, className, onCloseClick } = this.props;
        const cls = cx(className, 'modal', 'is-active');

        return (
            <div className={ cls }>
                <div className="modal-background" />
                <div className="modal-content">
                    { children }
                </div>
                { onCloseClick ? <button className="modal-close" onClick={ onCloseClick } /> : null }
            </div>
        );
    }
}
