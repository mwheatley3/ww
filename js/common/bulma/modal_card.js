// http://bulma.io/documentation/components/modal/

import React, { Component, PropTypes } from 'react';

import cx from 'classnames';

const { node, func, string } = PropTypes;

export default class ModalCard extends Component {
    static propTypes = {
        children: node,
        onCloseClick: func,
        title: node.isRequired,
        footer: node,
        className: string,
    };

    render() {
        const { children, title, footer, onCloseClick, className } = this.props;
        const cls = cx(className, 'modal', 'is-active');

        return (
            <div className={ cls }>
                <div className="modal-background" />
                <div className="modal-card">
                    <header className="modal-card-head">
                        <div className="modal-card-title">{ title }</div>
                        { onCloseClick ? <button className="delete" onClick={ onCloseClick } /> : null }
                    </header>
                    <section className="modal-card-body">
                        { children }
                    </section>
                    { footer ? <footer className="modal-card-foot">{ footer }</footer> : null }
                </div>
            </div>
        );
    }
}
