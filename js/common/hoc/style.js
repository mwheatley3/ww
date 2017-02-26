import hoist from 'hoist-non-react-statics';

export default function(...styles) {
    return function(Comp) {
        class Style extends Comp {
            static displayName = 'Style(' + (Comp.displayName || Comp.name) + ')';

            constructor(props: Object, context: Object) {
                for (let i = 0; i < styles.length; i++) {
                    styles[i].use();
                }

                super(props, context);
            }
        }

        hoist(Style, Comp);

        return Style;
    };
}
