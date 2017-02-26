import css from './nav.css';

import React, { Component, PropTypes } from 'react';

import { withRouter } from 'react-router-dom';
import { observer } from 'mobx-react';
import { withStore } from 'js/personal/store';
import style from 'js/common/hoc/style';

import { Link } from 'react-router-dom';
import imageURL from 'js/personal/imgs/hamburger.png';
import Menu from 'js/personal/components/menu';

const { object, func } = PropTypes;

@withRouter
@withStore
@style(css)
@observer
export default class Nav extends Component {
    static propTypes = {
        store: object.isRequired,
        push: func.isRequired,
    };

    constructor(...args) {
        super(...args);

        this.state = {
            menuVisible: false,
        };
    }

    onLogoutClick() {
        const { store, push } = this.props;

        store.auth.logout();
        push('/login');
    }

    toggleMenu = () => {
        this.setState(st => ({ menuVisible: !st.menuVisible }));
    }

    render() {
        const { store } = this.props;
        const { menuVisible } = this.state;

        return (
            <nav className="nav">
                <Link className="brand" to="/"><img src={ imageURL } /></Link>
                {
                    store.auth.loggedIn ?
                        <div className="user-icon">
                            <Menu visible={ menuVisible } onClose={ this.toggleMenu }>
                                <a onClick={ () => this.onLogoutClick() }>Log Out</a>
                            </Menu>
                        </div> :
                    null
                }
                <div className="user">{
                    store.auth.loggedIn ?
                        <span key="1">{ store.auth.user.value.email }</span> :
                        null
                }</div>
            </nav>
        );
    }
}
