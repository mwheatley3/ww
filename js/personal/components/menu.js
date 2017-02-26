import css from './menu.css';
import React, { Component, PropTypes, Children } from 'react';
import style from 'js/common/hoc/style';
import cx from 'classnames';

const { node, bool, func, string } = PropTypes;

@style(css)
export default class Menu extends Component {
    static propTypes = {
        children: node,
        visible: bool,
        onClose: func.isRequired,
        className: string,
    };

    componentWillMount() {
        document.addEventListener('click', this.onDocClick);
    }

    componentWillReceiveProps(props) {
        if (props.visible) {
            document.addEventListener('click', this.onDocClick);
        } else {
            document.removeEventListener('click', this.onDocClick);
        }
    }

    componentWillUnmount() {
        document.removeEventListener('click', this.onDocClick);
    }

    onDocClick = e => {
        const { onClose, visible } = this.props;
        if (visible) {
            onClose();
        }
    }

    render() {
        const { children, visible, className, onClose } = this.props;
        if (!visible) {
            return null;
        }

        return (
            <ul className={ cx('menu', className) }>{
                Children.map(children, ch => (
                    <li onClick={ () => onClose() } >{ ch }</li>
                ))
            }</ul>
        );
    }
}
