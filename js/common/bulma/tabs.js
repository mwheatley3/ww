// http://bulma.io/documentation/components/tabs/

import React, { Component, PropTypes } from 'react';
import cx from 'classnames';

const { string, bool, arrayOf, object, oneOf } = PropTypes;

export default class Tabs extends Component {
    static propTypes = {
        tabs: arrayOf(object).isRequired,
        align: oneOf(['left', 'center', 'right', 'full-width']),
        size: oneOf(['small', 'normal', 'medium', 'large']),
        boxed: bool,
        toggle: bool,
        active: string.isRequired,
    };

    static defaultProps = {
        align: 'left',
        size: 'normal',
    };

    render() {
        const { tabs, active, align, size, boxed, toggle } = this.props;

        const cls = cx('tabs', 'is-' + align, {
            ['is-' + size]: !!size && size !== 'normal',
            'is-boxed': boxed,
            'is-toggle': toggle,
        });

        return (
            <div className={ cls }>
                <ul>{
                    tabs.map(t => {
                        let tcls = '';
                        if (t.name === active) {
                            tcls = 'is-active';
                        }

                        return (
                            <li key={ t.name } className={ tcls }>
                                { t.render }
                            </li>
                        );
                    })
                }</ul>
            </div>
        );
    }
}
