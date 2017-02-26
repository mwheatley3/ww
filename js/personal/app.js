import React, { Component } from 'react';
import { BrowserRouter, Redirect, Route } from 'react-router-dom';
import { DragDropContext } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';

import Container from 'js/personal/components/container';
import Login from 'js/personal/components/login';
import Home from 'js/personal/components/home';

@DragDropContext(HTML5Backend)
export default class App extends Component {
    props: {
        store: Object,
    }

    static childContextTypes = {
        store: React.PropTypes.object.isRequired,
    };

    getChildContext() {
        const { store } = this.props;
        return { store };
    }

    componentWillMount() {
        this.props.store.init();

        window.addEventListener("unhandledrejection", function(e) {
            console.error(e);
        });

        window.addEventListener("rejectionhandled", function(e) {
            console.error(e);
        });
    }

    render() {
        return (
            <BrowserRouter>
                <Container>
                    <Route path="/" exact render={ () => <Redirect to="/client/home" push={ false } /> } />
                    <Route path="/login" component={ Login } />
                    <Route path="/client/home" component={ Home } />
                </Container>
            </BrowserRouter>
        );
    }
}
