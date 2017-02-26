import React, { PropTypes, Component } from 'react';
import { Creatable } from 'react-select';
import css from 'react-select/dist/react-select.css';
import style from 'js/common/hoc/style';
import css_combo_box from './combo_box.css';
import { omit } from 'lodash';

const { array, func, any } = PropTypes;

@style(css)
@style(css_combo_box)
export default class ComboBox extends Component {

    static propTypes = {
        options: array.isRequired,
        onChange: func.isRequired,
        value: any,
    }

    onChange = option => {
        const { onChange } = this.props;
        const newVal = option ? option.value : '';
        onChange(newVal);
    }

    render() {
        const props = omit(this.props, 'value', 'onChange');
        const value = { label: this.props.value, value: this.props.value };
        return (
            <Creatable { ...props } value={ value } onChange={ this.onChange } arrowRenderer={ () => {} }/>
        );
    }
}
