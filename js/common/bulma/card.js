// http://bulma.io/documentation/components/card/

import React, { Component, PropTypes } from 'react';

const { node, arrayOf } = PropTypes;

export default class Card extends Component {
    static propTypes = {
        children: node.isRequired,
        title: node,
        footerItems: arrayOf(node),
    };

    render() {
        const { children, title, footerItems, ...rest } = this.props;

        return (
            <div className="card" { ...rest }>
                { !title ? null : (
                    <header className="card-header">
                        <p className="card-header-title">{ title }</p>
                    </header>
                ) }
                <div className="card-content">
                    { children }
                </div>
                { !footerItems ? null : (
                    <footer className="card-footer">{
                        footerItems.map((it, i) => (
                            <div className="card-footer-item" key={ i }>{ it }</div>
                        ))
                    }</footer>
                ) }
            </div>
        );
    }
}
