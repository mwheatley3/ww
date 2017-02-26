import bcss from 'bulma/css/bulma.css';
import facss from 'font-awesome/css/font-awesome.css';
import css from './container.css';

import React, { Component, PropTypes } from 'react';
import style from 'js/common/hoc/style';
import { withStore } from 'js/personal/store';
import { observer } from 'mobx-react';

import DevTools from 'mobx-react-devtools';
import Nav from 'js/personal/components/nav';
import SidebarNav from 'js/personal/components/sidebar_nav';
import cx from 'classnames';

const { node, object } = PropTypes;

const sidebarKey = '__container_sidebar__';

let MobxTools = DevTools;
if (process.env.NODE_ENV === 'production') {
    MobxTools = () => null;
}

@withStore
@style(bcss, facss, css)
@observer
export default class Container extends Component {
    static propTypes = {
        children: node,
        store: object.isRequired,
    };

    constructor(...args) {
        super(...args);

        this.state = {
            sidebarOpen: sidebarKey in localStorage ? !!localStorage[sidebarKey] : true,
        };
    }

    onButtonClick = e => {
        this.setState(st => {
            localStorage[sidebarKey] = st.sidebarOpen ? '' : '1';
            return { sidebarOpen: !st.sidebarOpen };
        });
    }

    render() {
        const { children, store } = this.props;
        const { sidebarOpen } = this.state;
        const cls = cx('container-container', {
            'no-sidebar': !store.auth.loggedIn || !sidebarOpen,
        });

        return (
            <div className={ cls }>
                { store.auth.loggedIn ? [
                    <button key="button" onClick={ this.onButtonClick } />,
                    <SidebarNav key="nav" visible={ sidebarOpen } />,
                ] : null }
                <div className="main-content">
                    <Nav />
                    <div className="container">
                        <div>{ children }</div>
                    </div>
                </div>
                <MobxTools position={ { bottom: 0, right: 0 } } />
            </div>
        );
    }
}
