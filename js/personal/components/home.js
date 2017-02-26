import React, { Component } from 'react';
import requireAuth from 'js/personal/hoc/require_auth';

@requireAuth
export default class Home extends Component {
    static propTypes = {
    };

    render() {
        return (
            <div>
                <div>Home</div>
            </div>
        );
    }
}
