import hoist from 'hoist-non-react-statics';
import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { createPath } from 'react-router-dom/node_modules/history/PathUtils';
import Store, { withStore } from 'js/personal/store';
import { mergeQuery } from 'js/common/location';
import { autorun } from 'mobx';
import { observer } from 'mobx-react';

export default function(Comp) {
    class RequireAuthHOC extends Component {
        render() {
            return (
                <RequireAuth>
                    <Comp { ...this.props } />
                </RequireAuth>
            );
        }
    }

    hoist(RequireAuthHOC, Comp);
    return RequireAuthHOC;
}

@withRouter
@withStore
@observer
export class RequireAuth extends Component {
    props: {
        location: any | Object,
        store: any | Store,
        push: any | () => void,
        children?: any,
    };

    static defaultProps = {
        location: null, // provided by withRouter
        push: null, // provided by withRouter
        store: null, // provided by withStore
    };

    dispose: () => void
    state: { authed: boolean }

    constructor(...args: any) {
        super(...args);

        this.state = {
            authed: false,
        };
    }

    checkAuth() {
        const { location, push, store } = this.props;
        if (!store.auth.user.loaded) {
            return;
        }

        if (!store.auth.loggedIn) {
            const next = createPath(location);
            push(mergeQuery({ pathname: '/login' }, { next }));
            return;
        }

        if (!this.state.authed) {
            this.setState({ authed: true });
        }

        store.auth.loadProjects();
    }

    componentWillMount() {
        this.dispose = autorun(() => this.checkAuth());
    }

    componentWillUnmount() {
        this.dispose();
    }

    render() {
        const { authed } = this.state;
        const { store } = this.props;

        if (!authed || !store.auth.projects.value) {
            return null;
        }

        return this.props.children;
    }
}
