import css from './login.css';
import React, { Component, PropTypes } from 'react';
import { observer } from 'mobx-react';

import Button from 'js/common/bulma/button';

import style from 'js/common/hoc/style';
import { withStore } from 'js/personal/store';
import { getQuery } from 'js/common/location';

const { object, func } = PropTypes;

@withStore
@style(css)
@observer
export default class Login extends Component {
    static propTypes = {
        store: object.isRequired,
        location: object.isRequired,
        push: func.isRequired,
    };

    onSubmit = e => {
        e.preventDefault();

        const { store } = this.props;

        if (store.auth.user.loading) {
            return;
        }

        const { email, password } = this.refs;
        store.auth.login(email.value, password.value);
    }

    componentWillMount() {
        this.checkAuth();
    }

    componentDidUpdate() {
        this.checkAuth();
    }

    checkAuth() {
        const { store, push, location } = this.props;
        if (store.auth.loggedIn) {
            const next = getQuery(location, 'next') || '/';
            push(next);
        }
    }

    render() {
        const { store } = this.props;

        if (store.auth.loggedIn) {
            return null;
        }

        const isUnauthorized = store.auth.user.error && store.auth.user.error.statusCode === 401 || false;
        const loading = store.auth.user.loading;

        return (
            <div className="login">
                <form onSubmit={ this.onSubmit }>
                    <h1>Sign In</h1>
                    { (isUnauthorized && !loading) ? <div className="error">Invalid email or password</div> : null }
                    <div className="login-input">
                        <input type="text" ref="email" placeholder="Email" />
                    </div>
                    <div className="login-input">
                        <input type="password" ref="password" placeholder="Password" />
                    </div>
                    <Button component="button" loading={ loading } disabled={ loading }>Login</Button>
                </form>
            </div>
        );
    }
}
